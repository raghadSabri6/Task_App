package presenter

import (
	"task2/internal/app/dto"
	"task2/internal/domain/entity"
)

// TaskPresenter converts between domain entities and DTOs
type TaskPresenter struct{}

// NewTaskPresenter creates a new task presenter
func NewTaskPresenter() *TaskPresenter {
	return &TaskPresenter{}
}

// ToDTO converts a task entity to a DTO
func (p *TaskPresenter) ToDTO(task *entity.Task) *dto.TaskResponse {
	if task == nil {
		return nil
	}
	
	// Create task response
	taskResponse := &dto.TaskResponse{
		ID:          task.UUID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		DeletedAt:   task.DeletedAt,
	}
	
	// Add created by
	if task.CreatedBy != nil {
		taskResponse.CreatedBy = dto.UserSummary{
			ID:    task.CreatedBy.UUID,
			Name:  task.CreatedBy.Name,
			Email: task.CreatedBy.Email,
		}
	}
	
	// Add assigned to
	if task.AssignedTo != nil {
		taskResponse.AssignedTo = &dto.UserSummary{
			ID:    task.AssignedTo.UUID,
			Name:  task.AssignedTo.Name,
			Email: task.AssignedTo.Email,
		}
	}
	
	// Add users
	if task.Users != nil {
		taskResponse.Users = make([]dto.UserSummary, len(task.Users))
		for i, user := range task.Users {
			taskResponse.Users[i] = dto.UserSummary{
				ID:    user.UUID,
				Name:  user.Name,
				Email: user.Email,
			}
		}
	}
	
	return taskResponse
}

// ToDTOList converts a list of task entities to DTOs
func (p *TaskPresenter) ToDTOList(tasks []*entity.Task) *dto.TasksResponse {
	if tasks == nil {
		return &dto.TasksResponse{
			Tasks: []dto.TaskResponse{},
		}
	}
	
	// Create task responses
	taskResponses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponse := p.ToDTO(task)
		if taskResponse != nil {
			taskResponses[i] = *taskResponse
		}
	}
	
	return &dto.TasksResponse{
		Tasks: taskResponses,
	}
}