// Resume Event Handling

package events

import "gopkg.in/redis.v3"

// Sets the correct redis keys when the player enters a pause state
func SetResumeState(c *redis.Client, durration string) error {
	// Create Transaction
	tx := c.Multi()

	// Execute Transaction
	for {
		_, err := tx.Exec(func() error {
			tx.Set(pauseKey, "0", 0)
			tx.Set(pauseDurrationKey, durration, 0)
			tx.Del(pauseTimeKey)
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

	return nil
}
