package libplumraw

// types.go is the main types
import (
	"net"
	"net/http"
)

// Houses is a list of House IDs
type Houses IDs

type House struct {
	ID       string `json:"hid,omitempty"`      // UUID
	RoomIDs  IDs    `json:"rids,omitempty"`     // list of Room UUIDs
	Location string `json:"location,omitempty"` // zip code
	LatLong  struct {
		Latitude  float64 `json:"latitude_degrees_north,omitempty"` // decimal degrees North
		Longitude float64 `json:"longitude_degrees_west,omitempty"` // decimal degrees West
	}
	AccessToken string `json:"house_access_token,omitempty"`
	Name        string `json:"house_name,omitempty"`
	// TimeZone is seconds offset from UTC for the local time zone
	TimeZone int `json:"local_tz,omitempty"`
}

// Scenes is a list of Scene IDs
type Scenes IDs

type Scene struct {
	ID       string          `json:"sid"`
	Settings []SceneSettings `json:"settings"`
	HouseID  string          `json:"hid"`
	Name     string          `json:"scene_name"`
}

type SceneSettings struct {
	LLID  string `json:"llid"`
	Level int    `json:"level"` // range 0-255
	Fade  int    `json:"fade"`  // transition time, in milliseconds
}

type Room struct {
	// ID is the UUID identifying this room
	ID   string `json:"rid,omitempty"`
	Name string `json:"room_name,omitempty"`
	// House is the house in which this room exists
	HouseID string `json:"hid,omitempty"`
	// All the lightpads in a room. Note a lightpad being in a room means that
	// the load (the lights) are in the room, not necessarily the switch itself.
	LLIDs IDs `json:"llids,omitempty"`
}

type LogicalLoad struct {
	ID     string `json:"llid,omitempty"`
	Name   string `json:"logical_load_name,omitempty"`
	LPIDs  IDs    // LPIDs is a list of lightpad IDs in this logical load
	RoomID string `json:"rid,omitempty"`
}

type LogicalLoadConfig struct {
	GlowColor   LightpadGlowColor `json:"glowColor,omitempty"`
	GlowTimeout int               `json:"glowTimeout,omitempty"`
	GlowEnabled bool              `json:"glowEnabled"`
}

// LightpadSpec represents an individual light switch, as reported by the web
// API
type LightpadSpec struct {
	// ID is the Lightpad ID identifying this switch. Stringified version of a
	// UUID
	ID string `json:"lpid,omitempty"`
	// LLID is the Logical Load ID, and can be shared by multiple switches.
	// Stringified version of a UUID
	LLID string `json:"llid"`
	// LogicalLoad is a pionter to the actual logical load struct matching the
	// LLID
	Config         LightpadConfig `json:"config,omitempty"`
	IsProvisioned  bool           `json:"is_provisioned,omitempty"`  // Whether lightpad is provisioned to a house/room/logical load
	CustomGestures int            `json:"custom_gestures,omitempty"` // Unused
	Name           string         `json:"lightpad_name,omitempty"`   // Lightpad Name
}

// Lightpad represents an actual switch to which you can make calls
type DefaultLightpad struct {
	ID         string       `json:"lpid,omitempty"`
	LLID       string       `json:"llid"`
	Level      int          `json:"level,omitempty"` // range 0-255
	Power      int          `json:"power,omitempty"`
	IP         net.IP       `json:"ip"`   // IP address of this lightpad
	Port       int          `json:"port"` // port on which this lightpad listens
	HAT        string       `json:"hat"`  // house access token
	HttpClient *http.Client `json:"-"`

	// StateChanges is a channel down which the lightpad will send state change
	// events
	StateChanges chan Event `json:"-"`
}

