package routes

import (
	"github.com/BRAVO68WEB/go-rest-api/controllers"
	"github.com/gofiber/fiber/v2"
)

func GistsRoute(route fiber.Router) {
    route.Get("/", controllers.GetAllGists)
    route.Get("/:id", controllers.GetGist)
    route.Post("/", controllers.AddGist)
    route.Put("/:id", controllers.UpdateGist)
    route.Delete("/:id", controllers.DeleteGist)
}