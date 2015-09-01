// Pause Event Handling

package events

import "gopkg.in/redis.v3"

// Sets the correct redis keys when the player enters a pause state
func SetPauseState(c *redis.Client, start string) error {
	// Create Transaction
	tx := c.Multi()

	// Execute Transaction
	for {
		_, err := tx.Exec(func() error {
			tx.Set(pauseKey, "1", 0)
			tx.Set(pauseTimeKey, start, 0)
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
