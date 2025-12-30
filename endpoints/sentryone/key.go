package sentryone

import (
	"main/lib"
	"main/lib/api"
	"github.com/gofiber/fiber/v2"
)

func Key(fiber *fiber.Ctx) error {
	fiber = api.JsonEndpoint(fiber)
	var key = "example_key"
	
	data := lib.Object{ "key": key }
	return api.Response.OK(fiber, data)
}