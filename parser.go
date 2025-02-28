package main

import (
	"log"

	"github.com/marcos-venicius/daily-term/argumentparser"
)

func printErrorIfExists(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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

	printErrorIfExists(editor.argumentParser.AddCommand("q"))
	printErrorIfExists(editor.argumentParser.AddCommand("quit"))
	printErrorIfExists(editor.argumentParser.AddCommand("w"))
	printErrorIfExists(editor.argumentParser.AddCommand("wa"))

	printErrorIfExists(editor.argumentParser.AddCommand("nt", newTaskArguments...))
	printErrorIfExists(editor.argumentParser.AddCommand("new task", newTaskArguments...))

	printErrorIfExists(editor.argumentParser.AddCommand("dt", deletetaskArguments...))
	printErrorIfExists(editor.argumentParser.AddCommand("delete task", deletetaskArguments...))

	printErrorIfExists(editor.argumentParser.Finish())
}
