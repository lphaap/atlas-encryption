package endpoints

import (
	"main/lib"
	"main/lib/api"
	"github.com/gofiber/fiber/v2"
)

func Status(fiber *fiber.Ctx) error {
	fiber = api.JsonEndpoint(fiber)
	data := lib.Object{
		"status": "API is online",
	}
	return api.Response.OK(fiber, data)
}
