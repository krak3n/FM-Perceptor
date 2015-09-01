// Publish Volume Events to Redis

package events

import (
	"encoding/json"
	"strconv"

	"gopkg.in/redis.v3"
)

// Play POST request expected JSON payload
type publishVolumePayload struct {
	event string `json:"event"`
	level int    `json:"volume"`
}

// Publish a Play Event from the Player to Redis. This sets the current
// playing track, the user and start time as well as publishing to the
// event channel
func PublishVolumeEvent(c *redis.Client, level int) error {
	var err error

	if err = c.Set(volumeKey, strconv.Itoa(level), 0).Err(); err != nil {
		return err
	}

	// Generate message payload
	payload, err := json.Marshal(&publishVolumePayload{
		event: volumeEvent,
		level: level,
	})
	if err != nil {
		return err
	}

	// Publish Message
	if err = c.Publish(eventsChannel, string(payload[:])).Err(); err != nil {
		return err
	}

	// Also Publish mute change to false
	if err = PublishMuteEvent(c, false); err != nil {
		return err
	}

	return nil
}
