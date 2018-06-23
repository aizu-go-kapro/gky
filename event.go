package main

import "github.com/gdamore/tcell"

func (v *View) EventHandle(key tcell.Key) error {
	switch v.mode {
	case Normal:
		if err := v.NormalEvent(key); err != nil {
			return err
		}
	case Insert:
	case Visual:
	}
	return nil
}

func (v *View) NormalEvent(key tcell.Key) error {
	switch key {
	case tcell.KeyBackspace, tcell.KeyLeft, 'h', tcell.KeyRight, 'l', tcell.KeyUp, 'k', tcell.KeyDown, 'j', tcell.KeyEnter:
		v.buf.CursorMove(key)
	default:
	}
	return nil
}

func (buf *Buffer) CursorMove(key tcell.Key) {
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
