package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/soulxburn/soulxsnips/api"
	"github.com/soulxburn/soulxsnips/docs"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

// @title SouLxSnippets
// @version 1.0
// @description Backend API for storing and retrieving markdown snippets
// @tag.name md
// @BasePath
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
    host := os.Getenv("HOST")
	port := os.Getenv("PORT")
    if host == "" {
        host = "localhost"
    }
	if port == "" {
		port = "3000"
	}

	docs.SwaggerInfo.Host = host + ":" + port

	fiberApp := fiber.New()
	fiberApp.Get("/swagger/*", swagger.Handler)

	api.ConfigureMiddleware(fiberApp)
	api.ConfigureBasicAuth(fiberApp)
	api.ConfigureRoutes(fiberApp)

	fiberApp.Listen(":" + port)
}
