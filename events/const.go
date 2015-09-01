// Events Package Constants

package events

// Redis Keys
const (
	eventsChannel     string = "fm:events"
	currentTrackKey   string = "fm:player:current"
	startTimeKey      string = "fm:player:start_time"
	pauseKey          string = "fm:player:paused"
	pauseTimeKey      string = "fm:player:pause_time"
	pauseDurrationKey string = "fm:player:pause_duration"
	volumeKey         string = "fm:player:volume"
	muteKey           string = "fm:player:mute"
)

// Events
const (
	playEvent   string = "play"
	endEvent    string = "end"
	volumeEvent string = "volume_changed"
	muteEvent   string = "mute_changed"
)

// Event payload to publish
type publishEventPayload struct {
	event string `json:"event"`
	track string `json:"uri"`
	user  string `json:"user"`
}
