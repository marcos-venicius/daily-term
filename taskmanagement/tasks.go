package taskmanagement

import "github.com/marcos-venicius/daily-term/idcluster"

func CreateBoard() *Board {
	idCluster := idcluster.CreateIdCluster()

	return &Board{
		task:      nil,
		size:      0,
		idCluster: idCluster,
	}
}

func (board *Board) SelectedTaskId() *int {
	return &board.task.Id
}

func (board *Board) SelectNextTask() {
	if board.task.next != nil {
		board.task = board.task.next
	}
}

func (board *Board) SelectPreviousTask() {
	if board.task.prev != nil {
		board.task = board.task.prev
	}
}

func (board *Board) AddTask(name string) Task {
	task := Task{
		Id:    board.idCluster.NewId(),
		Name:  name,
		State: Todo,
		prev:  nil,
		next:  nil,
	}

	if board.task == nil {
		board.task = &task
		board.root = board.task
	} else {
		task.next = board.root
		board.root.prev = &task
		board.root = &task
		board.task = board.root
	}

	return task
}

func (board *Board) Tasks() []Task {
	var tasks []Task = make([]Task, board.size)

	current := board.root

	for current.prev != nil {
		current = current.prev
	}

	for current != nil {
		tasks = append(tasks, *current)

		current = current.next
	}

	return tasks
}
