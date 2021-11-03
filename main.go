package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/soulxburn/mdsnips/api"
	"github.com/soulxburn/mdsnips/client"
	"github.com/soulxburn/mdsnips/docs"
	"github.com/soulxburn/mdsnips/md"
	"go.mongodb.org/mongo-driver/mongo"

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
	api.ConfigureMiddleware(fiberApp)

	fiberApp.Get("/swagger/*", swagger.Handler)
	fiberApp.All("/", func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/swagger/index.html", http.StatusMovedPermanently)
	})

	api.ConfigureBasicAuth(fiberApp)

	mClient := getMongoConnection()
	md.ConfigureIndexes(mClient)

	mdService := md.InitMDService(mClient)
	mdHandlers := md.InitMDHandlers(mdService)
	mdHandlers.ConfigureRoutes(fiberApp)

	if err := fiberApp.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

// Initialize MongoClient
func getMongoConnection() *mongo.Client {
	mongoConn := os.Getenv("MDSNIPS_MONGO_CONN")
	mClient, err := client.InitMongoClient(mongoConn)
	if err != nil {
		log.Fatal(err)
	}

	return mClient
}
