package main

import argumentparser "github.com/marcos-venicius/daily-term/argument-parser"

func (editor *Editor) addTask(arguments []argumentparser.CommandArgument) {
	var name = arguments[0].Value.(string)

	editor.board.AddTask(name)

	editor.SetInfoMessage("new task added successfully")
}
