package service

import (
	"context"
	"errors"
	"task2/internal/domain/entity"
	"task2/internal/domain/repository"

	"github.com/google/uuid"
)

// TaskService provides domain logic for tasks
type TaskService struct {
	taskRepo repository.TaskRepository
	userRepo repository.UserRepository
}

// NewTaskService creates a new task service
func NewTaskService(taskRepo repository.TaskRepository, userRepo repository.UserRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(ctx context.Context, task *entity.Task) error {
	// Validate creator exists
	creator, err := s.userRepo.GetByUUID(ctx, task.CreatedByID)
	if err != nil {
		return errors.New("creator not found")
	}
	
	task.CreatedBy = creator
	return s.taskRepo.Create(ctx, task)
}

// GetTaskByUUID gets a task by UUID
func (s *TaskService) GetTaskByUUID(ctx context.Context, taskUUID uuid.UUID) (*entity.Task, error) {
	return s.taskRepo.GetByUUID(ctx, taskUUID)
}

// GetAllTasks gets all tasks
func (s *TaskService) GetAllTasks(ctx context.Context) ([]*entity.Task, error) {
	return s.taskRepo.GetAll(ctx)
}

// GetTasksCreatedByUser gets tasks created by a user
func (s *TaskService) GetTasksCreatedByUser(ctx context.Context, userUUID uuid.UUID) ([]*entity.Task, error) {
	return s.taskRepo.GetTasksCreatedByUser(ctx, userUUID)
}

// GetTasksAssignedToUser gets tasks assigned to a user
func (s *TaskService) GetTasksAssignedToUser(ctx context.Context, userUUID uuid.UUID) ([]*entity.Task, error) {
	return s.taskRepo.GetTasksAssignedToUser(ctx, userUUID)
}

// AssignTask assigns a task to a user
func (s *TaskService) AssignTask(ctx context.Context, taskUUID uuid.UUID, userUUID uuid.UUID, requestorUUID uuid.UUID) error {
	// Get the task
	task, err := s.taskRepo.GetByUUID(ctx, taskUUID)
	if err != nil {
		return errors.New("task not found")
	}
	
	// Check if requestor is authorized to assign the task
	if task.CreatedByID != requestorUUID {
		return errors.New("only the task creator can assign users")
	}
	
	// Check if user exists
	_, err = s.userRepo.GetByUUID(ctx, userUUID)
	if err != nil {
		return errors.New("user not found")
	}
	
	// Assign the task
	return s.taskRepo.AssignTaskToUser(ctx, taskUUID, userUUID)
}

// CompleteTask marks a task as completed
func (s *TaskService) CompleteTask(ctx context.Context, taskUUID uuid.UUID, userUUID uuid.UUID) error {
	// Get the task
	task, err := s.taskRepo.GetByUUID(ctx, taskUUID)
	if err != nil {
		return errors.New("task not found")
	}
	
	// Check if user is authorized to complete the task
	if !task.CanBeModifiedBy(userUUID) {
		return errors.New("you are not authorized to complete this task")
	}
	
	// Complete the task
	if err := task.Complete(); err != nil {
		return err
	}
	
	return s.taskRepo.Update(ctx, task)
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(ctx context.Context, taskUUID uuid.UUID, userUUID uuid.UUID) error {
	// Get the task
	task, err := s.taskRepo.GetByUUID(ctx, taskUUID)
	if err != nil {
		return errors.New("task not found")
	}
	
	// Check if user is authorized to delete the task
	if task.CreatedByID != userUUID {
		return errors.New("only the task creator can delete the task")
	}
	
	// Delete the task
	return s.taskRepo.Delete(ctx, taskUUID)
}