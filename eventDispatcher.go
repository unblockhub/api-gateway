package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	activeUsers = make(map[string]**websocket.Conn)
)

func createEvent(eventType string, error bool, data interface{}) map[string]interface{} {
	message := make(map[string]interface{})
	message["eventType"] = eventType
	message["error"] = error
	message["data"] = data
	return message
}

func runEventDispatcher() {
	// Stats
	activeUsersCounter := promauto.NewCounterFunc(prometheus.CounterOpts{Name: "api_totalusers"}, func() float64 {
		return float64(len(activeUsers))
	})
	_ = prometheus.Register(activeUsersCounter)


	// Routes
	app.Get("/api/feed", websocket.New(func(conn *websocket.Conn) {
		userId := getUserId(conn.Cookies(accessTokenCookieName, ""))
		if userId == "" {
			_ = conn.WriteJSON(createEvent("DISCONNECT", true, "Invalid authentication."))
			_ = conn.Close()
			return
		}
		existing, exists := activeUsers[userId]
		if exists {
			_ = (*existing).WriteJSON(createEvent("DISCONNECT", true, "Logged in from another location"))
			(*existing).Close()
		}

		activeUsers[userId] = &conn
		logger.Printf("User %s connected to news feed", userId)

		for true {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
		if activeUsers[userId] == &conn {
			logger.Printf("Deleting connection")
			delete(activeUsers, userId)
		}
	}))
	app.Get("/api/totalUsers", func(ctx *fiber.Ctx) error {
		return ctx.SendString(fmt.Sprintf("%d", len(activeUsers)))
	})
}
