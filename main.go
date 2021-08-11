package main

import (
	"os"
	"soulxsnips/src/route"

	_ "soulxsnips/docs"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

// @title SouLxSnippets
// @version 1.0
// @description Backend API for storing and retrieving markdown snippets
// @tag.name md
// @host localhost:3000
// @BasePath /
func main() {
	app := fiber.New()
	app.Get("/swagger/*", swagger.Handler)
	route.Configure(app)

	var port string
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		port = ":3000"
	}
	app.Listen(port)
}
