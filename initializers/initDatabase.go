package initializers

import (
	"log"
	"os"
	"os/exec"
	"task2/database"
	"task2/models"
)

func SyncDatabase() {
	databaseURL := os.Getenv("DB_URL")

	cmd := exec.Command(
		"migrate",
		"-database", databaseURL,
		"-path", "./migrations",
		"up",
	)

	err := cmd.Run()
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Database migrated successfully.")
}

func RegisterModels() {
	database.DB.RegisterModel((*models.UserTask)(nil))
	database.DB.RegisterModel((*models.User)(nil))
	database.DB.RegisterModel((*models.Task)(nil))
}
