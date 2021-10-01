package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// ConfigureMiddleware
// Configures various GoFiber middleware
// i.e. Recover, Limiter, etc.
func ConfigureMiddleware(app *fiber.App) {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Use(cors.New())

	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 30 * time.Second,
		KeyGenerator: func(ctx *fiber.Ctx) string {
			return getRequestIP(ctx)
		},
		LimitReached: func(ctx *fiber.Ctx) error {
			log.Printf("Too many requests received from: %s\n", getRequestIP(ctx))
			return ctx.SendStatus(http.StatusTooManyRequests)
		},
	}))

	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02T15:04:05-0700",
		TimeZone:   "UTC",
	}))

}

// getRequestIP
// Returns the 'x-forwarded-for' field,
// if not present returns the ip address on the context.
func getRequestIP(c *fiber.Ctx) string {
	ip := c.Get("x-forwarded-for")
	if ip == "" {
		ip = c.IP()
	}
	return ip
}
