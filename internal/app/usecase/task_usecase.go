package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"task2/internal/adapter/presenter"
	"task2/internal/app/dto"
	"task2/internal/domain/entity"
	"task2/internal/domain/service"

	"github.com/google/uuid"
)

// TaskUseCase handles application logic for tasks
type TaskUseCase struct {
	taskService   *service.TaskService
	userService   *service.UserService
	taskPresenter *presenter.TaskPresenter
}

// NewTaskUseCase creates a new task use case
func NewTaskUseCase(taskService *service.TaskService, userService *service.UserService) *TaskUseCase {
	return &TaskUseCase{
		taskService:   taskService,
		userService:   userService,
		taskPresenter: presenter.NewTaskPresenter(),
	}
}

// CreateTask creates a new task
func (uc *TaskUseCase) CreateTask(ctx context.Context, req *dto.CreateTaskRequest, creatorUUID uuid.UUID) (*dto.TaskResponse, error) {
	// Create task entity
	task, err := entity.NewTask(req.Title, req.Description, creatorUUID)
	if err != nil {
		return nil, err
	}
	
	// Create task
	if err := uc.taskService.CreateTask(ctx, task); err != nil {
		return nil, err
	}
	
	// Assign users if provided
	if len(req.Users) > 0 {
		// First validate that all user IDs exist in the database
		var invalidUsers []string
		var validUsers []uuid.UUID
		
		for _, userAssign := range req.Users {
			// Parse the string UUID to uuid.UUID
			userUUID, err := uuid.Parse(userAssign.ID)
			if err != nil {
				invalidUsers = append(invalidUsers, userAssign.ID+" (invalid format)")
				continue
			}
			
			// Check if user exists in database
			_, err = uc.userService.GetUserByUUID(ctx, userUUID)
			if err != nil {
				invalidUsers = append(invalidUsers, userAssign.ID+" (not found)")
				continue
			}
			
			validUsers = append(validUsers, userUUID)
		}
		
		// If there are invalid users, return an error
		if len(invalidUsers) > 0 {
			return nil, errors.New("some users could not be assigned to the task: " + strings.Join(invalidUsers, ", "))
		}
		
		// Assign task to valid users
		for _, userUUID := range validUsers {
			if err := uc.taskService.AssignTask(ctx, task.UUID, userUUID, creatorUUID); err != nil {
				log.Printf("Failed to assign task to user %s: %v", userUUID, err)
				return nil, fmt.Errorf("failed to assign task to user %s: %w", userUUID, err)
			}
		}
	}
	
	// Get the created task with all relationships
	createdTask, err := uc.taskService.GetTaskByUUID(ctx, task.UUID)
	if err != nil {
		return nil, err
	}
	
	// Convert to DTO
	return uc.taskPresenter.ToDTO(createdTask), nil
}

// GetTaskByUUID gets a task by UUID
func (uc *TaskUseCase) GetTaskByUUID(ctx context.Context, taskUUID uuid.UUID) (*dto.TaskResponse, error) {
	// Get task
	task, err := uc.taskService.GetTaskByUUID(ctx, taskUUID)
	if err != nil {
		return nil, err
	}
	
	// Convert to DTO
	return uc.taskPresenter.ToDTO(task), nil
}

// GetAllTasks gets all tasks
func (uc *TaskUseCase) GetAllTasks(ctx context.Context) (*dto.TasksResponse, error) {
	// Get all tasks
	tasks, err := uc.taskService.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}
	
	// Convert to DTOs
	return uc.taskPresenter.ToDTOList(tasks), nil
}

// GetTasksCreatedByUser gets tasks created by a user
func (uc *TaskUseCase) GetTasksCreatedByUser(ctx context.Context, userUUID uuid.UUID) (*dto.TasksResponse, error) {
	// Get tasks created by user
	tasks, err := uc.taskService.GetTasksCreatedByUser(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	
	// Convert to DTOs
	return uc.taskPresenter.ToDTOList(tasks), nil
}

// GetTasksAssignedToUser gets tasks assigned to a user
func (uc *TaskUseCase) GetTasksAssignedToUser(ctx context.Context, userUUID uuid.UUID) (*dto.TasksResponse, error) {
	// Get tasks assigned to user
	tasks, err := uc.taskService.GetTasksAssignedToUser(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	
	// Convert to DTOs
	return uc.taskPresenter.ToDTOList(tasks), nil
}

// CompleteTask completes a task
func (uc *TaskUseCase) CompleteTask(ctx context.Context, taskUUID uuid.UUID, userUUID uuid.UUID) (*dto.TaskResponse, error) {
	// Complete the task
	if err := uc.taskService.CompleteTask(ctx, taskUUID, userUUID); err != nil {
		return nil, err
	}
	
	// Get the updated task
	task, err := uc.taskService.GetTaskByUUID(ctx, taskUUID)
	if err != nil {
		return nil, err
	}
	
	// Convert to DTO
	return uc.taskPresenter.ToDTO(task), nil
}

// DeleteTask deletes a task
func (uc *TaskUseCase) DeleteTask(ctx context.Context, taskUUID uuid.UUID, userUUID uuid.UUID) error {
	return uc.taskService.DeleteTask(ctx, taskUUID, userUUID)
}

// AssignTask assigns a task to a user
func (uc *TaskUseCase) AssignTask(ctx context.Context, taskUUID uuid.UUID, userUUID uuid.UUID, requestorUUID uuid.UUID) (*dto.TaskResponse, error) {
	// Assign the task
	if err := uc.taskService.AssignTask(ctx, taskUUID, userUUID, requestorUUID); err != nil {
		return nil, err
	}
	
	// Get the updated task
	task, err := uc.taskService.GetTaskByUUID(ctx, taskUUID)
	if err != nil {
		return nil, err
	}
	
	// Convert to DTO
	return uc.taskPresenter.ToDTO(task), nil
}