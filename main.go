package main

import (
	"log"
	"main/endpoints"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env", "../.env")
	if err != nil {
		log.Println("Warning: .env file not found")
	} else {
		log.Println(".env file loaded successfully")
	}

	app := fiber.New()
	endpoints.Router(app)

	log.Fatal(app.Listen(":8001"))
}
