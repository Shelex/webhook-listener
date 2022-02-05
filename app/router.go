package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ProvideRouter() *fiber.App {
	router := fiber.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Access-Control-Allow-Origin, Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, DELETE, OPTIONS",
	}))
	//router.Use(logger.New())
	return router
}
