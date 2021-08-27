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
	mongoConn := os.Getenv("SOULXSNIPS_MONGO_CONN")
	mClient, err := client.InitMongoClient(mongoConn)
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
