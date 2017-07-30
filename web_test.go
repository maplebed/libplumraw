package libplumraw

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHouses(t *testing.T) {
	// expected response
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprintln(w, `["houseid1", "houseid2"]`)
	})
	wc := newMockHTTP(hf)
	hs, err := wc.GetHouses()
	assert.NoError(t, err)
	expect := Houses{"houseid1", "houseid2"}
	assert.Equal(t, expect, hs)

	// permission denied
	hf = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	})
	wc = newMockHTTP(hf)
	_, err = wc.GetHouses()
	assert.Error(t, err)
}

func TestGetHouse(t *testing.T) {
	// what I expect to send
	expCallStr := `{"hid":"sending-houseid1"}`
	// what the server will return
	respStr := `{"rids":["roomid1"],"location":"012345","hid":"rcv-houseid1","latlong":{"latitude_degrees_north":34.567,"longitude_degrees_west":123.456},"house_access_token":"bonnie-cap","house_name":"sarah","local_tz":-25200}`
	// and the house that that resp creates
	expHouse := House{
		ID:       "rcv-houseid1",
		RoomIDs:  IDs{"roomid1"},
		Location: "012345",
		LatLong: struct {
			Latitude  float64 `json:"latitude_degrees_north,omitempty"`
			Longitude float64 `json:"longitude_degrees_west,omitempty"`
		}{34.567, 123.456},
		AccessToken: "bonnie-cap",
		Name:        "sarah",
		TimeZone:    -25200,
	}
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		bod, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, expCallStr, string(bod))
		fmt.Fprintln(w, respStr)
	})
	wc := newMockHTTP(hf)
	hs, err := wc.GetHouse("sending-houseid1")
	assert.NoError(t, err)
	assert.Equal(t, expHouse, hs)

	// permission denied
	hf = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	})
	wc = newMockHTTP(hf)
	_, err = wc.GetHouse("sending-houseid2")
	assert.Error(t, err)
}

func TestGetScenes(t *testing.T) {
	// what I expect to send
	expCallStr := `{"hid":"sending-houseid1"}`
	// expected response
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		bod, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, expCallStr, string(bod))
		fmt.Fprintln(w, `["sceneid1", "sceneid2"]`)
	})
	wc := newMockHTTP(hf)
	sc, err := wc.GetScenes("sending-houseid1")
	assert.NoError(t, err)
	expect := Scenes{"sceneid1", "sceneid2"}
	assert.Equal(t, expect, sc)

	// permission denied
	hf = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	})
	wc = newMockHTTP(hf)
	_, err = wc.GetScenes("sending-houseid1")
	assert.Error(t, err)
}

func TestGetScene(t *testing.T) {
	// what I expect to send
	expCallStr := `{"sid":"sceneid1"}`
	// what the server will return
	respStr := `{"settings":[{"llid":"load-id1","level":70,"fade":10000}],"hid":"rcv-houseid1","scene_name":"pastoral","sid":"sceneid1"}`
	// and the scene that that resp creates
	expScene := Scene{
		ID:       "sceneid1",
		Settings: []SceneSettings{{"load-id1", 70, 10000}},
		HouseID:  "rcv-houseid1",
		Name:     "pastoral",
	}
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		bod, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, expCallStr, string(bod))
		fmt.Fprintln(w, respStr)
	})
	wc := newMockHTTP(hf)
	sc, err := wc.GetScene("sceneid1")
	assert.NoError(t, err)
	assert.Equal(t, expScene, sc)

	// permission denied
	hf = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	})
	wc = newMockHTTP(hf)
	_, err = wc.GetScene("sceneid2")
	assert.Error(t, err)
}

func TestGetRoom(t *testing.T) {
	// what I expect to send
	expCallStr := `{"rid":"roomid1"}`
	// what the server will return
	respStr := `{"rid":"roomid1","hid":"houseid1","llids":["loadid1"],"room_name":"dungeon"}`
	// and the room that that resp creates
	expRoom := Room{
		ID:      "roomid1",
		Name:    "dungeon",
		HouseID: "houseid1",
		LLIDs:   IDs{"loadid1"},
	}
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		bod, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, expCallStr, string(bod))
		fmt.Fprintln(w, respStr)
	})
	wc := newMockHTTP(hf)
	room, err := wc.GetRoom("roomid1")
	assert.NoError(t, err)
	assert.Equal(t, expRoom, room)

	// permission denied
	hf = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	})
	wc = newMockHTTP(hf)
	_, err = wc.GetRoom("roomid1")
	assert.Error(t, err)
}

func TestGetLogicalLoad(t *testing.T) {
	// what I expect to send
	expCallStr := `{"llid":"loadid1"}`
	// what the server will return
	respStr := `{"rid":"roomid1","lpids":["padid1"],"logical_load_name":"chair","llid":"loadid1"}`
	// and the load that that resp creates
	expLoad := LogicalLoad{
		ID:     "loadid1",
		Name:   "chair",
		RoomID: "roomid1",
		LPIDs:  IDs{"padid1"},
	}
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		bod, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, expCallStr, string(bod))
		fmt.Fprintln(w, respStr)
	})
	wc := newMockHTTP(hf)
	load, err := wc.GetLogicalLoad("loadid1")
	assert.NoError(t, err)
	assert.Equal(t, expLoad, load)

	// permission denied
	hf = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	})
	wc = newMockHTTP(hf)
	_, err = wc.GetLogicalLoad("loadid1")
	assert.Error(t, err)
}

func TestGetLightpad(t *testing.T) {
	// what I expect to send
	expCallStr := `{"lpid":"pad-id"}`
	// what the server will return
	respStr := `{
  "config": {
    "amqpEnabled": true,
    "cluster": "production",
    "defaultLevel": 255,
    "dimEnabled": true,
    "fadeOffTime": 0,
    "fadeOnTime": 0,
    "forceGlow": false,
    "glowColor": {
      "green": 12,
      "red": 29
    },
    "glowEnabled": true,
    "glowFade": 1000,
    "glowIntensity": 0.4,
    "glowTimeout": 10,
    "glowTracksDimmer": false,
    "logRemote": false,
    "maxWattage": 420,
    "minimumLevel": 51,
    "name": "",
    "occupancyAction": "on",
    "occupancyTimeout": 5,
    "pirSensitivity": 192,
    "rememberLastDimLevel": false,
    "serialNumber": "BL012345ABCDE",
    "slowFadeTime": 15000,
    "touchRate": 0.6654275059700012,
    "trackingSpeed": 1000,
    "uuid": "pad-uuid",
    "versionLocked": false
  },
  "custom_gestures": 0,
  "is_provisioned": true,
  "lightpad_name": "pad-uuid",
  "llid": "load-uuid",
  "lpid": "pad-uuid"
}`
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		bod, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, expCallStr, string(bod))
		fmt.Fprintln(w, respStr)
	})
	wc := newMockHTTP(hf)
	_, err := wc.GetLightpad("pad-id")
	assert.NoError(t, err)
	// assert.Equal(t, expHouse, pad)

	// permission denied
	hf = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	})
	wc = newMockHTTP(hf)
	_, err = wc.GetLightpad("pad-id")
	assert.Error(t, err)
}

func newMockHTTP(handler http.Handler) WebConnection {
	ts := httptest.NewServer(handler)
	conf := WebConnectionConfig{
		PlumAPIURL: ts.URL,
	}
	return NewWebConnection(conf)
}
