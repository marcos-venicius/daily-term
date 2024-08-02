package taskmanagement

import (
	"time"

	"github.com/marcos-venicius/daily-term/idcluster"
)

// these are all task states
const (
	Todo       TaskState = iota
	InProgress TaskState = iota
	Done       TaskState = iota
)

type TaskState int

type Task struct {
	Id        int
	Name      string
	State     TaskState // default is Todo
	CreatedAt time.Time
	UpdatedAt time.Time // updated automatically every time the state or name changes
}

type Board struct {
	tasks        []Task
	idCluster    *idcluster.IdCluster
	SelectedTask *Task // current selected task
}
