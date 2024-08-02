package taskmanagement

import (
	"sort"
	"time"

	"github.com/marcos-venicius/daily-term/idcluster"
)

func CreateBoard() *Board {
	idCluster := idcluster.CreateIdCluster()

	return &Board{
		tasks:        []Task{},
		idCluster:    idCluster,
		SelectedTask: nil,
	}
}

func (board *Board) AddTask(name string) Task {
	task := Task{
		Id:        board.idCluster.NewId(),
		Name:      name,
		State:     Todo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	board.tasks = append(board.tasks, task)

	if len(board.tasks) == 1 {
		board.SelectedTask = &board.tasks[0]
	}

	return task
}

func (board *Board) Tasks() []Task {
	sort.SliceStable(board.tasks, func(i, j int) bool {
		return board.tasks[i].UpdatedAt.Unix() < board.tasks[j].UpdatedAt.Unix()
	})

	return board.tasks
}
