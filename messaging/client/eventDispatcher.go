package client

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/unblockhub/api-gateway/auth"
	"github.com/unblockhub/api-gateway/cache"
	"github.com/unblockhub/api-gateway/messaging/node"
	"log"
	"os"
)

var (
	logger      = log.New(os.Stdout, "[GATEWAY][MESSAGING][CLIENT]", 0)
	activeUsers = make(map[string]**websocket.Conn)
	nodeName, _ = os.Hostname()
)

func createEvent(eventType string, error bool, data interface{}) map[string]interface{} {
	message := make(map[string]interface{})
	message["eventType"] = eventType
	message["error"] = error
	message["data"] = data
	return message
}

func getTotalUsers() float64 {
	return float64(len(cache.RedisClient.Keys(cache.Ctx, "websocket:clients").Val()))
}

func RunEventDispatcher(app *fiber.App) {
	// Unregister when told to
	node.Subscribe("API_DISCONNECT", func(raw []byte) {
		var event node.DisconnectUserMessage
		err := json.Unmarshal(raw, &event)
		if err != nil {
			logger.Printf("Failed to unmarshal payload")
			return
		}

		user := activeUsers[event.User]
		reason := make(map[string]interface{})

		_ = (*user).WriteJSON(reason)
		_ = (*user).Close()
		delete(activeUsers, event.User)
	})

	// Stats
	activeUsersCounter := promauto.NewCounterFunc(prometheus.CounterOpts{Name: "api_users_total"}, getTotalUsers)
	_ = prometheus.Register(activeUsersCounter)

	// Routes
	app.Get("/feed", websocket.New(func(conn *websocket.Conn) {
		userId := auth.GetUserId(conn.Cookies(auth.AccessTokenCookieName, ""))
		if userId == "" {
			_ = conn.WriteJSON(createEvent("DISCONNECT", true, "Invalid authentication."))
			_ = conn.Close()
			return
		}

		// Duplicate connection protection
		currentNodeName := cache.RedisClient.Get(cache.Ctx, fmt.Sprintf("gateway:clients:%s", userId)).Val()
		if currentNodeName != nodeName {
			// Duplicate exists on another node, dispatch to get it deleted
			eventData := make(map[string]interface{})
			eventData["disconnect_reason"] = "Logged in from another device"
			eventData["user_id"] = userId
			node.Publish(eventData, "API_EVENTS")
		}
		existing, exists := activeUsers[userId]
		if exists {
			_ = (*existing).WriteJSON(createEvent("DISCONNECT", true, "Logged in from another location"))
			_ = (*existing).Close()
		}

		activeUsers[userId] = &conn
		cache.RedisClient.Set(cache.Ctx, fmt.Sprintf("websocket:clients:%s", userId), nodeName, 0)
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
	app.Get("/totalUsers", func(ctx *fiber.Ctx) error {
		return ctx.SendString(fmt.Sprintf("%f", getTotalUsers()))
	})
}
