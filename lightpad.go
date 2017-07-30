package libplumraw

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"
)

// SetLogicalLoadLevel is used to both toggle and dim switches
func (l *DefaultLightpad) SetLogicalLoadLevel(level int) error {
	pd := struct {
		Level int    `json:"level"`
		LLID  string `json:"llid"`
	}{level, l.LLID}
	fmt.Printf("%+v\n", pd)
	resp, err := l.makePadPOSTRequest(pathSetLogicalLoadLevel, pd)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to set load level: status %s", resp.Status)
	}
	return nil
}

// SetLogicalLoadConfig
func (l *DefaultLightpad) SetLogicalLoadConfig(conf LogicalLoadConfig) error {
	pd := struct {
		Config LogicalLoadConfig `json:"config"`
		LLID   string            `json:"llid"`
	}{conf, l.LLID}
	resp, err := l.makePadPOSTRequest(pathSetLogicalLoadConfig, pd)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to set load config: status %s", resp.Status)
	}
	return nil
}

// SetLightpadConfig
func (l *DefaultLightpad) SetLightpadConfig(conf LightpadConfig) error {
	pd := struct {
		Config LightpadConfig `json:"config"`
		LLID   string         `json:"llid"`
	}{conf, l.LLID}
	resp, err := l.makePadPOSTRequest(pathSetLogicalLoadConfig, pd)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to set pad config level: status %s", resp.Status)
	}
	return nil

}

func (l *DefaultLightpad) GetLogicalLoadMetrics() (LogicalLoadMetrics, error) {
	pd := struct {
		LLID string `json:"llid"`
	}{
		LLID: l.LLID,
	}
	resp, err := l.makePadPOSTRequest(pathGetLogicalLoadMetrics, pd)
	if err != nil {
		return LogicalLoadMetrics{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return LogicalLoadMetrics{}, fmt.Errorf("failed to get load metrics level: status %s", resp.Status)
	}
	lm := LogicalLoadMetrics{}
	err = json.NewDecoder(resp.Body).Decode(&lm)
	if err != nil {
		return LogicalLoadMetrics{}, err
	}
	return lm, nil
}

func (l *DefaultLightpad) SetLogicalLoadGlow(glow ForceGlow) error {
	resp, err := l.makePadPOSTRequest(pathSetLogicalLoadGlow, glow)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to set load glow: status %s", resp.Status)
	}
	return nil
}

func (l *DefaultLightpad) makePadPOSTRequest(urlPath string, postData interface{}) (*http.Response, error) {
	userAgent := fmt.Sprintf("%s/%s", DefaultUserAgent, Version)
	if UserAgentAddition != "" {
		userAgent = fmt.Sprintf("%s %s", userAgent, strings.TrimSpace(UserAgentAddition))
	}
	api, err := url.Parse(fmt.Sprintf("https://%s:%d", l.IP, l.Port))
	if err != nil {
		return nil, err
	}
	api.Path = path.Join(api.Path, urlPath)
	pd, err := json.Marshal(postData)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", api.String(), bytes.NewReader(pd))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	// sha256sum the HAT for the header
	encHat := fmt.Sprintf("%x", sha256.Sum256([]byte(l.HAT)))
	req.Header.Set("X-Plum-House-Access-Token", encHat)
	if l.HttpClient == nil {
		l.HttpClient = &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
	}
	return l.HttpClient.Do(req)
}

// Subscribe returns a channel that will send you state changes from the
// lightpad. When you want to close the connection and are done listening for
// state changes on the lightpad, cancel the context passed in and it will clean
// itself up.
func (l *DefaultLightpad) Subscribe(ctx context.Context) error {
	logrus.WithField("ipaddr", l.IP).Debug("about to connect to lightpad")
	if l.StateChanges == nil {
		l.StateChanges = make(chan Event, 5)
	}
	addrStr := fmt.Sprintf("%s:%d", l.IP, 2708)
	conn, err := net.Dial("tcp", addrStr)
	if err != nil {
		logrus.WithField("error", err).Debug("failed to connect to lightpad")
		return err
	}
	go func() {
		for {
			// fmt.Println("beg loop")
			select {
			case <-ctx.Done():
				logrus.WithField("lpid", l.ID).Debug("we've been cancelled")
				return // we've been cancelled
			default:
			}
			// fmt.Println("in loop")
			message, err := bufio.NewReader(conn).ReadString('\n')
			// fmt.Println("got mess")
			if err != nil {
				fmt.Printf("got error from subscribed conn %s\n", err)
				continue
			}
			// fmt.Println(message)
			message = strings.TrimSuffix(strings.TrimSpace(message), ".")
			lpe := lightpadEvent{}
			err = json.Unmarshal([]byte(message), &lpe)
			if err != nil {
				l.StateChanges <- lightpadEvent{Error: err}
			}
			switch lpe.Type {
			case "dimmerchange":
				ev := LPEDimmerChange{}
				err = json.Unmarshal([]byte(message), &ev)
				if err != nil {
					l.StateChanges <- lightpadEvent{Error: err}
				}
				l.StateChanges <- ev
			case "power":
				ev := LPEPower{}
				err = json.Unmarshal([]byte(message), &ev)
				if err != nil {
					l.StateChanges <- lightpadEvent{Error: err}
				}
				l.StateChanges <- ev
			case "pirSignal":
				ev := LPEPIRSignal{}
				err = json.Unmarshal([]byte(message), &ev)
				if err != nil {
					l.StateChanges <- lightpadEvent{Error: err}
				}
				l.StateChanges <- ev
			default:
				l.StateChanges <- LPEUnknown{
					lightpadEvent{
						Type: "unknown",
					},
					message,
				}

			}
		}
	}()
	return nil
}
