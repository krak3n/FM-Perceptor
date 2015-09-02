// Serves the Websocket Server

package main

import (
	"fmt"

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
	// Set Log Level
	log.SetLevel(log.DebugLevel)

	// Defaults
	viper.SetDefault("port", "5000")
	viper.SetDefault("redis_host", "127.0.0.1")
	viper.SetDefault("redis_port", "6379")

	// From file
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

	// From environment vars - Only top level are configured from env vars
	viper.SetEnvPrefix("perceptor")
	viper.AutomaticEnv()
}

// Application Entrypoint
func main() {
	// Message Hub
	hub := socket.NewHub()
	go hub.Run()

	// WS Service
	ws := socket.NewWSService(hub)

	// Redis Subscription
	s := pubsub.NewSubscription(
		hub,
		viper.GetString("redis_host"),
		viper.GetString("redis_port"))
	go s.Consume()

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

	// Start Serving Web Application
	log.Debugf("Starting Server on :%s", viper.GetString("port"))
	graceful.ListenAndServe(fmt.Sprintf(":%s", viper.GetString("port")), c)
}
