package main

import (
	"github.com/nsf/termbox-go"
)

func main() {

  editor := CreateEditor()

	defer termbox.Close()

  termbox.Flush()

	for editor.running {
    editor.mode.Display()

    if editor.mode.IsCommand() {
      editor.commandInput.Draw()
    }

    termbox.Flush()
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	}
}
