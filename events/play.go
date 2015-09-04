// Publish Play Events to Redis

package events

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v3"
)

// Publish a Play Event from the Player to Redis. This sets the current
// playing track, the user and start time as well as publishing to the
// event channel
func PublishPlayEvent(c *redis.Client, track string, user string, start string) error {
	var err error

	// Create content for current key
	data, err := json.Marshal(&currentTrack{
		Track: track,
		User:  user,
	})
	if err != nil {
		log.Errorf("Failed to marshal current key: %s", err)
	}

	// Create Transaction
	tx := c.Multi()

	// Execute Transaction
	for {
		_, err := tx.Exec(func() error {
			tx.Set(currentTrackKey, string(data[:]), 0)
			tx.Set(startTimeKey, start, 0)
			tx.Set(pauseKey, "0", 0)
			tx.Del(pauseTimeKey)
			tx.Del(pauseDurrationKey)
			return nil
		})
		if err == redis.TxFailedErr {
			// Retry.
			continue
		} else if err != nil {
			return err
		}
		break
	}

	// Generate message payload
	payload, err := json.Marshal(&publishEventPayload{
		Event: playEvent,
		Track: track,
		User:  user,
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
