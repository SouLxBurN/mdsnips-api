package api

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

// ConfigureBasicAuth
// Configures and attaches gofiber basic auth middleware
func ConfigureBasicAuth(app *fiber.App) {
	baseUser := os.Getenv("soulxsnips_user")
	basePass := os.Getenv("soulxsnips_pass")
	users := map[string]string{
		baseUser: basePass,
	}
	app.Use(basicauth.New(basicauth.Config{
		Users: users,
		Realm: "Forbidden",
		Authorizer: func(user, pass string) bool {
			creds, ok := users[user]
			if ok {
				if pass == creds {
					return true
				}
			}
			return false
		},
	}))
}
