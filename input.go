package main

import (
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const (
	PREFERED_HORIZONTAL_THRESHOLD = 5
	TABSTOP_LENGTH                = 8
)

var arrowLeft = '←'
var arrowRight = '→'

type Input struct {
	text           []byte
	line_voffset   int
	cursor_boffset int // cursor offset in bytes
	cursor_voffset int // visual cursor offset in termbox cells
	cursor_coffset int // cursor offset in unicode code points
	width          int
	height         int
	x              int
	y              int
}

func CreateInput(width, height, x, y int) *Input {
	input := &Input{
		width:  width,
		height: height,
		x:      x,
		y:      y,
	}

	return input
}

func (input *Input) AdjustVOffset(width int) {
	ht := PREFERED_HORIZONTAL_THRESHOLD

	max_h_threshold := (width - 1) / 2

	if ht > max_h_threshold {
		ht = max_h_threshold
	}

	threshold := width - 1

	if input.line_voffset != 0 {
		threshold = width - ht
	}

	if input.cursor_voffset-input.line_voffset >= threshold {
		input.line_voffset = input.cursor_voffset + (ht - width + 1)
	}

	if input.line_voffset != 0 && input.cursor_voffset-input.line_voffset < ht {
		input.line_voffset = input.cursor_voffset - ht

		if input.line_voffset < 0 {
			input.line_voffset = 0
		}
	}
}

func (input *Input) MoveCursorTo(boffset int) {
	input.cursor_boffset = boffset
	input.cursor_voffset, input.cursor_coffset = voffset_coffset(input.text, boffset)
}

func (input *Input) RuneUnderCursor() (rune, int) {
	return utf8.DecodeRune(input.text[input.cursor_boffset:])
}

func (input *Input) RuneBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(input.text[:input.cursor_boffset])
}

func (input *Input) MoveCursorOneRuneBackward() {
	if input.cursor_boffset == 0 {
		return
	}

	_, size := input.RuneBeforeCursor()

	input.MoveCursorTo(input.cursor_boffset - size)
}

func (input *Input) MoveCursorOneRuneForward() {
	if input.cursor_boffset == len(input.text) {
		return
	}

	_, size := input.RuneUnderCursor()

	input.MoveCursorTo(input.cursor_boffset + size)
}

func (input *Input) MoveCursorToBeginningOfTheLine() {
	input.MoveCursorTo(0)
}

func (input *Input) MoveCursorToEndOfTheLine() {
	input.MoveCursorTo(len(input.text))
}

func (input *Input) DeleteRuneBackward() {
	if input.cursor_boffset == 0 {
		return
	}

	input.MoveCursorOneRuneBackward()

	_, size := input.RuneUnderCursor()

	input.text = byte_slice_remove(input.text, input.cursor_boffset, input.cursor_boffset+size)
}

func (input *Input) DeleteRuneForward() {
	if input.cursor_boffset == len(input.text) {
		return
	}

	_, size := input.RuneUnderCursor()

	input.text = byte_slice_remove(input.text, input.cursor_boffset, input.cursor_boffset+size)
}

func (input *Input) DeleteTheRestOfTheLine() {
	input.text = input.text[:input.cursor_boffset]
}

func (input *Input) InsertRune(r rune) {
	var buf [utf8.UTFMax]byte

	n := utf8.EncodeRune(buf[:], r)

	input.text = byte_slice_insert(input.text, input.cursor_boffset, buf[:n])
	input.MoveCursorOneRuneForward()
}

// Please, keep in mind that cursor depends on the value of line_voffset, which
// is being set on Draw() call, so.. call this method after Draw() one.
func (input *Input) CursorX() int {
	return input.cursor_voffset - input.line_voffset
}

func (input *Input) Draw() {
	fill(input.x, input.y-1, input.width, 1, termbox.Cell{Ch: '─', Fg: termbox.ColorRed})

	input.drawInput(input.x, input.y, input.width, input.height)

	termbox.SetCursor(input.x+input.CursorX(), input.y)
}

func (input *Input) Reset() {
	input.text = []byte{}
	input.line_voffset = 0
	input.cursor_boffset = 0
	input.cursor_voffset = 0
	input.cursor_coffset = 0

	termbox.HideCursor()
}

func (input *Input) startListeningEvents(editor *Editor) {
	go func() {
		for editor.mode.IsCommand() && editor.running {
			event := termbox.PollEvent()

			switch event.Type {
			case termbox.EventKey:
				switch event.Key {
				case termbox.KeyArrowLeft:
				case termbox.KeyCtrlH:
					input.MoveCursorOneRuneBackward()
				case termbox.KeyArrowRight:
				case termbox.KeyCtrlL:
					input.MoveCursorOneRuneForward()
				case termbox.KeyBackspace2:
					input.DeleteRuneBackward()
				case termbox.KeyDelete, termbox.KeyCtrlD:
					input.DeleteRuneForward()
				case termbox.KeyTab:
					input.InsertRune('\t')
				case termbox.KeySpace:
					input.InsertRune(' ')
				case termbox.KeyCtrlK:
					input.DeleteTheRestOfTheLine()
				case termbox.KeyHome, termbox.KeyCtrlA:
					input.MoveCursorToBeginningOfTheLine()
				case termbox.KeyEnd, termbox.KeyCtrlE:
					input.MoveCursorToEndOfTheLine()
				case termbox.KeyEnter:
					command := string(input.text)
					input.Reset()
					editor.exec(command)
					return
				default:
					if event.Ch != 0 {
						input.InsertRune(event.Ch)
					}
				}
			}
		}
	}()
}

func (input *Input) drawInput(x, y, w, h int) {
	input.AdjustVOffset(w)

	const coldef = termbox.ColorDefault
	const colred = termbox.ColorRed

	fill(x, y, w, h, termbox.Cell{Ch: ' '})

	t := input.text
	lx := 0
	tabstop := 0

	for {
		rx := lx - input.line_voffset

		if len(t) == 0 {
			break
		}

		if lx == tabstop {
			tabstop += TABSTOP_LENGTH
		}

		if rx >= w {
			termbox.SetCell(x+w-1, y, arrowRight, colred, coldef)
			break
		}

		r, size := utf8.DecodeRune(t)

		if r == '\t' {
			for ; lx < tabstop; lx++ {
				rx = lx - input.line_voffset

				if rx >= w {
					goto next
				}

				if rx >= 0 {
					termbox.SetCell(x+rx, y, ' ', coldef, coldef)
				}
			}
		} else {
			if rx >= 0 {
				termbox.SetCell(x+rx, y, r, coldef, coldef)
			}
			lx += runewidth.RuneWidth(r)
		}
	next:
		t = t[size:]
	}

	if input.line_voffset != 0 {
		termbox.SetCell(x, y, arrowLeft, colred, coldef)
	}
}
