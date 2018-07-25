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
		buf.renderY--
	case tcell.KeyDown, tcell.KeyEnter, 'j':
		buf.Cursor.y++
		buf.renderY++
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
		buf.renderY--
	case tcell.KeyDown, tcell.KeyEnter, 'j':
		buf.Cursor.y++
		buf.renderY++
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
	} else if buf.Cursor.x >= len(buf.data[buf.renderY])+line {
		buf.Cursor.x = len(buf.data[buf.renderY]) + line
	}
}

func (buf *Buffer) CursorRender() {
	line := buf.getLine() - 1
	if buf.renderY <= line && buf.renderY >= 0 {
		buf.Render(buf.renderY - buf.Cursor.y)
	} else if buf.renderY > line {
		buf.renderY = line
	} else {
		buf.renderY = 0
	}
}

func (buf *Buffer) setCursor() {
	screen.ShowCursor(buf.Cursor.x, buf.Cursor.y)
}
