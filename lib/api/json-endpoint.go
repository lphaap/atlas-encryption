package api

import (
	"github.com/gofiber/fiber/v2"
)

func JsonEndpoint(fiber *fiber.Ctx) *fiber.Ctx {
	fiber.Accepts("application/json")
	fiber.Accepts("text", "json")
	fiber.Context().SetContentType("application/json")
	return fiber
}