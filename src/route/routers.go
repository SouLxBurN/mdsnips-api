package route

import (
	"soulxsnips/src/handler"

	"github.com/gofiber/fiber/v2"
)

func Configure(app *fiber.App) {
	app.Post("/md", handler.CreateMD)
	app.Patch("/md", handler.UpdateMD)
	app.Get("/md/:id", handler.GetMD)
	app.Get("/md", handler.GetAllMD)
	app.Delete("/md/:id", handler.DeleteMD)
}
