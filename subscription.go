// Temporary Redis Event Subscription to pass events from
// API's to the Player clients

package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v3"
)

// Redis event subscription
type subscription struct {
	client  *redis.Client
	channel *string
}

// Consime events from redis, passing them onto the WS hub for
// broadcast to clients
func (s *subscription) consume() {
	// Subscribe to channel, exiting the program on fail
	pubsub := s.client.PubSub()
	err := pubsub.Subscribe(*s.channel)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure connection the channel is closed on exit
	defer pubsub.Close()

	// Loop to recieve events
	for {
		msg, err := pubsub.Receive() // recieve a message from the channel
		if err != nil {
			log.Error(err)
		} else {
			switch m := msg.(type) { // Switch the mesage type
			case *redis.Subscription:
				log.Debugf("Subscribed: %s", m.Channel)
			case *redis.Message:
				// place the messsage on the hub broadcast channel
				h.broadcast <- []byte(m.Payload)
			default:
				log.Errorf("Unknwon Message Type")
			}
		}
	}
}

// Create a new Redis Subscription
func NewSubscription(c *redis.Client, ch *string) *subscription {
	return &subscription{
		client:  c,
		channel: ch,
	}
}
