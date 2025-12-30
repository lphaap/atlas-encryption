package crypto

import (
	"main/lib"
	"main/lib/api"
	"github.com/gofiber/fiber/v2"
)

func Generate(fiber *fiber.Ctx) error {
	fiber = api.JsonEndpoint(fiber)
	key, _ := lib.Key(32)
	
	data := lib.Object{ "key": key }
	return api.Response.Created(fiber, data)
}