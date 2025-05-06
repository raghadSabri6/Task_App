package controllers

import (
	"log"
	"net/http"
	"strings"
	"task2/database"
	"task2/helperFunc"
	"task2/middlewares"
	"task2/models"
	"task2/schemas"

	"github.com/google/uuid"
)

type Tasks struct {
	L *log.Logger
}

func NewTasks(l *log.Logger) *Tasks {
	return &Tasks{L: l}
}

func (t *Tasks) CompleteTask(w http.ResponseWriter, r *http.Request) {
	uuidStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	uuidStr = strings.TrimSuffix(uuidStr, "/complete")
	taskUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusBadRequest, "Invalid task UUID", nil)
		return
	}

	task, err := models.FetchTaskByUUID(taskUUID)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusNotFound, "Task not found", nil)
		return
	}

	userUUID := helperFunc.GetUserUUIDFromRequest(r)
	if err := task.Complete(userUUID); err != nil {
		helperFunc.RespondJSON(w, http.StatusForbidden, err.Error(), nil)
		return
	}

	updatedTask, err := models.FetchTaskByUUID(taskUUID)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusInternalServerError, "Failed to fetch updated task", nil)
		return
	}

	helperFunc.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"task": updatedTask})
}

func (t *Tasks) CreateTask(w http.ResponseWriter, r *http.Request) {
	userUUID := helperFunc.GetUserUUIDFromRequest(r)

	ctx := r.Context()

	val := ctx.Value(middlewares.BindKey)
	taskReq, ok := val.(*schemas.CreateTaskRequest)
	if !ok {
		helperFunc.RespondJSON(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	task := &models.Task{
		Title:       taskReq.Title,
		Description: taskReq.Description,
	}

	if len(taskReq.Users) > 0 {
		var users []models.User
		for _, userAssign := range taskReq.Users {
			var dbAssignedUser models.User
			err := database.DB.NewSelect().Model(&dbAssignedUser).Where("uuid = ?", userAssign.ID).Scan(ctx)
			if err != nil {
				helperFunc.RespondJSON(w, http.StatusBadRequest, "One or more assigned users do not exist", nil)
				return
			}
			users = append(users, models.User{UUID: userAssign.ID})
		}

		var userPtrs []*models.User
		for i := range users {
			userPtrs = append(userPtrs, &users[i])
		}
		task.Users = userPtrs

		// Set the primary assignee to the first user in the list
		if len(users) > 0 {
			assignedID := users[0].UUID
			task.AssignedToID = &assignedID
		}
	}

	if err := task.ValidateAndCreateTask(userUUID); err != nil {
		helperFunc.RespondJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	createdTask, err := models.FetchTaskByUUID(task.UUID)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusInternalServerError, "Failed to fetch created task", nil)
		return
	}

	helperFunc.RespondJSON(w, http.StatusCreated, "", map[string]interface{}{"task": createdTask})
}

func (t *Tasks) DeleteTask(w http.ResponseWriter, r *http.Request) {
	uuidStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	taskUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusBadRequest, "Invalid task UUID", nil)
		return
	}

	task, err := models.FetchTaskByUUID(taskUUID)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusNotFound, "Task not found", nil)
		return
	}

	userUUID := helperFunc.GetUserUUIDFromRequest(r)
	if err := task.DeleteTask(userUUID); err != nil {
		if err.Error() == "unauthorized: you can only delete your own tasks" {
			helperFunc.RespondJSON(w, http.StatusForbidden, err.Error(), nil)
		} else {
			helperFunc.RespondJSON(w, http.StatusInternalServerError, "Failed to delete task", nil)
		}
		return
	}

	helperFunc.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"task": task})
}

func (t *Tasks) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	uuidStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	taskUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusBadRequest, "Invalid task UUID", nil)
		return
	}

	task, err := models.FetchTaskByUUID(taskUUID)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusNotFound, "Task not found", nil)
		return
	}

	helperFunc.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"task": task})
}

func (t *Tasks) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := models.GetTasks()
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusInternalServerError, "Failed to fetch tasks", nil)
		return
	}

	helperFunc.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"tasks": tasks})
}

func (t *Tasks) GetUserTasks(w http.ResponseWriter, r *http.Request) {
	userUUID := helperFunc.GetUserUUIDFromRequest(r)

	if userUUID == uuid.Nil {
		helperFunc.RespondJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	tasks, err := models.GetUserTasks(userUUID)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusInternalServerError, "Failed to fetch user tasks", nil)
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	helperFunc.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"tasks": tasks})
}

func (t *Tasks) GetTasksCreatedByUser(w http.ResponseWriter, r *http.Request) {
	userUUID := helperFunc.GetUserUUIDFromRequest(r)

	tasks, err := models.GetTasksCreatedByUser(userUUID)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusInternalServerError, "Failed to fetch tasks created by user", nil)
		return
	}

	helperFunc.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"tasks": tasks})
}

func (t *Tasks) GetTasksAssignedToUser(w http.ResponseWriter, r *http.Request) {
	userUUID := helperFunc.GetUserUUIDFromRequest(r)

	if userUUID == uuid.Nil {
		helperFunc.RespondJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	tasks, err := models.GetTasksAssignedToUser(userUUID)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusInternalServerError, "Failed to fetch tasks assigned to user", nil)
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	helperFunc.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"tasks": tasks})
}

func (t *Tasks) AssignTask(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	parts := strings.Split(path, "/")

	if len(parts) < 3 || parts[1] != "assign" {
		helperFunc.RespondJSON(w, http.StatusBadRequest, "Invalid path format", nil)
		return
	}

	taskUUID, err := uuid.Parse(parts[0])
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusBadRequest, "Invalid task UUID", nil)
		return
	}

	assignedUserUUID, err := uuid.Parse(parts[2])
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusBadRequest, "Invalid user UUID", nil)
		return
	}

	task, err := models.FetchTaskByUUID(taskUUID)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusNotFound, "Task not found", nil)
		return
	}

	creatorUUID := helperFunc.GetUserUUIDFromRequest(r)
	ctx := r.Context()

	if err := task.ValidateAndAssignTask(ctx, database.DB, creatorUUID, assignedUserUUID); err != nil {
		helperFunc.RespondJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	updatedTask, err := models.FetchTaskByUUID(task.UUID)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusInternalServerError, "Failed to fetch updated task", nil)
		return
	}

	helperFunc.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"task": updatedTask})
}
