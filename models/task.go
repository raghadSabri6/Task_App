package models

import (
	"context"
	"errors"
	"fmt"
	"task2/database"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Task struct {
	bun.BaseModel `bun:"table:tasks"`

	ID          int64      `bun:",pk,autoincrement"`
	UUID        uuid.UUID  `bun:",type:uuid,default:uuid_generate_v4()" json:"id"`
	Title       string     `bun:",notnull" json:"title"`
	Description string     `json:"description"`
	Completed   bool       `bun:",default:false"`
	CreatedAt   time.Time  `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time  `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt   *time.Time `bun:",soft_delete" json:"deleted_at,omitempty"`

	CreatedByID uuid.UUID `bun:",type:uuid,notnull"`
	CreatedBy   *User     `bun:"rel:belongs-to,join:created_by_id=uuid"`

	AssignedToID *uuid.UUID `bun:",type:uuid"`
	AssignedTo   *User      `bun:"rel:belongs-to,join:assigned_to_id=uuid"`

	Users []*User `bun:"m2m:user_tasks" json:"users,omitempty"`
}

func FetchTaskByUUID(taskUUID uuid.UUID) (*Task, error) {
	task := new(Task)
	err := database.DB.NewSelect().Model(task).
		Relation("Users").
		Relation("CreatedBy").
		Relation("AssignedTo").
		Where("task.uuid = ?", taskUUID).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return task, nil
}

func GetTasks() ([]Task, error) {
	var tasks []Task
	err := database.DB.NewSelect().Model(&tasks).
		Relation("Users").
		Relation("CreatedBy").
		Relation("AssignedTo").
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (task *Task) updateTask() error {
	_, err := database.DB.NewUpdate().Model(task).
		WherePK().
		Exec(context.Background())
	return err
}

func (task *Task) Complete(userUUID uuid.UUID) error {
	if task.Completed {
		return errors.New("task is already completed")
	}

	err := database.DB.NewSelect().Model(task).
		Relation("Users").
		Where("task.uuid = ?", task.UUID).
		Scan(context.Background())
	if err != nil {
		return errors.New("task not found")
	}

	if task.CreatedByID != userUUID {
		ctx := context.Background()
		count, err := database.DB.NewSelect().
			TableExpr("user_tasks ut").
			Join("users u ON u.id = ut.user_id").
			Where("ut.task_id = ? AND u.uuid = ?", task.ID, userUUID).
			Count(ctx)

		if err != nil || count == 0 {
			return errors.New("you are not authorized to complete this task")
		}
	}

	task.Completed = true
	return task.updateTask()
}

func (task *Task) ValidateAndCreateTask(userUUID uuid.UUID) error {
	if task.Title == "" {
		return errors.New("task title is required")
	}

	task.UUID = uuid.New()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	ctx := context.Background()

	var creatorUser User
	err := database.DB.NewSelect().
		Model(&creatorUser).
		Where("uuid = ?", userUUID).
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("creator user not found: %w", err)
	}

	task.CreatedByID = userUUID

	tx, err := database.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NewInsert().Model(task).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	if task.Users != nil && len(task.Users) > 0 {
		for _, user := range task.Users {
			var dbUser User
			err := tx.NewSelect().Model(&dbUser).Where("uuid = ?", user.UUID).Scan(ctx)
			if err != nil {
				return fmt.Errorf("user with UUID %s not found: %w", user.UUID, err)
			}

			userTask := &UserTask{
				TaskID: task.ID,
				UserID: dbUser.ID,
			}

			_, err = tx.NewInsert().Model(userTask).Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to insert user_task relation: %w", err)
			}
		}
	}

	return tx.Commit()
}

func (task *Task) DeleteTask(userUUID uuid.UUID) error {
	if task.CreatedByID != userUUID {
		return errors.New("unauthorized: you can only delete your own tasks")
	}

	_, err := database.DB.NewDelete().Model(task).WherePK().Exec(context.Background())
	return err
}

func (task *Task) CreateTask(ctx context.Context, db *bun.DB) error {
	_, err := db.NewInsert().
		Model(task).
		Exec(ctx)
	return err
}

func (task *Task) FetchCreatedTask(ctx context.Context, db *bun.DB) (*Task, error) {
	var createdTask Task
	err := db.NewSelect().
		Model(&createdTask).
		Relation("Users").
		Relation("CreatedBy").
		Relation("AssignedTo").
		Where("task.uuid = ?", task.UUID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &createdTask, nil
}

func (task *Task) ValidateAndAssignTask(ctx context.Context, db *bun.DB, creatorUUID, assignedUserUUID uuid.UUID) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}

	if task.CreatedByID != creatorUUID {
		return errors.New("unauthorized: only the task creator can assign users")
	}

	alreadyAssigned, err := db.NewSelect().
		Model((*UserTask)(nil)).
		Where("task_id = ?", task.ID).
		Where("user_id IN (SELECT id FROM users WHERE uuid = ?)", assignedUserUUID).
		Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check existing assignment: %w", err)
	}
	if alreadyAssigned {
		return errors.New("user is already assigned to this task")
	}

	var assignedUser User
	err = db.NewSelect().
		Model(&assignedUser).
		Where("uuid = ?", assignedUserUUID).
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("assigned user not found: %w", err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	userTask := &UserTask{
		TaskID: task.ID,
		UserID: assignedUser.ID,
	}
	if _, err := tx.NewInsert().Model(userTask).Exec(ctx); err != nil {
		return fmt.Errorf("failed to insert into user_tasks: %w", err)
	}

	assignedID := assignedUserUUID
	task.AssignedToID = &assignedID

	if _, err := tx.NewUpdate().Model(task).Column("assigned_to_id").WherePK().Exec(ctx); err != nil {
		return fmt.Errorf("failed to update task with assignee: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func GetUserTasks(userUUID uuid.UUID) ([]Task, error) {
	if userUUID == uuid.Nil {
		return nil, fmt.Errorf("invalid user UUID")
	}

	var tasks []Task
	ctx := context.Background()

	err := database.DB.NewSelect().
		Model(&tasks).
		Where("created_by_id = ?", userUUID).
		Relation("Users").
		Relation("CreatedBy").
		Relation("AssignedTo").
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks created by user: %w", err)
	}

	if tasks == nil {
		return []Task{}, nil
	}

	return tasks, nil
}

func GetTasksCreatedByUser(userUUID uuid.UUID) ([]Task, error) {
	var tasks []Task
	ctx := context.Background()

	err := database.DB.NewSelect().
		Model(&tasks).
		Where("created_by_id = ?", userUUID).
		Relation("Users").
		Relation("CreatedBy").
		Relation("AssignedTo").
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks created by user: %w", err)
	}

	return tasks, nil
}

func GetTasksAssignedToUser(userUUID uuid.UUID) ([]Task, error) {
	if userUUID == uuid.Nil {
		return nil, fmt.Errorf("invalid user UUID")
	}

	var tasks []Task
	ctx := context.Background()

	err := database.DB.NewSelect().
		Model(&tasks).
		Where("assigned_to_id = ?", userUUID).
		Relation("Users").
		Relation("CreatedBy").
		Relation("AssignedTo").
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks assigned to user: %w", err)
	}

	if tasks == nil {
		return []Task{}, nil
	}

	return tasks, nil
}
