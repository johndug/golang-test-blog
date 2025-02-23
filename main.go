package main

import (
	"log"
	"net/http"
	"test-ai-api/init/db"
	"test-ai-api/routes"
)

func main() {
	database, err := db.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	handler := routes.SetupRoutes(database)
	log.Printf("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
