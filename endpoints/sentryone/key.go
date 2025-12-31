package sentryone

import (
	"os"
	"main/lib"
	"main/lib/api"
	"github.com/gofiber/fiber/v2"
)

type SentryOneKeyRequest struct {
	Encrypted string `json:"encrypted"`
}

func Key(fiber *fiber.Ctx) error {
	var request SentryOneKeyRequest

	if err := fiber.BodyParser(&request); err != nil {
		error := lib.Object{"error": "Parameter encrypted is required"}
		return api.Response.BadRequest(fiber, error)
	}

	if request.Encrypted == "" {
		error := lib.Object{"error": "Parameter encrypted is required"}
		return api.Response.BadRequest(fiber, error)
	}

	key := os.Getenv("SENTRYONE_KEY")
	if key == "" {
		error := lib.Object{
			"error": "Encryption key not found",
		}
		return api.Response.InternalServerError(fiber, error)
	}
	
	decrypted, err := lib.Decrypt(request.Encrypted, key)
	if err != nil {
		error := lib.Object{
			"error": "Decryption failed",
		}
		return api.Response.InternalServerError(fiber, error)
	}

	data := lib.Object{"decrypted": decrypted}
	return api.Response.OK(fiber, data)
}