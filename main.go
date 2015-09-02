// Serves the Websocket Server

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/thisissoon/FM-Perceptor/middleware"
	"github.com/thisissoon/FM-Perceptor/rest"
	"github.com/thisissoon/FM-Perceptor/socket"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

// Entrypoint - Runs the WS Server
func main() {
	log.SetLevel(log.DebugLevel)

	// Message Hub
	hub := socket.NewHub()
	go hub.Run()

	// WS Service
	ws := socket.NewWSService(hub)

	// Redis Connection
	s := NewSubscription(hub)
	go s.consume()

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

	//
	c.Get("/playlist/next", rest.GetNextTrackHandler)

	graceful.ListenAndServe(":9000", c)
}
