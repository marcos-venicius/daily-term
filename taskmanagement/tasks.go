package taskmanagement

import "time"

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
