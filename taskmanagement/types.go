package taskmanagement

import "time"

// these are all task states
const (
	Todo       TaskState = iota
	InProgress TaskState = iota
	Done       TaskState = iota
)

type TaskState int

type Task struct {
	Name      string
	State     TaskState // default is Todo
	CreatedAt time.Time
	UpdatedAt time.Time // updated automatically every time the state or name changes
}

type Board struct {
	tasks []Task
}
