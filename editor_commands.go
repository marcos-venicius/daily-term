package main

import argumentparser "github.com/marcos-venicius/daily-term/argument-parser"

func (editor *Editor) Quit() {
	editor.Stop()
}

func (editor *Editor) addTask(arguments []argumentparser.CommandArgument) {
	var name = arguments[0].Value.(string)

	editor.board.AddTask(name)

	editor.SetInfoMessage("new task added successfully")
}

func (editor *Editor) deleteTask(arguments []argumentparser.CommandArgument) {
	if len(arguments) == 0 {
		err := editor.board.DeleteCurrentSelectedTask()

		if err == nil {
			editor.SetInfoMessage("task deleted successfully")
		} else {
			editor.SetErrorMessage(err.Error())
		}
	} else {
		taskId := arguments[0].Value.(int)

		err := editor.board.DeleteTaskById(taskId)

		if err == nil {
			editor.SetInfoMessage("task deleted successfully")
		} else {
			editor.SetErrorMessage(err.Error())
		}
	}
}
