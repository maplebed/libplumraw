package libplumraw

// web.go has all the functions that call out to the Plum web service. These are
// used to fetch house, room, load, and lightpad configurations.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
)

func (c *defaultWebConnection) GetHouses() (Houses, error) {
	resp, err := c.makePlumWebGETRequest(pathGetHouses)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to reach Plum: status %s", resp.Status)
	}
	hids := make(Houses, 0, 0)
	// spew.Dump(resp)
	json.NewDecoder(resp.Body).Decode(&hids)
	return hids, nil
}

func (c *defaultWebConnection) GetHouse(hid string) (House, error) {
	postData := struct {
		HID string `json:"hid"`
	}{hid}
	resp, err := c.makePlumWebPOSTRequest(pathGetHouse, postData)
	if err != nil {
		return House{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return House{}, fmt.Errorf("failed to reach Plum: status %s", resp.Status)
	}
	house := House{}
	err = json.NewDecoder(resp.Body).Decode(&house)
	if err != nil {
		return House{}, err
	}
	sort.Strings(house.RoomIDs)
	return house, nil
}

func (c *defaultWebConnection) GetScenes(hid string) (Scenes, error) {
	postData := struct {
		HID string `json:"hid"`
	}{hid}
	resp, err := c.makePlumWebPOSTRequest(pathGetScenes, postData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to reach Plum: status %s", resp.Status)
	}
	sids := make(Scenes, 0, 0)
	json.NewDecoder(resp.Body).Decode(&sids)
	return sids, nil
}

func (c *defaultWebConnection) GetScene(sid string) (Scene, error) {
	postData := struct {
		SID string `json:"sid"`
	}{sid}
	resp, err := c.makePlumWebPOSTRequest(pathGetScene, postData)
	if err != nil {
		return Scene{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Scene{}, fmt.Errorf("failed to reach Plum: status %s", resp.Status)
	}
	scene := Scene{}
	err = json.NewDecoder(resp.Body).Decode(&scene)
	if err != nil {
		return Scene{}, err
	}
	return scene, nil
}

func (c *defaultWebConnection) GetRoom(rid string) (Room, error) {
	postData := struct {
		RID string `json:"rid"`
	}{rid}
	resp, err := c.makePlumWebPOSTRequest(pathGetRoom, postData)
	if err != nil {
		return Room{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Room{}, fmt.Errorf("failed to reach Plum: status %s", resp.Status)
	}
	room := Room{}
	err = json.NewDecoder(resp.Body).Decode(&room)
	if err != nil {
		return Room{}, err
	}
	sort.Strings(room.LLIDs)
	return room, nil
}

func (c *defaultWebConnection) GetLogicalLoad(llid string) (LogicalLoad, error) {
	postData := struct {
		LLID string `json:"llid"`
	}{llid}
	resp, err := c.makePlumWebPOSTRequest(pathGetLogicalLoad, postData)
	if err != nil {
		return LogicalLoad{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return LogicalLoad{}, fmt.Errorf("failed to reach Plum: status %s", resp.Status)
	}
	ll := LogicalLoad{}
	err = json.NewDecoder(resp.Body).Decode(&ll)
	if err != nil {
		return LogicalLoad{}, err
	}
	sort.Strings(ll.LPIDs)
	return ll, nil
}

func (c *defaultWebConnection) GetLightpad(lpid string) (LightpadSpec, error) {
	postData := struct {
		LPID string `json:"lpid"`
	}{lpid}
	resp, err := c.makePlumWebPOSTRequest(pathGetLightpad, postData)
	if err != nil {
		return LightpadSpec{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return LightpadSpec{}, fmt.Errorf("failed to reach Plum: status %s", resp.Status)
	}
	lp := LightpadSpec{}
	err = json.NewDecoder(resp.Body).Decode(&lp)
	if err != nil {
		return LightpadSpec{}, err
	}
	lp.ID = lpid
	return lp, nil
}

func (c *defaultWebConnection) makePlumWebGETRequest(urlPath string) (*http.Response, error) {
	userAgent := fmt.Sprintf("%s/%s", DefaultUserAgent, Version)
	if UserAgentAddition != "" {
		userAgent = fmt.Sprintf("%s %s", userAgent, strings.TrimSpace(UserAgentAddition))
	}
	api, err := url.Parse(c.config.PlumAPIURL)
	if err != nil {
		return nil, err
	}
	api.Path = path.Join(api.Path, urlPath)
	req, err := http.NewRequest("GET", api.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	// spew.Dump(c.config)
	req.SetBasicAuth(c.config.Email, c.config.Password)
	// spew.Dump(req)
	return c.HttpClient.Do(req)
}

func (c *defaultWebConnection) makePlumWebPOSTRequest(urlPath string, postData interface{}) (*http.Response, error) {
	userAgent := fmt.Sprintf("%s/%s", DefaultUserAgent, Version)
	if UserAgentAddition != "" {
		userAgent = fmt.Sprintf("%s %s", userAgent, strings.TrimSpace(UserAgentAddition))
	}
	api, err := url.Parse(c.config.PlumAPIURL)
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
	req.SetBasicAuth(c.config.Email, c.config.Password)
	return c.HttpClient.Do(req)
}
