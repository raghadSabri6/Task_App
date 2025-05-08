package entity

// UserTask represents the relationship between users and tasks
type UserTask struct {
	UserID int64
	TaskID int64
	
	// References to related entities
	User *User
	Task *Task
}