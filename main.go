package main

import (
	"log"

	"github.com/azkaainurridho514/bimskrip/database"
	"github.com/azkaainurridho514/bimskrip/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
)


func main() {
	database.InitDB()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	router.SetupRoutes(app)
	// log.Fatal(app.Listen("192.168.245.114:8080"))
	log.Fatal(app.Listen("192.168.100.6:8080"))
}