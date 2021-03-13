package main

import (
	"github.com/gofiber/websocket/v2"
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
}
