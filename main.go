package main

import (
	"log"
	"main/endpoints"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := fiber.New()
	endpoints.Router(app)

	log.Fatal(app.Listen(":8001"))
}
