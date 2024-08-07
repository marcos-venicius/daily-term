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

func (board *Board) SetCustomTaskState(state TaskState) error {
	if board.task == nil {
		return errors.New("You have no selected task")
	}

	if board.task.State == state {
		switch state {
		case Todo:
			return errors.New("This task is already todo")
		case InProgress:
			return errors.New("This task is already in progress")
		case Completed:
			return errors.New("This task is already completed")
		default:
			return errors.New("This task is already in unknown state")
		}
	}

	board.task.State = state

	return nil
}

func (board *Board) CurrentTask() *Task {
	return board.task
}

func (board *Board) MoveCurrentSelectedTaskToTodo() error {
	if board.task == nil {
		return errors.New("You have no selected task")
	}

	if board.task.State == Todo {
		return errors.New("This task is already todo")
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
	if board.task != nil && board.task.Next != nil {
		board.task = board.task.Next
	}
}

func (board *Board) SelectPreviousTask() {
	if board.task != nil && board.task.Prev != nil {
		board.task = board.task.Prev
	}
}

func (board *Board) DeleteCurrentSelectedTask() error {
	if board.root == nil {
		return errors.New("You don't have any selected task")
	}

	if board.task.Next == nil && board.task.Prev == nil {
		board.root = nil
		board.task = nil
		return nil
	} else if board.task.Prev != nil && board.task.Next != nil {
		board.task.Next.Prev, board.task.Prev.Next = board.task.Prev, board.task.Next
		board.task = board.task.Next
		return nil
	} else if board.task.Prev == nil && board.task.Next != nil {
		board.task = board.task.Next
		board.task.Prev = nil
		board.root = board.task
		return nil
	} else if board.task.Next == nil && board.task.Prev != nil {
		board.task = board.task.Prev
		board.task.Next = nil
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
			if current.Next == nil && current.Prev == nil {
				board.root = nil
				current = nil
				return nil
			} else if current.Prev != nil && current.Next != nil {
				current.Next.Prev, current.Prev.Next = current.Prev, current.Next
				current = current.Next
				board.task = current
				return nil
			} else if current.Prev == nil && current.Next != nil {
				current = current.Next
				current.Prev = nil
				board.root = current
				board.task = current
				return nil
			} else if current.Next == nil && current.Prev != nil {
				current = current.Prev
				current.Next = nil
				board.root = current
				board.task = current
				return nil
			}

			return errors.New("Something went wrong")
		}

		current = current.Next
	}

	return errors.New("Task not found")
}

func (board *Board) AddTask(name string) Task {
	task := Task{
		Id:    board.idCluster.NewId(),
		Name:  name,
		State: Todo,
		Prev:  nil,
		Next:  nil,
	}

	if board.task == nil {
		board.task = &task
		board.root = board.task
	} else {
		task.Next = board.root
		board.root.Prev = &task
		board.root = &task
		board.task = board.root
	}

	if board.root.Prev != nil {
		panic("Invalid root")
	}

	return task
}

func (board *Board) Tasks() []Task {
	if board.task == nil {
		return []Task{}
	}

	var tasks []Task

	current := board.root

	for current.Prev != nil {
		current = current.Prev
	}

	for current != nil {
		tasks = append(tasks, *current)

		current = current.Next
	}

	return tasks
}
