package main

import (
	"fmt"

	"github.com/marcos-venicius/daily-term/argument-parser"
	"github.com/marcos-venicius/daily-term/taskmanagement"
	"github.com/nsf/termbox-go"
)

// These are all the modes available in the application
const (
	NormalMode  EditorMode = iota // when the user wants to be able to select another editor mode
	CommandMode EditorMode = iota // when the user wants execute some command like quit (q)
	DeleteMode  EditorMode = iota // when the user wants delete a task
)

type EditorMode int

type Editor struct {
	mode           EditorMode // current editor mode (default is Normal Mode)
	termbox_event  chan termbox.Event
	running        bool
	commandInput   *Input
	argumentParser *argumentparser.ArgumentParser
	errorMessage   string
	infoMessage    string
	width          int
	height         int
	board          *taskmanagement.Board
	fps            float64
}

func CreateEditor() *Editor {
	windowWidth, windowHeight := termbox.Size()

	var termbox_event chan termbox.Event = make(chan termbox.Event, 20)

	termbox.SetInputMode(termbox.InputEsc)

	argumentParser := argumentparser.CreateArgumentParser()
	board := taskmanagement.CreateBoard()

	// TODO: remove this tasks
	board.AddTask("Lorem ipsum dolor sit ammet consectur")
	board.AddTask("This is a brand new task")
	board.AddTask("Another nice task")
	board.AddTask("Other task")

	editor := &Editor{
		mode:           NormalMode,
		termbox_event:  termbox_event,
		running:        true,
		commandInput:   CreateInput(windowWidth, 1, 0, windowHeight-1),
		argumentParser: argumentParser,
		board:          board,
		width:          windowWidth,
		height:         windowHeight,
		fps:            50,
	}

	go func() {
		for editor.running {
			editor.termbox_event <- termbox.PollEvent()
		}
	}()

	editor.listenEvents()

	return editor
}

func (editor *Editor) Stop() {
	editor.running = false
	termbox.Interrupt()
}

func (editor *Editor) SetDeleteMode() {
	if editor.board.HasTasks() {
		editor.mode = DeleteMode
	}
}

func (editor *Editor) SetNormalMode() {
	editor.mode = NormalMode
}

func (editor *Editor) SetCommandMode() {
	editor.mode = CommandMode
}

func (mode *EditorMode) IsDelete() bool {
	return *mode == DeleteMode
}

func (mode *EditorMode) IsNormal() bool {
	return *mode == NormalMode
}

func (mode *EditorMode) IsCommand() bool {
	return *mode == CommandMode
}

func (mode *EditorMode) Display() {
	switch *mode {
	case NormalMode:
		tbprint(0, 0, termbox.ColorWhite, termbox.ColorDefault, "NORMAL")
		break
	case CommandMode:
		tbprint(0, 0, termbox.ColorWhite, termbox.ColorDefault, "COMMAND")
		break
	case DeleteMode:
		tbprint(0, 0, termbox.ColorRed, termbox.ColorDefault, "DELETE")
		break
	default:
		tbprint(0, 0, termbox.ColorWhite, termbox.ColorDefault, "UNKNOWN")
		break
	}
}

func (editor *Editor) DisplayTasks() {
	const startingRow = 2

	for row, task := range editor.board.Tasks() {
		color := termbox.ColorWhite

		switch task.State {
		case taskmanagement.InProgress:
			color = termbox.ColorYellow
			break
		case taskmanagement.Completed:
			color = termbox.ColorGreen
			break
		}

		selectedSymbol := task.Symbol(*editor.board.SelectedTaskId())

		if task.Id == *editor.board.SelectedTaskId() {
			if editor.mode.IsDelete() {
				selectedSymbol = '-'
				color = termbox.ColorLightRed
			} else {
				color = termbox.ColorLightCyan
			}
		}

		text := fmt.Sprintf("%c [%04d] %v", selectedSymbol, task.Id, task.Name)

		tbprint(0, startingRow+row, color, termbox.ColorDefault, text)
	}
}

func (editor *Editor) DisplayError() {
	if editor.errorMessage != "" && editor.mode.IsNormal() {
		errorMessage := fmt.Sprintf("ERROR: %v", editor.errorMessage)

		tbprint(0, editor.height-1, termbox.ColorRed, termbox.ColorDefault, errorMessage)
	}
}

func (editor *Editor) DisplayInfo() {
	if editor.infoMessage != "" && editor.mode.IsNormal() {
		tbprint(0, editor.height-1, termbox.ColorWhite, termbox.ColorDefault, editor.infoMessage)
	}
}

func (editor *Editor) listenNormalModeEvents(event termbox.Event) {
	if !editor.running {
		return
	}

	if event.Type != termbox.EventKey {
		return
	}

	if len(editor.errorMessage) > 0 {
		editor.errorMessage = ""
	}

	if len(editor.infoMessage) > 0 {
		editor.infoMessage = ""
	}

	switch event.Ch {
	case ':':
		editor.SetCommandMode()
		break
	case 'd':
		editor.SetDeleteMode()
		break
	case 'q':
		editor.Stop()
		break
	case 'j':
		editor.board.SelectNextTask()
		break
	case 'k':
		editor.board.SelectPreviousTask()
		break
	case 't':
		editor.ChangeCurrentTaskStateFor(taskmanagement.Todo)
		break
	case 'i':
		editor.ChangeCurrentTaskStateFor(taskmanagement.InProgress)
		break
	case 'c':
		editor.ChangeCurrentTaskStateFor(taskmanagement.Completed)
		break
	default:
		break
	}
}

func (editor *Editor) listenCommandModeEvents(event termbox.Event) {
	if !editor.running {
		return
	}

	if event.Type != termbox.EventKey {
		return
	}

	if event.Key == termbox.KeyEsc {
		editor.commandInput.Reset()
		editor.SetNormalMode()
	}
}

func (editor *Editor) listenDeleteModeEvents(event termbox.Event) {
	if !editor.running {
		return
	}

	if event.Type != termbox.EventKey {
		return
	}

	if event.Key == termbox.KeyEsc {
		editor.SetNormalMode()
	}

	switch event.Ch {
	case 'd':
		editor.exec("delete task")
		break
	default:
		break
	}
}

func (editor *Editor) listenEvents() {
	go func() {
		for editor.running {
			event := <-editor.termbox_event

			if editor.mode.IsNormal() {
				editor.listenNormalModeEvents(event)
			} else if editor.mode.IsCommand() {
				editor.listenCommandModeEvents(event)
			} else if editor.mode.IsDelete() {
				editor.listenDeleteModeEvents(event)
			}

			editor.commandInput.handleEvents(editor, event)
		}
	}()
}

func (editor *Editor) SetErrorMessage(message string) {
	editor.infoMessage = ""
	editor.errorMessage = message
}

func (editor *Editor) SetInfoMessage(message string) {
	editor.errorMessage = ""
	editor.infoMessage = message
}

func (editor *Editor) exec(command string) {
	editor.SetNormalMode()

	cmd, err := editor.argumentParser.ParseFromString(command)

	if err != nil {
		editor.SetErrorMessage(err.Error())
		return
	}

	switch cmd.Name {
	case "quit", "q":
		editor.Quit()
		break
	case "new task":
		editor.addTask(cmd.Arguments)
		break
	case "delete task":
		editor.deleteTask(cmd.Arguments)
		break
	default:
		editor.SetErrorMessage(fmt.Sprintf(`Unhandled command "%v"`, cmd.Name))
		break
	}
}
