// Publish Play Events to Redis

package events

import (
	"encoding/json"

	"gopkg.in/redis.v3"
)

// Represents the publish event json payload
type publishEventPayload struct {
	event string `json:"event"`
	track string `json:"uri"`
	user  string `json:"user"`
}

// Publish a Play Event from the Player to Redis. This sets the current
// playing track, the user and start time as well as publishing to the
// event channel
func PublishPlayEvent(c *redis.Client, track string, user string, start string) error {
	var err error

	// Save Current Track
	err = c.Set("fm:player:current", string(track), 0).Err()
	if err != nil {
		return err
	}

	// Save Start Time
	err = c.Set("fm:player:start_time", string(start), 0).Err()
	if err != nil {
		return err
	}

	// Generate message payload
	payload, err := json.Marshal(&publishEventPayload{
		event: "play",
		track: track,
		user:  user,
	})
	if err != nil {
		return err
	}

	// Publish Message
	err = c.Publish("fm:events", string(payload[:])).Err()
	if err != nil {
		return err
	}

	return nil
}
