package main

import (
	"log"
	"time"

	"github.com/marcos-venicius/daily-term/taskmanagement"
	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()

	if err != nil {
		log.Fatal(err)
	}

	repository, err := taskmanagement.CreateRepository()

	if err != nil {
		log.Fatal(err)
	}

	editor := CreateEditor(repository)

	defer termbox.Close()
	defer repository.CloseRepository()
	defer close(editor.termbox_event)

	editor.InitParser()

	termbox.Flush()

	for editor.running {
		update := time.Now()

		editor.mode.Display()

		editor.DisplayTasks()

		if editor.mode.IsCommand() {
			editor.commandInput.Draw()
		}

		editor.DisplayError()
		editor.DisplayInfo()

		termbox.Flush()
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		time.Sleep(time.Duration((update.Sub(time.Now()).Seconds()*1000.0)+1000.0/editor.fps) * time.Millisecond)
	}
}
