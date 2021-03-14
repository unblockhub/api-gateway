package main

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/unblockhub/api-gateway/cache"
	"github.com/unblockhub/api-gateway/messaging/client"
	"github.com/unblockhub/api-gateway/messaging/node"
	"log"
	"os"
)

var (
	logger      = log.New(os.Stdout, "[GATEWAY]", 0)
	app         = fiber.New()

)


func main() {
	// Init modules
	cache.Init()
	node.Init()

	// Metrics
	prometheus := fiberprometheus.New("api-gateway")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})


	// Modules
	client.RunEventDispatcher(app)

	logger.Fatal(app.Listen(":8080"))
}


