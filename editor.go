package main

import (
	"log"

	"github.com/nsf/termbox-go"
)

const (
	EDITOR_MODE_NORMAL  EditorMode = iota // when the user wants to be able to select another editor mode
	EDITOR_MODE_COMMAND EditorMode = iota // when the user wants execute some command like quit (q)
)

type EditorMode int

type Editor struct {
	mode          EditorMode // current editor mode (default is Normal Mode)
	termbox_event chan termbox.Event
	running       bool
	commandInput  *Input
}

func CreateEditor() *Editor {
	err := termbox.Init()

	if err != nil {
		log.Fatal(err)
	}

  windowWidth, windowHeight := termbox.Size()

	var termbox_event chan termbox.Event = make(chan termbox.Event, 20)

	termbox.SetInputMode(termbox.InputEsc)

	editor := &Editor{
		mode:          EDITOR_MODE_NORMAL,
		termbox_event: termbox_event,
		running:       true,
    commandInput: CreateInput(windowWidth, 1, 0, windowHeight - 1),
	}

	go func() {
		for editor.running {
			editor.termbox_event <- termbox.PollEvent()
		}
	}()

	editor.listenEvents()

	return editor
}

func (editor *Editor) SetNormalMode() {
	editor.mode = EDITOR_MODE_NORMAL
}

func (editor *Editor) SetCommandMode() {
	editor.mode = EDITOR_MODE_COMMAND
}

func (mode *EditorMode) IsNormal() bool {
	return *mode == EDITOR_MODE_NORMAL
}

func (mode *EditorMode) IsCommand() bool {
	return *mode == EDITOR_MODE_COMMAND
}

func (mode *EditorMode) Display() {
	switch *mode {
	case EDITOR_MODE_NORMAL:
		tbprint(0, 0, termbox.ColorRed, termbox.ColorDefault, "NORMAL")
		break
	case EDITOR_MODE_COMMAND:
		tbprint(0, 0, termbox.ColorRed, termbox.ColorDefault, "COMMAND")
		break
	default:
		tbprint(0, 0, termbox.ColorRed, termbox.ColorDefault, "UNKNOWN")
		break
	}
}

func (editor *Editor) listenNormalModeEvents() {
	event := <-editor.termbox_event

	if event.Type != termbox.EventKey {
		return
	}

	switch event.Ch {
	case rune(':'):
		editor.commandInput.startListeningEvents(editor)
		editor.SetCommandMode()
		return
	case rune('q'):
		editor.running = false
		return
	default:
		return
	}
}

func (editor *Editor) listenCommandModeEvents() {
	event := <-editor.termbox_event

	if event.Type != termbox.EventKey {
		return
	}

	if event.Key == termbox.KeyEsc {
    editor.commandInput.Reset()
		editor.SetNormalMode()
	}
}

func (editor *Editor) listenEvents() {
	go func() {
		for editor.running {
			if editor.mode.IsNormal() {
				editor.listenNormalModeEvents()
			} else if editor.mode.IsCommand() {
				editor.listenCommandModeEvents()
			}
		}
	}()
}
