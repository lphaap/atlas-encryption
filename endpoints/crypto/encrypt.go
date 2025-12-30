package crypto

import (
	"main/lib"
	"main/lib/api"

	"github.com/gofiber/fiber/v2"
)

type EncryptRequest struct {
	Plain string `json:"plain"`
	Key   string `json:"key"`
}

func Encrypt(fiber *fiber.Ctx) error {
	var request EncryptRequest

	if err := fiber.BodyParser(&request); err != nil {
		error := lib.Object{"error": "Parameters plain and key are required"}
		return api.Response.BadRequest(fiber, error)
	}

	if request.Plain == "" || request.Key == "" {
		error := lib.Object{"error": "Parameters plain and key are required"}
		return api.Response.BadRequest(fiber, error)
	}

	encrypted, err := lib.Encrypt(request.Plain, request.Key)
	if err != nil {
		error := lib.Object{"error": "Encryption failed"}
		return api.Response.InternalServerError(fiber, error)
	}

	data := lib.Object{"encrypted": encrypted}
	return api.Response.OK(fiber, data)
}
