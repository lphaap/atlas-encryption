package endpoints

import (
	sentryone "main/endpoints/sentryone"
	sentrytwo "main/endpoints/sentrytwo"
	atlas "main/endpoints/atlas"
	crypto "main/endpoints/crypto"
	
	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) *fiber.App {

	app.Get("/", Welcome)
	app.Get("/status", Status)
	
	app.Get("/crypto/generate", crypto.Generate)
	app.Post("/crypto/decrypt", crypto.Decrypt)
	app.Post("/crypto/encrypt", crypto.Encrypt)
	
	app.Post("/sentryone/key", sentryone.Key)
	app.Post("/sentrytwo/key", sentrytwo.Key)
	app.Post("/atlas/key", atlas.Key)

	return app
}