type LightpadConfig struct {
	AMQPEnabled          bool              `json:"amqpEnabled,omitempty"`  // always true
	Cluster              string            `json:"cluster,omitempty"`      // production or development cluster
	DefaultLevel         int               `json:"defaultLevel,omitempty"` // range 0-255 default power level
	DimEnabled           bool              `json:"dimEnabled,omitempty"`   // true if switch is a dimmer, false for ON/OFF only
	FadeOffTime          int               `json:"fadeOffTime,omitempty"`  // milliseconds to fade from on to off
	FadeOnTime           int               `json:"fadeOnTime,omitempty"`   // milliseconsd to fade from off to on
	ForceGlow            bool              `json:"forceGlow,omitempty"`    // bool is glow currently forced
	GlowColor            LightpadGlowColor `json:"glowColor,omitempty"`
	GlowEnabled          bool              `json:"glowEnabled,omitempty"`          // turn on glow when PIR detects motion
	GlowFade             int               `json:"glowFade,omitempty"`             // millisec to fade off the glow ring
	GlowIntensity        float64           `json:"glowIntensity,omitempty"`        // range 0-1 glow brightness
	GlowTimeout          int               `json:"glowTimeout,omitempty"`          // seconds glow remains on for motion
	GlowTracksDimmer     bool              `json:"glowTracksDimmer,omitempty"`     // glow same as light level
	LogRemote            bool              `json:"logRemote,omitempty"`            // always false
	MaxWattage           int               `json:"maxWattage,omitempty"`           // 420
	MinimumLevel         int               `json:"minimumLevel,omitempty"`         // range 0-255 minimum dimmable power level
	Name                 string            `json:"name,omitempty"`                 // unused
	OccupancyAction      string            `json:"occupancyAction,omitempty"`      // unused
	OccupancyTimeout     int               `json:"occupancyTimeout,omitempty"`     // seconds before no occupancy
	PIRSensitivity       int               `json:"pirSensitivity,omitempty"`       // range 0-255 sensitivity of motion sensor 0 to ~6ft
	RememberLastDimLevel bool              `json:"rememberLastDimLevel,omitempty"` // return to last dim level
	SerialNumber         string            `json:"serialNumber,omitempty"`         // serial number
	SlowFadeTime         int               `json:"slowFadeTime,omitempty"`         // milliseconds
	TouchRate            float64           `json:"touchRate,omitempty"`            // range 0-1 touch sensitivity
	TrackingSpeed        int               `json:"trackingSpeed,omitempty"`        // 1000
	UUID                 string            `json:"uuid,omitempty"`                 // not sure what this UUID is used for
	VersionLocked        bool              `json:"versionLocked,omitempty"`        // true if this switch is software upgradable
}

// LightpadGlowColor indicates the color of the glow ring and its brightness.
// Each color has the range 0-255
type LightpadGlowColor struct {
	White int `json:"white,omitempty"`
	Red   int `json:"red,omitempty"`
	Green int `json:"green,omitempty"`
	Blue  int `json:"blue,omitempty"`
}

// ForceGlow is used to temporarily force a speciifc glow
type ForceGlow struct {
	LightpadGlowColor
	// Intensity is range 0-1 - Brightness Level Of Glow Ring
	Intensity float64 `json:"intensity"`
	// Timeout is int milliseconds to remain on
	Timeout int `json:"timeout"`
	// LLID is the Logical Load ID to identify the switch
	LLID string `json:"llid"`
}

// LogicalLoadMetrics is the response from the lightpad to a request for metrics
// and covers a logical load and all of the switches that share that load
type LogicalLoadMetrics struct {
	// Level is the current power level of the logical load
	Level int `json:"level,omitempty"`
	// Power is the current wattage of the logical load
	Power int `json:"power,omitempty"`
	// Lightpads is a list of the power level for each switch
	Metrics []LightpadMetric `json:"lightpad_metrics,omitempty"`
}

// LightpadMetric is an individual switch's load
type LightpadMetric struct {
	// ID is identifying a specific switch
	ID string `json:"lpid,omitempty"`
	// Level is the current power level of the switch
	Level int `json:"level,omitempty"`
	// Power is the current wattage of the switch
	Power int `json:"power,omitempty"`
}

// TODO give lightpad events more structure instead of just a map[string]interface{}
type LightpadEventType uint32

const (
	UndefEvent LightpadEventType = iota
	DimmerChange
	Power
	PIRSignal
	ConfigChange
)

// LightpadEvent is what's emitted from the lightpad when it changes state.
// Examples: load change, PIR activated, etc.
type Event interface{}

type lightpadEvent struct {
	Type  string
	Error error
}

type LPEUnknown struct {
	lightpadEvent
	Message string
}
type LPEDimmerChange struct {
	lightpadEvent
	Level int
}
type LPEPower struct {
	lightpadEvent
	Watts int
}
type LPEPIRSignal struct {
	lightpadEvent
	Signal int
}
type LPEConfigChange struct {
	lightpadEvent
}

type IDs []string

// make a list of IDs sortable and comparable
func (ids IDs) Len() int           { return len(ids) }
func (ids IDs) Swap(i, j int)      { ids[i], ids[j] = ids[j], ids[i] }
func (ids IDs) Less(i, j int) bool { return ids[i] < ids[j] }
func (ids IDs) Equals(other IDs) bool {
	if len(ids) != len(other) {
		return false
	}
	for j, id := range ids {
		if id != other[j] {
			return false
		}
	}
	return true
}
