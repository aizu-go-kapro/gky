package main

import (
	"fmt"

	"github.com/gdamore/tcell"
)

func (buf *Buffer) CursorMove(key MoveC) {
	switch key {
	case tcell.KeyBackspace, tcell.KeyBackspace2, tcell.KeyLeft, 'h':
		buf.Cursor.x--
	case tcell.KeyRight, 'l':
		buf.Cursor.x++
	case tcell.KeyUp, 'k':
		buf.Cursor.y--
		buf.render_y--
	case tcell.KeyDown, tcell.KeyEnter, 'j':
		buf.Cursor.y++
		buf.render_y++
	}

	buf.scrollCursor()
	buf.setCursor()
	screen.Show()
}

func (buf *Buffer) CursorMoveVisual(key MoveC) {
	switch key {
	case tcell.KeyBackspace, tcell.KeyBackspace2, tcell.KeyLeft, 'h':
		buf.Cursor.x--
	case tcell.KeyRight, 'l':
		buf.Cursor.x++
	case tcell.KeyUp, 'k':
		buf.Cursor.y--
		buf.render_y--
	case tcell.KeyDown, tcell.KeyEnter, 'j':
		buf.Cursor.y++
		buf.render_y++
	}

	buf.scrollCursor()
	buf.SetStyleHighlight()
	buf.SetStyleDefault()
	buf.SetHighlightEnd()
	buf.setCursor()
	screen.Show()
}

func (buf *Buffer) scrollCursor() {
	line := len(fmt.Sprintf("%d", buf.getLine()))

	y := buf.getLine() - 1
	if y > screenHeight {
		y = screenHeight - 1
	}

	if buf.Cursor.y < 0 {
		buf.Cursor.y = 0
		buf.CursorRender()
	} else if buf.Cursor.y >= y {
		buf.Cursor.y = y
		buf.CursorRender()
	}

	if buf.Cursor.x < line+1 {
		buf.Cursor.x = line + 1
	} else if buf.Cursor.x >= len(buf.data[buf.render_y])+line {
		buf.Cursor.x = len(buf.data[buf.render_y]) + line
	}
}

func (buf *Buffer) CursorRender() {
	line := buf.getLine() - 1
	if buf.render_y <= line && buf.render_y >= 0 {
		buf.Render(buf.render_y - buf.Cursor.y)
	} else if buf.render_y > line {
		buf.render_y = line
	} else {
		buf.render_y = 0
	}
}

func (buf *Buffer) setCursor() {
	screen.ShowCursor(buf.Cursor.x, buf.Cursor.y)
}
