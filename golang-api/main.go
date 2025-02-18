package main

import (
	"log"
	"net/http"

	"golang-api/config"
	"golang-api/routes"

	"github.com/rs/cors"
)

func main() {
	// Inisialisasi database & logger
	config.InitDB()
	config.InitLogger()

	router := routes.SetupRoutes()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	config.Logger.Info("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler.Handler(router)))
}
