package libplumraw

import (
	"context"
	"net/http"
	"time"
)

const (
	Version = "0.0.1"

	// DefaultLightpadPort is the port on which lightpads listen for HTTP requests
	DefaultLightpadPort = 8443
	//DefaultLightpadHeartbeatPort the lightpads use to broadcast UDP status.
	//Lightpads send out a heartbeat once every ~5 minutes.
	DefaultLightpadHeartbeatPort = 43770
	DefaultUserAgent             = "libplumraw"
	DefaultPlumAPIHOST           = "https://production.plum.technology"

	// website API paths
	pathGetHouses      = "/v2/getHouses"
	pathGetHouse       = "/v2/getHouse"
	pathGetScenes      = "/v2/getScenes"
	pathGetScene       = "/v2/getScene"
	pathGetRoom        = "/v2/getRoom"
	pathGetLogicalLoad = "/v2/getLogicalLoad"
	pathGetLightpad    = "/v2/getLightpad"

	// lightpad API paths
	pathSetLogicalLoadLevel   = "/v2/setLogicalLoadLevel"
	pathSetLogicalLoadConfig  = "/v2/setLogicalLoadConfig"
	pathGetLogicalLoadMetrics = "/v2/getLogicalLoadMetrics"
	pathSetLogicalLoadGlow    = "/v2/setLogicalLoadGlow"
)

var UserAgentAddition string

// Use this for lightpads but not web connections
// var lightpadHttpClient = &http.Client{Transport: &http.Transport{
// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// }}

type WebConnection interface {
	GetHouses() (Houses, error)
	GetHouse(string) (House, error)
	GetScenes(string) (Scenes, error)
	GetScene(string) (Scene, error)
	GetRoom(string) (Room, error)
	GetLogicalLoad(string) (LogicalLoad, error)
	GetLightpad(string) (LightpadSpec, error)
}

type Lightpad interface {
	SetLogicalLoadLevel(level int) error
	SetLogicalLoadConfig(conf LogicalLoadConfig) error
	GetLogicalLoadMetrics() (LogicalLoadMetrics, error)
	SetLogicalLoadGlow(glow ForceGlow) error
	Subscribe(context.Context) (chan Event, error)
}

type WebConnectionConfig struct {
	Email      string
	Password   string
	PlumAPIURL string // default https://production.plum.technology/
}

type defaultWebConnection struct {
	config     WebConnectionConfig
	HttpClient *http.Client
}

//
func NewWebConnection(conf WebConnectionConfig) WebConnection {
	c := &defaultWebConnection{
		config: conf,
	}
	c.HttpClient = &http.Client{}
	if c.config.PlumAPIURL == "" {
		c.config.PlumAPIURL = "https://production.plum.technology/"
	}
	return c
}

// TestWebConnection implements the WebConnection interface and is for using
// this library in upstream tests - instead of calling out to any URL it just
// returns the error or objects with which it was configured
type TestWebConnection struct {
	Houses       Houses
	House        House
	Scenes       Scenes
	Scene        Scene
	Room         Room
	LogicalLoad  LogicalLoad
	LightpadSpec LightpadSpec
	Error        *error
}

func NewTestWebConnection() *TestWebConnection {
	return &TestWebConnection{}
}

func (t *TestWebConnection) GetHouses() (Houses, error) {
	if t.Error != nil {
		return nil, *t.Error
	}
	return t.Houses, nil
}

func (t *TestWebConnection) GetHouse(hid string) (House, error) {
	if t.Error != nil {
		return House{}, *t.Error
	}
	return t.House, nil
}

func (t *TestWebConnection) GetScenes(hid string) (Scenes, error) {
	if t.Error != nil {
		return nil, *t.Error
	}
	return t.Scenes, nil
}

func (t *TestWebConnection) GetScene(hid string) (Scene, error) {
	if t.Error != nil {
		return Scene{}, *t.Error
	}
	return t.Scene, nil
}

func (t *TestWebConnection) GetRoom(rid string) (Room, error) {
	if t.Error != nil {
		return Room{}, *t.Error
	}
	return t.Room, nil
}

func (t *TestWebConnection) GetLogicalLoad(llid string) (LogicalLoad, error) {
	if t.Error != nil {
		return LogicalLoad{}, *t.Error
	}
	return t.LogicalLoad, nil
}

func (t *TestWebConnection) GetLightpad(lpid string) (LightpadSpec, error) {
	if t.Error != nil {
		return LightpadSpec{}, *t.Error
	}
	return t.LightpadSpec, nil
}

// TestLightpad implements the Lightpad interface and is for using this library
// in upstream tests - instead of calling out to an actual lightpad it just
// returns the error or objects with which it was configured
type TestLightpad struct {
	LogicalLoadMetrics LogicalLoadMetrics
	Error              *error
}

func (t *TestLightpad) SetLogicalLoadLevel(level int) error {
	if t.Error != nil {
		return *t.Error
	}
	return nil
}
func (t *TestLightpad) SetLogicalLoadConfig(conf LogicalLoadConfig) error {
	if t.Error != nil {
		return *t.Error
	}
	return nil
}
func (t *TestLightpad) GetLogicalLoadMetrics() (LogicalLoadMetrics, error) {
	if t.Error != nil {
		return LogicalLoadMetrics{}, *t.Error
	}
	return t.LogicalLoadMetrics, nil
}
func (t *TestLightpad) SetLogicalLoadGlow(glow ForceGlow) error {
	if t.Error != nil {
		return *t.Error
	}
	return nil
}

// TestLightpadHeartbeat sends a lightpad announcement every 2 seconds until the
// context used to intitialize it is cancelled.
type TestLightpadHeartbeat struct {
	LightpadAnnouncement
}

func (t *TestLightpadHeartbeat) Listen(ctx context.Context) chan LightpadAnnouncement {
	responses := make(chan LightpadAnnouncement, 0)
	tick := time.NewTicker(2 * time.Second).C
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			responses <- t.LightpadAnnouncement
		}
	}()
	return responses
}
