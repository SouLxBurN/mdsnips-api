package api

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/soulxburn/mdsnips/client"
	"github.com/soulxburn/mdsnips/md"
)

// ConfiugureRoutes
// Imports and configures various routes for
// all modules.
func ConfigureRoutes(app *fiber.App) {
	mongoConn := os.Getenv("MDSNIPS_MONGO_CONN")
	mClient, err := client.InitMongoClient(mongoConn)
	if err != nil {
		log.Fatal("Failed to establish connection to Mongo")
	}
	mdService := md.InitMDService(mClient)
	mdHandlers := md.InitMDHandlers(mdService)

	// Redirect root requests to swagger documentation
	app.Post("/md", mdHandlers.CreateMDHandler)
	app.Patch("/md", mdHandlers.UpdateMDHandler)
	app.Get("/md/:id", mdHandlers.GetMDHandler)
	app.Get("/md", mdHandlers.GetAllMDHandler)
	app.Delete("/md/:id", mdHandlers.DeleteMDHandler)
}
