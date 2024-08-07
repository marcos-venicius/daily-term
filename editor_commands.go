package main

import (
	"github.com/marcos-venicius/daily-term/argumentparser"
	"github.com/marcos-venicius/daily-term/taskmanagement"
)

func (editor *Editor) Quit() {
	editor.Stop()
}

func (editor *Editor) ChangeCurrentTaskStateFor(state taskmanagement.TaskState) {
	var err error

	previousTaskState := editor.board.CurrentTask().State

	switch state {
	case taskmanagement.Todo:
		err = editor.board.MoveCurrentSelectedTaskToTodo()
		break
	case taskmanagement.InProgress:
		err = editor.board.MoveCurrentSelectedTaskToInProgress()
		break
	case taskmanagement.Completed:
		err = editor.board.MoveCurrentSelectedTaskToCompleted()
		break
	default:
		break
	}

	if !editor.setErrorMessageIfNNil(err) {
		if editor.setErrorMessageIfNNil(editor.repository.SaveBoard(editor.board)) {
			editor.setErrorMessageIfNNil(editor.board.SetCustomTaskState(previousTaskState)) // rollback
		}
	}
}

func (editor *Editor) addTask(arguments []argumentparser.CommandArgument) {
	var name = arguments[0].Value.(string)

	editor.board.AddTask(name)

	err := editor.repository.SaveBoard(editor.board)

	if err != nil {
		editor.board.DeleteCurrentSelectedTask()

		editor.SetErrorMessage(err.Error())
	} else {
		editor.SetInfoMessage("new task added successfully")
	}
}

func (editor *Editor) deleteTask(arguments []argumentparser.CommandArgument) {
	success := false

	if len(arguments) == 0 {
		err := editor.board.DeleteCurrentSelectedTask()

		if err == nil {
			success = true

			editor.SetInfoMessage("task deleted successfully")
		} else {
			editor.SetErrorMessage(err.Error())
		}
	} else {
		taskId := arguments[0].Value.(int)

		err := editor.board.DeleteTaskById(taskId)

		if err == nil {
			success = true

			editor.SetInfoMessage("task deleted successfully")
		} else {
			editor.SetErrorMessage(err.Error())
		}
	}

	if success {
		editor.setErrorMessageIfNNil(editor.repository.SaveBoard(editor.board))
	}
}
