package taskmanagement

import (
	"github.com/marcos-venicius/daily-term/idcluster"
)

// these are all task states
const (
	Todo       TaskState = iota
	InProgress TaskState = iota
	Completed  TaskState = iota
)

type TaskState int

type Task struct {
	Id    int       `json:"id"`
	Name  string    `json:"name"`
	State TaskState `json:"state"` // default is Todo
	Prev  *Task     `json:"prev"`  // previous task in the board
	Next  *Task     `json:"next"`  // next task in the board
}

type Board struct {
	task      *Task // current selected task
	root      *Task // tree root node
	idCluster *idcluster.IdCluster
}
