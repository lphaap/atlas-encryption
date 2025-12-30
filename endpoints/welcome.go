package endpoints

import (
	"main/lib"
	"main/lib/api"
	"github.com/gofiber/fiber/v2"
)

func Welcome(fiber *fiber.Ctx) error {
	fiber = api.JsonEndpoint(fiber)
	var data = lib.Object{"greetings": "API is online"}
	return api.Response.OK(fiber, data)
}
