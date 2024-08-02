package taskmanagement

import (
	"sort"
	"time"
)

func CreateBoard() *Board {
	return &Board{}
}

func (board *Board) AddTask(name string) Task {
	task := Task{
		Name:      name,
		State:     Todo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	board.tasks = append(board.tasks, task)

	return task
}

func (board *Board) Tasks() []Task {
	sort.SliceStable(board.tasks, func(i, j int) bool {
		return board.tasks[i].UpdatedAt.Unix() < board.tasks[j].UpdatedAt.Unix()
	})

	return board.tasks
}
