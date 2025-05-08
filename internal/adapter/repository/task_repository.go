package repository

import (
	"context"
	"errors"
	"task2/internal/domain/entity"
	"task2/internal/infrastructure/persistence"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// TaskRepository implements the domain.TaskRepository interface
type TaskRepository struct {
	db *bun.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *bun.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

// Create creates a new task
func (r *TaskRepository) Create(ctx context.Context, task *entity.Task) error {
	// Convert domain entity to persistence model
	dbTask := &persistence.Task{
		UUID:        task.UUID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedByID: task.CreatedByID,
	}

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert task
	if _, err := tx.NewInsert().Model(dbTask).Exec(ctx); err != nil {
		return err
	}

	// Update task ID
	task.ID = dbTask.ID

	// Add users if provided
	if task.Users != nil && len(task.Users) > 0 {
		for _, user := range task.Users {
			// Get user ID from UUID
			var dbUser persistence.User
			err := tx.NewSelect().Model(&dbUser).Where("uuid = ?", user.UUID).Scan(ctx)
			if err != nil {
				return err
			}

			// Create user-task relationship
			userTask := &persistence.UserTask{
				TaskID: dbTask.ID,
				UserID: dbUser.ID,
			}

			if _, err := tx.NewInsert().Model(userTask).Exec(ctx); err != nil {
				return err
			}
		}

		// Set assigned user if provided
		if task.AssignedToID != nil {
			dbTask.AssignedToID = task.AssignedToID
			if _, err := tx.NewUpdate().Model(dbTask).Column("assigned_to_id").WherePK().Exec(ctx); err != nil {
				return err
			}
		}
	}

	// Commit transaction
	return tx.Commit()
}

// GetByUUID gets a task by UUID
func (r *TaskRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Task, error) {
	dbTask := new(persistence.Task)

	// Get task with relationships
	err := r.db.NewSelect().
		Model(dbTask).
		Relation("Users").
		Relation("CreatedBy").
		Relation("AssignedTo").
		Where("task.uuid = ?", uuid).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	// Convert to domain entity
	task := &entity.Task{
		ID:           dbTask.ID,
		UUID:         dbTask.UUID,
		Title:        dbTask.Title,
		Description:  dbTask.Description,
		Completed:    dbTask.Completed,
		CreatedAt:    dbTask.CreatedAt,
		UpdatedAt:    dbTask.UpdatedAt,
		DeletedAt:    dbTask.DeletedAt,
		CreatedByID:  dbTask.CreatedByID,
		AssignedToID: dbTask.AssignedToID,
	}

	// Convert relationships
	if dbTask.CreatedBy != nil {
		task.CreatedBy = &entity.User{
			ID:    dbTask.CreatedBy.ID,
			UUID:  dbTask.CreatedBy.UUID,
			Name:  dbTask.CreatedBy.Name,
			Email: dbTask.CreatedBy.Email,
		}
	}

	if dbTask.AssignedTo != nil {
		task.AssignedTo = &entity.User{
			ID:    dbTask.AssignedTo.ID,
			UUID:  dbTask.AssignedTo.UUID,
			Name:  dbTask.AssignedTo.Name,
			Email: dbTask.AssignedTo.Email,
		}
	}

	if dbTask.Users != nil {
		task.Users = make([]*entity.User, len(dbTask.Users))
		for i, user := range dbTask.Users {
			task.Users[i] = &entity.User{
				ID:    user.ID,
				UUID:  user.UUID,
				Name:  user.Name,
				Email: user.Email,
			}
		}
	}

	return task, nil
}

// GetAll gets all tasks
func (r *TaskRepository) GetAll(ctx context.Context) ([]*entity.Task, error) {
	var dbTasks []persistence.Task

	// Get all tasks with relationships
	err := r.db.NewSelect().
		Model(&dbTasks).
		Relation("Users").
		Relation("CreatedBy").
		Relation("AssignedTo").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	// Convert to domain entities
	tasks := make([]*entity.Task, len(dbTasks))
	for i, dbTask := range dbTasks {
		tasks[i] = &entity.Task{
			ID:           dbTask.ID,
			UUID:         dbTask.UUID,
			Title:        dbTask.Title,
			Description:  dbTask.Description,
			Completed:    dbTask.Completed,
			CreatedAt:    dbTask.CreatedAt,
			UpdatedAt:    dbTask.UpdatedAt,
			DeletedAt:    dbTask.DeletedAt,
			CreatedByID:  dbTask.CreatedByID,
			AssignedToID: dbTask.AssignedToID,
		}

		// Convert relationships
		if dbTask.CreatedBy != nil {
			tasks[i].CreatedBy = &entity.User{
				ID:    dbTask.CreatedBy.ID,
				UUID:  dbTask.CreatedBy.UUID,
				Name:  dbTask.CreatedBy.Name,
				Email: dbTask.CreatedBy.Email,
			}
		}

		if dbTask.AssignedTo != nil {
			tasks[i].AssignedTo = &entity.User{
				ID:    dbTask.AssignedTo.ID,
				UUID:  dbTask.AssignedTo.UUID,
				Name:  dbTask.AssignedTo.Name,
				Email: dbTask.AssignedTo.Email,
			}
		}

		if dbTask.Users != nil {
			tasks[i].Users = make([]*entity.User, len(dbTask.Users))
			for j, user := range dbTask.Users {
				tasks[i].Users[j] = &entity.User{
					ID:    user.ID,
					UUID:  user.UUID,
					Name:  user.Name,
					Email: user.Email,
				}
			}
		}
	}

	return tasks, nil
}

// Update updates a task
func (r *TaskRepository) Update(ctx context.Context, task *entity.Task) error {
	// Convert domain entity to persistence model
	dbTask := &persistence.Task{
		ID:           task.ID,
		UUID:         task.UUID,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		UpdatedAt:    task.UpdatedAt,
		AssignedToID: task.AssignedToID,
	}

	// Update task
	_, err := r.db.NewUpdate().
		Model(dbTask).
		Column("title", "description", "completed", "updated_at", "assigned_to_id").
		WherePK().
		Exec(ctx)

	return err
}

// Delete deletes a task
func (r *TaskRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	// Get task
	dbTask := new(persistence.Task)
	err := r.db.NewSelect().
		Model(dbTask).
		Where("uuid = ?", uuid).
		Scan(ctx)

	if err != nil {
		return err
	}

	// Delete task
	_, err = r.db.NewDelete().
		Model(dbTask).
		WherePK().
		Exec(ctx)

	return err
}

// GetTasksCreatedByUser gets tasks created by a user
func (r *TaskRepository) GetTasksCreatedByUser(ctx context.Context, userUUID uuid.UUID) ([]*entity.Task, error) {
	var dbTasks []persistence.Task

	// Get tasks created by user
	err := r.db.NewSelect().
		Model(&dbTasks).
		Where("created_by_id = ?", userUUID).
		Relation("Users").
		Relation("CreatedBy").
		Relation("AssignedTo").
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	// Convert to domain entities
	tasks := make([]*entity.Task, len(dbTasks))
	for i, dbTask := range dbTasks {
		tasks[i] = &entity.Task{
			ID:           dbTask.ID,
			UUID:         dbTask.UUID,
			Title:        dbTask.Title,
			Description:  dbTask.Description,
			Completed:    dbTask.Completed,
			CreatedAt:    dbTask.CreatedAt,
			UpdatedAt:    dbTask.UpdatedAt,
			DeletedAt:    dbTask.DeletedAt,
			CreatedByID:  dbTask.CreatedByID,
			AssignedToID: dbTask.AssignedToID,
		}

		if dbTask.CreatedBy != nil {
			tasks[i].CreatedBy = &entity.User{
				ID:    dbTask.CreatedBy.ID,
				UUID:  dbTask.CreatedBy.UUID,
				Name:  dbTask.CreatedBy.Name,
				Email: dbTask.CreatedBy.Email,
			}
		}

		if dbTask.AssignedTo != nil {
			tasks[i].AssignedTo = &entity.User{
				ID:    dbTask.AssignedTo.ID,
				UUID:  dbTask.AssignedTo.UUID,
				Name:  dbTask.AssignedTo.Name,
				Email: dbTask.AssignedTo.Email,
			}
		}

		if dbTask.Users != nil {
			tasks[i].Users = make([]*entity.User, len(dbTask.Users))
			for j, user := range dbTask.Users {
				tasks[i].Users[j] = &entity.User{
					ID:    user.ID,
					UUID:  user.UUID,
					Name:  user.Name,
					Email: user.Email,
				}
			}
		}
	}

	return tasks, nil
}

// GetTasksAssignedToUser gets tasks assigned to a user
func (r *TaskRepository) GetTasksAssignedToUser(ctx context.Context, userUUID uuid.UUID) ([]*entity.Task, error) {
	var dbTasks []persistence.Task

	// Get tasks assigned to user
	err := r.db.NewSelect().
		Model(&dbTasks).
		Where("assigned_to_id = ?", userUUID).
		Relation("Users").
		Relation("CreatedBy").
		Relation("AssignedTo").
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	// Convert to domain entities
	tasks := make([]*entity.Task, len(dbTasks))
	for i, dbTask := range dbTasks {
		tasks[i] = &entity.Task{
			ID:           dbTask.ID,
			UUID:         dbTask.UUID,
			Title:        dbTask.Title,
			Description:  dbTask.Description,
			Completed:    dbTask.Completed,
			CreatedAt:    dbTask.CreatedAt,
			UpdatedAt:    dbTask.UpdatedAt,
			DeletedAt:    dbTask.DeletedAt,
			CreatedByID:  dbTask.CreatedByID,
			AssignedToID: dbTask.AssignedToID,
		}

		if dbTask.CreatedBy != nil {
			tasks[i].CreatedBy = &entity.User{
				ID:    dbTask.CreatedBy.ID,
				UUID:  dbTask.CreatedBy.UUID,
				Name:  dbTask.CreatedBy.Name,
				Email: dbTask.CreatedBy.Email,
			}
		}

		if dbTask.AssignedTo != nil {
			tasks[i].AssignedTo = &entity.User{
				ID:    dbTask.AssignedTo.ID,
				UUID:  dbTask.AssignedTo.UUID,
				Name:  dbTask.AssignedTo.Name,
				Email: dbTask.AssignedTo.Email,
			}
		}

		if dbTask.Users != nil {
			tasks[i].Users = make([]*entity.User, len(dbTask.Users))
			for j, user := range dbTask.Users {
				tasks[i].Users[j] = &entity.User{
					ID:    user.ID,
					UUID:  user.UUID,
					Name:  user.Name,
					Email: user.Email,
				}
			}
		}
	}

	return tasks, nil
}

// AssignTaskToUser assigns a task to a user
func (r *TaskRepository) AssignTaskToUser(ctx context.Context, taskUUID uuid.UUID, userUUID uuid.UUID) error {
	// Get task
	dbTask := new(persistence.Task)
	err := r.db.NewSelect().
		Model(dbTask).
		Where("uuid = ?", taskUUID).
		Scan(ctx)

	if err != nil {
		return err
	}

	// Get user
	dbUser := new(persistence.User)
	err = r.db.NewSelect().
		Model(dbUser).
		Where("uuid = ?", userUUID).
		Scan(ctx)

	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if user is already assigned
	exists, err := tx.NewSelect().
		Model((*persistence.UserTask)(nil)).
		Where("task_id = ? AND user_id = ?", dbTask.ID, dbUser.ID).
		Exists(ctx)

	if err != nil {
		return err
	}

	// Add user to task if not already assigned
	if !exists {
		userTask := &persistence.UserTask{
			TaskID: dbTask.ID,
			UserID: dbUser.ID,
		}

		if _, err := tx.NewInsert().Model(userTask).Exec(ctx); err != nil {
			return err
		}
	}

	// Update assigned user
	dbTask.AssignedToID = &userUUID
	if _, err := tx.NewUpdate().Model(dbTask).Column("assigned_to_id").WherePK().Exec(ctx); err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// CompleteTask completes a task
func (r *TaskRepository) CompleteTask(ctx context.Context, taskUUID uuid.UUID) error {
	// Get task
	dbTask := new(persistence.Task)
	err := r.db.NewSelect().
		Model(dbTask).
		Where("uuid = ?", taskUUID).
		Scan(ctx)

	if err != nil {
		return err
	}

	// Check if already completed
	if dbTask.Completed {
		return errors.New("task is already completed")
	}

	// Update task
	dbTask.Completed = true
	_, err = r.db.NewUpdate().
		Model(dbTask).
		Column("completed").
		WherePK().
		Exec(ctx)

	return err
}

// AddUserToTask adds a user to a task
func (r *TaskRepository) AddUserToTask(ctx context.Context, taskUUID uuid.UUID, userUUID uuid.UUID) error {
	// Get task
	dbTask := new(persistence.Task)
	err := r.db.NewSelect().
		Model(dbTask).
		Where("uuid = ?", taskUUID).
		Scan(ctx)

	if err != nil {
		return err
	}

	// Get user
	dbUser := new(persistence.User)
	err = r.db.NewSelect().
		Model(dbUser).
		Where("uuid = ?", userUUID).
		Scan(ctx)

	if err != nil {
		return err
	}

	// Check if user is already assigned
	exists, err := r.db.NewSelect().
		Model((*persistence.UserTask)(nil)).
		Where("task_id = ? AND user_id = ?", dbTask.ID, dbUser.ID).
		Exists(ctx)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("user is already assigned to this task")
	}

	// Add user to task
	userTask := &persistence.UserTask{
		TaskID: dbTask.ID,
		UserID: dbUser.ID,
	}

	_, err = r.db.NewInsert().Model(userTask).Exec(ctx)
	return err
}
