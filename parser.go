package main

import "github.com/marcos-venicius/daily-term/argumentparser"

func (editor *Editor) InitParser() {
	newTaskArguments := []argumentparser.CommandArgumentSyntax{
		{
			Name:     "Task name (string)",
			Required: true,
			Type:     argumentparser.StringArgumentType,
		},
	}

	deletetaskArguments := []argumentparser.CommandArgumentSyntax{
		{
			Name:     "Task id (int)",
			Required: false,
			Type:     argumentparser.IntArgumentType,
		},
	}

	editor.argumentParser.AddCommand("q")
	editor.argumentParser.AddCommand("quit")

	editor.argumentParser.AddCommand("nt", newTaskArguments...)
	editor.argumentParser.AddCommand("new task", newTaskArguments...)

	editor.argumentParser.AddCommand("dt", deletetaskArguments...)
	editor.argumentParser.AddCommand("delete task", deletetaskArguments...)

	editor.argumentParser.Finish()
}
