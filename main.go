package main

import (
	"log"

	"github.com/marcos-venicius/daily-term/argument-parser"
	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()

	if err != nil {
		log.Fatal(err)
	}

	defer termbox.Close()

	editor := CreateEditor()

	defer close(editor.termbox_event)

	editor.argumentParser.AddCommand("quit")
	editor.argumentParser.AddCommand("q")
	editor.argumentParser.AddCommand(
		"new task",
		argumentparser.CommandArgumentSyntax{
			Name:     "Task name (string)",
			Required: true,
			Type:     argumentparser.StringArgumentType,
		},
	)

	editor.argumentParser.Finish()

	termbox.Flush()

	for editor.running {
		editor.mode.Display()

		if editor.mode.IsCommand() {
			editor.commandInput.Draw()
		}

		editor.DisplayError()

		termbox.Flush()
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	}
}
