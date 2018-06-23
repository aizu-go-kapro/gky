package main

import "github.com/gdamore/tcell"

type MoveC interface{}

func (v *View) EventHandle(ev *tcell.EventKey) error {
	switch v.mode {
	case Normal:
		if err := v.NormalEvent(ev); err != nil {
			return err
		}
	case Insert:
	case Visual:
	}
	return nil
}

func (v *View) NormalEvent(ev *tcell.EventKey) error {
	switch ev.Key() {
	case tcell.KeyBackspace, tcell.KeyLeft, tcell.KeyRight, tcell.KeyUp, tcell.KeyDown, tcell.KeyEnter:
		v.buf.CursorMove(MoveC(ev.Key()))
	default:
	}

	switch ev.Rune() {
	case 'j', 'h', 'l', 'k':
		v.buf.CursorMove(MoveC(ev.Rune()))
	}
	return nil
}

func (buf *Buffer) CursorMove(key MoveC) {
	switch key {
	case tcell.KeyBackspace, tcell.KeyLeft, 'h':
		buf.Cursor.x--
	case tcell.KeyRight, 'l':
		buf.Cursor.x++
	case tcell.KeyUp, 'k':
		buf.Cursor.y--
	case tcell.KeyDown, tcell.KeyEnter, 'j':
		buf.Cursor.y++
	}

	buf.scrollCursor()
	buf.setCursor()
	screen.Show()
}

func (buf *Buffer) scrollCursor() {
	//TODO: x < 0 -> y--, x > width -> y++
	if buf.Cursor.x < 0 {
		buf.Cursor.x = 0
	} else if buf.Cursor.x >= screenWidth {
		buf.Cursor.x = screenWidth - 1
	}

	if buf.Cursor.y < 0 {
		buf.Cursor.y = 0
	} else if buf.Cursor.y >= screenHeight {
		buf.Cursor.y = screenHeight - 1
	}
}
