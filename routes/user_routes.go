package routes

import (
	"log"
	"net/http"
	"os"
	"task2/controllers"
	"task2/middlewares"
	"task2/models"
	"task2/schemas"
)

func UserRoutes(mux *http.ServeMux) {
	logger := log.New(os.Stdout, "API: ", log.LstdFlags)
	userController := controllers.NewUsers(logger)

	mux.HandleFunc("/signup", middlewares.MethodCheck("POST", middlewares.BindAndValidate(&models.User{}, userController.Signup)))
	mux.HandleFunc("/login", middlewares.MethodCheck("POST", middlewares.BindAndValidate(&schemas.LoginRequest{}, userController.Login)))
	mux.HandleFunc("/users", middlewares.MethodCheck("GET", middlewares.AuthMiddleware(userController.GetUsers)))
}
