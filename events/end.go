// Publish End Events to Redis

package events

import (
	"encoding/json"

	"gopkg.in/redis.v3"
)

// Publish track end event to Redis, this will delete current track,
// start time keys as well as publish the end event
func PublishEndEvent(c *redis.Client, track string, user string) error {
	var err error

	// Delete Current Track
	err = c.Del(currentTrackKey).Err()
	if err != nil {
		return err
	}

	// Delete Start Time
	err = c.Del(startTimeKey).Err()
	if err != nil {
		return err
	}

	// Generate message payload
	payload, err := json.Marshal(&publishEventPayload{
		event: "end",
		track: track,
		user:  user,
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
