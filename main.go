package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/soulxburn/mdsnips/api"
	"github.com/soulxburn/mdsnips/docs"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

// @title MDSnips
// @version 1.0
// @description API for storing and retrieving markdown snippets.\nBuilt live on stream @twitch.tv/soulxburn
// @tag.name md
// @BasePath
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	host := os.Getenv("HOST")
	if host == "" || strings.Contains(host, "localhost") {
		host = "localhost:" + port
	}
	docs.SwaggerInfo.Host = host

	fiberApp := fiber.New()
	fiberApp.Get("/swagger/*", swagger.Handler)
	fiberApp.All("/", func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/swagger/index.html", http.StatusMovedPermanently)
	})

	api.ConfigureMiddleware(fiberApp)
	api.ConfigureBasicAuth(fiberApp)
	api.ConfigureRoutes(fiberApp)

	fiberApp.Listen(":" + port)
}
