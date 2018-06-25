package main

import "github.com/gdamore/tcell"

func (buf *Buffer) CursorMove(key MoveC) {
	switch key {
	case tcell.KeyBackspace, tcell.KeyLeft, 'h':
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

func (buf *Buffer) scrollCursor() {
	//TODO: x < 0 -> y--, x > width -> y++
	if buf.Cursor.x < 0 {
		buf.Cursor.x = 0
	} else if buf.Cursor.x >= screenWidth {
		buf.Cursor.x = screenWidth - 1
	}

	if buf.Cursor.y < 0 {
		buf.Cursor.y = 0
		buf.CursorRender()
	} else if buf.Cursor.y >= screenHeight {
		buf.Cursor.y = screenHeight - 1
		buf.CursorRender()
	}
}

func (buf *Buffer) CursorRender() {
	line := buf.getLine()
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
