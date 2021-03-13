package main

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"time"
)

var (
	logger   = log.New(os.Stdout, "[GATEWAY]", 0)
	app = fiber.New()
)


func main() {
	// TODO: Prefork

	// Metrics
	prometheus := fiberprometheus.New("api-gateway")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})


	// Modules
	runEventDispatcher()

	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}


	logger.Fatal(app.Listen(":8080"))
}


