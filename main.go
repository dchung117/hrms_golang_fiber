package main

import (
	"log"

	"github.com/dchung117/hrms_golang_fiber/db"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// connect to mongodb instance
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}
	// create fiber app
	app := fiber.New()

	// define routes
	app.Get("/employee", db.GetEmployee)
	app.Post("/employee", db.AddEmployee)
	app.Put("/employee/:id", db.UpdateEmployee)
	app.Delete("/employee/:id", db.DeleteEmployee)
}
