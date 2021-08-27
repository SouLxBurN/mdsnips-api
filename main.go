package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/soulxburn/soulxsnips/api"
	_ "github.com/soulxburn/soulxsnips/docs"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

// @title SouLxSnippets
// @version 1.0
// @description Backend API for storing and retrieving markdown snippets
// @tag.name md
// @host localhost:3000
// @BasePath
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	fiberApp := fiber.New()

	fiberApp.Get("/swagger/*", swagger.Handler)

	api.ConfigureMiddleware(fiberApp)
	api.ConfigureBasicAuth(fiberApp)
	api.ConfigureRoutes(fiberApp)

	port := flag.Int("port", 3000, "port to listen on")
	flag.Parse()
	fiberApp.Listen(":" + strconv.Itoa(*port))
}
