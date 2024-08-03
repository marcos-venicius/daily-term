package taskmanagement

import (
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
	Id    int
	Name  string
	State TaskState // default is Todo

	prev *Task // previous task in the board
	next *Task // next task in the board
}

type Board struct {
	task      *Task // current selected task
	root      *Task // tree root node
	idCluster *idcluster.IdCluster
}
