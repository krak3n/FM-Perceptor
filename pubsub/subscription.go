// Temporary Redis Event Subscription to pass events from
// API's to the Player clients

package pubsub

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/thisissoon/FM-Perceptor/socket"
	"gopkg.in/redis.v3"
)

// Redis event subscription
type Subscription struct {
	hub     socket.Hub
	client  *redis.Client
	channel string
}

// Consume events from redis, passing them onto the WS hub for
// broadcast to clients
func (s *Subscription) Consume() {
	for {
		client := s.client
		// Connect to Redis Pubsub Channel
		pubsub := client.PubSub()
		err := pubsub.Subscribe(s.channel)
		// On error sleep for 1 second, log and continue to the next loop iteration
		if err != nil {
			log.Errorf("Redis Connection Error: %s", err)
			time.Sleep(time.Second)
			continue
		}
	ReceiveLoop: // Inner loop label
		for {
			msg, err := pubsub.Receive() // recieve a message from the channel
			if err != nil {
				// On error, close the channel and break out of the loop inner loop
				log.Error("Redis Error: %s", err)
				pubsub.Close()
				break ReceiveLoop
			} else {
				switch m := msg.(type) { // Switch the mesage type
				case *redis.Subscription:
					log.Debugf("Subscribed: %s", m.Channel)
				case *redis.Message:
					// place the messsage on the hub broadcast channel
					s.hub.Broadcast <- []byte(m.Payload)
				case error:
					// On error, close the channel and break out of the loop inner loop
					log.Errorf("Redis Error: %s", m)
					pubsub.Close()
					break ReceiveLoop
				}
			}
		}
	}
}

// Create a new Redis Subscription
func NewSubscription(h socket.Hub) *Subscription {
	// Create a new redis client - this does not connect to redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Create a new subscription
	return &Subscription{
		hub:     h,
		client:  client,
		channel: "fm:events",
	}
}
