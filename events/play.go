// Publish Play Events to Redis

package events

import (
	"encoding/json"

	"gopkg.in/redis.v3"
)

// Publish a Play Event from the Player to Redis. This sets the current
// playing track, the user and start time as well as publishing to the
// event channel
func PublishPlayEvent(c *redis.Client, track string, user string, start string) error {
	var err error

	// Save Current Track
	err = c.Set(currentTrackKey, string(track), 0).Err()
	if err != nil {
		return err
	}

	// Save Start Time
	err = c.Set(startTimeKey, string(start), 0).Err()
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
	err = c.Publish(eventsChannel, string(payload[:])).Err()
	if err != nil {
		return err
	}

	return nil
}
