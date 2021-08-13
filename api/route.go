package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/soulxburn/soulxsnips/md"
)

func ConfigureRoutes(app *fiber.App) {
	app.Post("/md", md.CreateMDHandler)
	app.Patch("/md", md.UpdateMDHandler)
	app.Get("/md/:id", md.GetMDHandler)
	app.Get("/md", md.GetAllMDHandler)
	app.Delete("/md/:id", md.DeleteMDHandler)
}
