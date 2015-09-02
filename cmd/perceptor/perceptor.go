// Serves the Websocket Server

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thisissoon/FM-Perceptor/middleware"
	"github.com/thisissoon/FM-Perceptor/pubsub"
	"github.com/thisissoon/FM-Perceptor/rest"
	"github.com/thisissoon/FM-Perceptor/socket"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

func init() {
	viper.SetConfigName("perceptor")        // name of config file (without extension)
	viper.AddConfigPath("/etc/perceptor/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.perceptor") // call multiple times to add many search paths
	viper.AddConfigPath("$PWD/.perceptor")  // call multiple times to add many search paths
	err := viper.ReadInConfig()             // Find and read the config file
	if err != nil {                         // Handle errors reading the config file
		log.Warnf("No config file found or is not properly formatted: %s", err)
	} else {
		log.Info("Config Loaded from File")
	}
}

// Application Entrypoint
func main() {
	log.SetLevel(log.DebugLevel)

	// Message Hub
	hub := socket.NewHub()
	go hub.Run()

	// WS Service
	ws := socket.NewWSService(hub)

	// Redis Connection
	s := pubsub.NewSubscription(hub)
	go s.Consume()

	// Serve the WS Server
	log.Debug("Starting Websocket Server on :9000")

	c := web.New()

	// Middlewares
	c.Use(middleware.SetupEnv)
	c.Use(middleware.HMACVerification)
	c.Use(middleware.CustomValidators)
	c.Use(middleware.RedisClient)

	// WS Connection Handler
	c.Get("/", ws.Handler)

	// Event REST endpoints
	c.Post("/events/play", rest.PlayCreateHandler)
	c.Post("/events/end", rest.EndCreateHandler)
	c.Post("/events/pause", rest.PauseCreateHandler)
	c.Post("/events/resume", rest.ResumeCreateHandler)

	// Updates to Mute / Volume States
	c.Put("/volume", rest.VolumeUpdateHandler)
	c.Put("/mute", rest.MuteUpdateHandler)

	// Get the next track from the playlist
	c.Get("/playlist/next", rest.GetNextTrackHandler)

	graceful.ListenAndServe(":9000", c)
}
