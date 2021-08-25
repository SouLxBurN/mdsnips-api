package api

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/soulxburn/soulxsnips/client"
	"github.com/soulxburn/soulxsnips/md"
)

// ConfiugureRoutes
// Imports and configures various routes for
// all modules.
func ConfigureRoutes(app *fiber.App) {
	mongoHost := os.Getenv("SOULXSNIPS_MONGO_HOST")
	mongoPort := os.Getenv("SOULXSNIPS_MONGO_PORT")
	mongoUser := os.Getenv("SOULXSNIPS_MONGO_USER")
	mongoPass := os.Getenv("SOULXSNIPS_MONGO_PASS")
	mClient, err := client.InitMongoClient(mongoHost, mongoPort, mongoUser, mongoPass)
	if err != nil {
		log.Fatal("Failed to establish connection to Mongo")
	}
	mdService := md.InitMDService(mClient)
	mdHandlers := md.InitMDHandlers(mdService)

	app.Post("/md", mdHandlers.CreateMDHandler)
	app.Patch("/md", mdHandlers.UpdateMDHandler)
	app.Get("/md/:id", mdHandlers.GetMDHandler)
	app.Get("/md", mdHandlers.GetAllMDHandler)
	app.Delete("/md/:id", mdHandlers.DeleteMDHandler)
}
