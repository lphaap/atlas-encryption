package crypto

import (
	"main/lib"
	"main/lib/api"

	"github.com/gofiber/fiber/v2"
)

type DecryptRequest struct {
	Encrypted string `json:"encrypted"`
	Key       string `json:"key"`
}

func Decrypt(fiber *fiber.Ctx) error {
	var request DecryptRequest

	if err := fiber.BodyParser(&request); err != nil {
		error := lib.Object{"error": "Parameters encrypted and key are required"}
		return api.Response.BadRequest(fiber, error)
	}

	if request.Encrypted == "" || request.Key == "" {
		error := lib.Object{"error": "Parameters encrypted and key are required"}
		return api.Response.BadRequest(fiber, error)
	}

	decrypted, err := lib.Decrypt(request.Encrypted, request.Key)
	if err != nil {
		error := lib.Object{"error": "Decryption failed"}
		return api.Response.InternalServerError(fiber, error)
	}

	data := lib.Object{"decrypted": decrypted}
	return api.Response.OK(fiber, data)
}