package routes

import (
	"log"
	"net/http"
	"os"
	"strings"
	"task2/controllers"
	"task2/helperFunc"
	"task2/middlewares"
	"task2/schemas"
)

func TaskRoutes(mux *http.ServeMux) {
	logger := log.New(os.Stdout, "API: ", log.LstdFlags)
	taskController := controllers.NewTasks(logger)

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler := middlewares.AuthMiddleware(middlewares.BindAndValidate(&schemas.CreateTaskRequest{}, taskController.CreateTask))
			handler(w, r)
		case http.MethodGet:
			handler := middlewares.AuthMiddleware(http.HandlerFunc(taskController.GetTasks))
			handler.ServeHTTP(w, r)
		default:
			helperFunc.RespondJSON(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		}
	})

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/tasks/")
		parts := strings.Split(path, "/")

		handler := http.NotFoundHandler()

		switch {
		case len(parts) == 1:
			if r.Method == http.MethodGet {
				handler = middlewares.MethodCheck(http.MethodGet,
					middlewares.AuthMiddleware(http.HandlerFunc(taskController.GetTaskByID)))
			} else if r.Method == http.MethodDelete {
				handler = middlewares.MethodCheck(http.MethodDelete,
					middlewares.AuthMiddleware(http.HandlerFunc(taskController.DeleteTask)))
			}
		case len(parts) == 2 && parts[1] == "complete":
			handler = middlewares.MethodCheck(http.MethodPost,
				middlewares.AuthMiddleware(http.HandlerFunc(taskController.CompleteTask)))
		case len(parts) == 3 && parts[1] == "assign":
			handler = middlewares.MethodCheck(http.MethodPost,
				middlewares.AuthMiddleware(http.HandlerFunc(taskController.AssignTask)))
		}

		handler.ServeHTTP(w, r)
	})

	mux.HandleFunc("/users/tasks", func(w http.ResponseWriter, r *http.Request) {
		handler := middlewares.MethodCheck(http.MethodGet,
			middlewares.AuthMiddleware(http.HandlerFunc(taskController.GetUserTasks)))
		handler(w, r)
	})

	mux.HandleFunc("/users/tasks/created", func(w http.ResponseWriter, r *http.Request) {
		handler := middlewares.MethodCheck(http.MethodGet,
			middlewares.AuthMiddleware(http.HandlerFunc(taskController.GetTasksCreatedByUser)))
		handler(w, r)
	})

	mux.HandleFunc("/users/tasks/assigned", func(w http.ResponseWriter, r *http.Request) {
		handler := middlewares.MethodCheck(http.MethodGet,
			middlewares.AuthMiddleware(http.HandlerFunc(taskController.GetTasksAssignedToUser)))
		handler(w, r)
	})
}