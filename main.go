package main

import (
	"learn-api/config"
	_ "learn-api/docs" // Import for Swagger docs
	"learn-api/routes"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Learn API
// @version 1.0
// @description This is a sample server for managing users.
// @host localhost:8000
// @BasePath /

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
func main() {
	config.Connect()
	r := routes.SetupRoutes()

	// Swagger endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
