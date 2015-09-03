// Handle Mute Events

package events

import (
	"encoding/json"

	"gopkg.in/redis.v3"
)

// Play POST request expected JSON payload
type publishMutePayload struct {
	Event  string `json:"event"`
	Active bool   `json:"mute"`
}

// Publish a Play Event from the Player to Redis. This sets the current
// playing track, the user and start time as well as publishing to the
// event channel
func PublishMuteEvent(c *redis.Client, active bool) error {
	var err error

	var state string
	if active {
		state = "1"
	} else {
		state = "0"
	}

	// Set mute state on Redis
	if err = c.Set(muteKey, state, 0).Err(); err != nil {
		return err
	}

	// Generate message payload
	payload, err := json.Marshal(&publishMutePayload{
		Event:  muteEvent,
		Active: active,
	})
	if err != nil {
		return err
	}

	// Publish Message
	err = c.Publish(eventsChannel, string(payload[:])).Err()
	if err != nil {
		return err
	}

	return nil
}
