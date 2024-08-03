package taskmanagement

import (
	"errors"

	"github.com/marcos-venicius/daily-term/idcluster"
)

func CreateBoard() *Board {
	idCluster := idcluster.CreateIdCluster()

	return &Board{
		task:      nil,
		root:      nil,
		idCluster: idCluster,
	}
}

func (board *Board) MoveCurrentSelectedTaskToTodo() error {
	if board.task == nil {
		return errors.New("You have no selected task")
	}

	if board.task.State == Todo {
		return errors.New("This task is already in todo mode")
	}

	board.task.State = Todo

	return nil
}

func (board *Board) MoveCurrentSelectedTaskToInProgress() error {
	if board.task == nil {
		return errors.New("You have no selected task")
	}

	if board.task.State == InProgress {
		return errors.New("This task is already in progress mode")
	}

	board.task.State = InProgress

	return nil
}

func (board *Board) MoveCurrentSelectedTaskToCompleted() error {
	if board.task == nil {
		return errors.New("You have no selected task")
	}

	if board.task.State == Completed {
		return errors.New("This task is already completed")
	}

	board.task.State = Completed

	return nil
}

func (board *Board) HasTasks() bool {
	return board.root != nil
}

func (board *Board) SelectedTaskId() *int {
	if board.task != nil {
		return &board.task.Id
	}

	return nil
}

func (board *Board) SelectNextTask() {
	if board.task != nil && board.task.next != nil {
		board.task = board.task.next
	}
}

func (board *Board) SelectPreviousTask() {
	if board.task != nil && board.task.prev != nil {
		board.task = board.task.prev
	}
}

func (board *Board) DeleteCurrentSelectedTask() error {
	if board.root == nil {
		return errors.New("You don't have any selected task")
	}

	if board.task.next == nil && board.task.prev == nil {
		board.root = nil
		board.task = nil
		return nil
	} else if board.task.prev != nil && board.task.next != nil {
		board.task.next.prev, board.task.prev.next = board.task.prev, board.task.next
		board.task = board.task.next
		return nil
	} else if board.task.prev == nil && board.task.next != nil {
		board.task = board.task.next
		board.task.prev = nil
		board.root = board.task
		return nil
	} else if board.task.next == nil && board.task.prev != nil {
		board.task = board.task.prev
		board.task.next = nil
		board.root = board.task
		return nil
	}

	return errors.New("Something went wrong")
}

func (board *Board) DeleteTaskById(id int) error {
	if board.root == nil {
		return errors.New("You don't have any selected task")
	}

	current := board.root

	for current != nil {
		if current.Id == id {
			if current.next == nil && current.prev == nil {
				board.root = nil
				current = nil
				return nil
			} else if current.prev != nil && current.next != nil {
				current.next.prev, current.prev.next = current.prev, current.next
				current = current.next
				board.task = current
				return nil
			} else if current.prev == nil && current.next != nil {
				current = current.next
				current.prev = nil
				board.root = current
				board.task = current
				return nil
			} else if current.next == nil && current.prev != nil {
				current = current.prev
				current.next = nil
				board.root = current
				board.task = current
				return nil
			}

			return errors.New("Something went wrong")
		}

		current = current.next
	}

	return errors.New("Task not found")
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
	if board.task == nil {
		return []Task{}
	}

	var tasks []Task

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
