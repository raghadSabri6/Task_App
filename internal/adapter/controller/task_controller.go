package controller

import (
	"net/http"
	"strings"
	"task2/internal/app/dto"
	"task2/internal/app/usecase"
	"task2/internal/infrastructure/middleware"
	"task2/pkg/utils"

	"github.com/google/uuid"
)

// TaskController handles HTTP requests for tasks
type TaskController struct {
	taskUseCase *usecase.TaskUseCase
}

// NewTaskController creates a new task controller
func NewTaskController(taskUseCase *usecase.TaskUseCase) *TaskController {
	return &TaskController{
		taskUseCase: taskUseCase,
	}
}

// CreateTask handles the creation of a new task
func (c *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Get user UUID from context
	userUUID := utils.GetUserUUIDFromRequest(r)
	
	// Get request body from context
	ctx := r.Context()
	val := ctx.Value(middleware.BindKey)
	taskReq, ok := val.(*dto.CreateTaskRequest)
	if !ok {
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	
	// Create task
	task, err := c.taskUseCase.CreateTask(ctx, taskReq, userUUID)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	
	utils.RespondJSON(w, http.StatusCreated, "", map[string]interface{}{"task": task})
}

// GetTaskByID handles getting a task by ID
func (c *TaskController) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	// Extract task UUID from path
	uuidStr := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	taskUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid task UUID", nil)
		return
	}
	
	// Get task
	task, err := c.taskUseCase.GetTaskByUUID(r.Context(), taskUUID)
	if err != nil {
		utils.RespondJSON(w, http.StatusNotFound, "Task not found", nil)
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"task": task})
}

// GetAllTasks handles getting all tasks
func (c *TaskController) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	// Get all tasks
	tasksResp, err := c.taskUseCase.GetAllTasks(r.Context())
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to fetch tasks", nil)
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"tasks": tasksResp.Tasks})
}

// GetTasksCreatedByUser handles getting tasks created by a user
func (c *TaskController) GetTasksCreatedByUser(w http.ResponseWriter, r *http.Request) {
	// Get user UUID from context
	userUUID := utils.GetUserUUIDFromRequest(r)
	
	// Get tasks created by user
	tasksResp, err := c.taskUseCase.GetTasksCreatedByUser(r.Context(), userUUID)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to fetch tasks created by user", nil)
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"tasks": tasksResp.Tasks})
}

// GetTasksAssignedToUser handles getting tasks assigned to a user
func (c *TaskController) GetTasksAssignedToUser(w http.ResponseWriter, r *http.Request) {
	// Get user UUID from context
	userUUID := utils.GetUserUUIDFromRequest(r)
	
	// Get tasks assigned to user
	tasksResp, err := c.taskUseCase.GetTasksAssignedToUser(r.Context(), userUUID)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to fetch tasks assigned to user", nil)
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"tasks": tasksResp.Tasks})
}

// CompleteTask handles completing a task
func (c *TaskController) CompleteTask(w http.ResponseWriter, r *http.Request) {
	// Extract task UUID from path
	uuidStr := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	uuidStr = strings.TrimSuffix(uuidStr, "/complete")
	taskUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid task UUID", nil)
		return
	}
	
	// Get user UUID from context
	userUUID := utils.GetUserUUIDFromRequest(r)
	
	// Complete task
	task, err := c.taskUseCase.CompleteTask(r.Context(), taskUUID, userUUID)
	if err != nil {
		utils.RespondJSON(w, http.StatusForbidden, err.Error(), nil)
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"task": task})
}

// DeleteTask handles deleting a task
func (c *TaskController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Extract task UUID from path
	uuidStr := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	taskUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid task UUID", nil)
		return
	}
	
	// Get user UUID from context
	userUUID := utils.GetUserUUIDFromRequest(r)
	
	// Delete task
	err = c.taskUseCase.DeleteTask(r.Context(), taskUUID, userUUID)
	if err != nil {
		if err.Error() == "only the task creator can delete the task" {
			utils.RespondJSON(w, http.StatusForbidden, err.Error(), nil)
		} else {
			utils.RespondJSON(w, http.StatusInternalServerError, "Failed to delete task", nil)
		}
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, "Task deleted successfully", nil)
}

// AssignTask handles assigning a task to a user
func (c *TaskController) AssignTask(w http.ResponseWriter, r *http.Request) {
	// Extract task UUID and user UUID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	parts := strings.Split(path, "/")
	
	if len(parts) < 3 || parts[1] != "assign" {
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid path format", nil)
		return
	}
	
	taskUUID, err := uuid.Parse(parts[0])
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid task UUID", nil)
		return
	}
	
	assignedUserUUID, err := uuid.Parse(parts[2])
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid user UUID", nil)
		return
	}
	
	// Get user UUID from context
	requestorUUID := utils.GetUserUUIDFromRequest(r)
	
	// Assign task
	task, err := c.taskUseCase.AssignTask(r.Context(), taskUUID, assignedUserUUID, requestorUUID)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	
	utils.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"task": task})
}