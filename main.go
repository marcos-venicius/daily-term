package main

import (
	"log"

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

	editor.InitParser()

	termbox.Flush()
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}

	for editor.running {
		editor.mode.Display()

		editor.DisplayTasks()

		if editor.mode.IsCommand() {
			editor.commandInput.Draw()
		}

		editor.DisplayError()
		editor.DisplayInfo()

		termbox.Flush()
		if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
			panic(err)
		}

		event := termbox.PollEvent()

		editor.ListenEvents(event)
	}
}
