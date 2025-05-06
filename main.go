package main

import (
	"log"
	"net/http"
	"task2/database"
	"task2/initializers"
	"task2/routes"
)

func init() {
	initializers.LoadEnvVariables()

	database.ConnectToDB()

	initializers.RegisterModels()

	initializers.SyncDatabase()
}

func main() {
	mux := http.NewServeMux()

	routes.UserRoutes(mux)
	routes.TaskRoutes(mux)

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
