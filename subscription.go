// Temporary Redis Event Subscription to pass events from
// API's to the Player clients

package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v3"
)

// Redis event subscription
type subscription struct {
	addr    string
	channel string
}

// Consime events from redis, passing them onto the WS hub for
// broadcast to clients
func (s *subscription) consume() {
	for {
		client := redis.NewClient(&redis.Options{
			Addr:     s.addr,
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		pubsub := client.PubSub()
		err := pubsub.Subscribe(s.channel)
		if err != nil {
			log.Errorf("Redis Connection: %s", err)
			time.Sleep(time.Second)
			continue
		}
	ReceiveLoop:
		for {
			msg, err := pubsub.Receive() // recieve a message from the channel
			if err != nil {
				log.Error(err)
				pubsub.Close()
				break ReceiveLoop
			} else {
				switch m := msg.(type) { // Switch the mesage type
				case *redis.Subscription:
					log.Debugf("Subscribed: %s", m.Channel)
				case *redis.Message:
					// place the messsage on the hub broadcast channel
					h.broadcast <- []byte(m.Payload)
				case error:
					log.Errorf("Redis Error: %s", m)
					pubsub.Close()
					break ReceiveLoop
				}
			}
		}
	}
}

// Create a new Redis Subscription
func NewSubscription() *subscription {
	return &subscription{
		addr:    "localhost:6379",
		channel: "fm:events",
	}
}
