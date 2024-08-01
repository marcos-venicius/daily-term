package main

import (
	"github.com/nsf/termbox-go"
	"log"
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
