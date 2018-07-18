package main

import (
	"fmt"

	"github.com/gdamore/tcell"
)

func (buf *Buffer) SetHighlightBegine() {
	if buf.HighlightBegine.x == -1 && buf.HighlightBegine.y == -1 {
		buf.HighlightBegine.x = buf.Cursor.x
		buf.HighlightBegine.y = buf.Cursor.y
	}
}

func (buf *Buffer) SetHighlightEnd() {
	buf.HighlightEnd.x = buf.Cursor.x
	buf.HighlightEnd.y = buf.Cursor.y
}

func (buf *Buffer) ClearHighlight() {
	buf.HighlightBegine = &Location{x: -1, y: -1}
	buf.HighlightEnd = &Location{x: -1, y: -1}
}

func isAhead(buf *Buffer) bool {
	if buf.Cursor.y > buf.HighlightBegine.y || (buf.Cursor.y == buf.HighlightBegine.y && buf.Cursor.x > buf.HighlightBegine.x) {
		return true
	} else {
		return false
	}
}

func isEqualLine(c1, c2 *Location) bool {
	if c1.y == c2.y {
		return true
	} else {
		return false
	}
}

func highlightLine(start int, end int, line int, buf *Buffer, isReverse bool) {
	spaces := len(fmt.Sprintf("%d", len(buf.data))) + 1
	if isReverse {
		for i := start; i < end; i++ {
			screen.SetContent(i, line, buf.data[line][i-spaces], nil, tcell.StyleDefault.Reverse(true))
		}
	} else {
		for i := start; i < end; i++ {
			screen.SetContent(i, line, buf.data[line][i-spaces], nil, tcell.StyleDefault)
		}
	}
}

func (buf *Buffer) SetStyleHighlight() {
	spaces := len(fmt.Sprintf("%d", len(buf.data))) + 1
	if isAhead(buf) {
		if isEqualLine(buf.HighlightBegine, buf.Cursor) {
			highlightLine(buf.HighlightBegine.x, buf.Cursor.x+1, buf.Cursor.y, buf, true)
		} else {
			if buf.Cursor.y == buf.HighlightBegine.y+1 {
				highlightLine(buf.HighlightBegine.x, len(buf.data[buf.HighlightBegine.y])+spaces, buf.HighlightBegine.y, buf, true)
			} else {
				highlightLine(spaces, len(buf.data[buf.Cursor.y-1])+spaces, buf.Cursor.y-1, buf, true)
			}
			highlightLine(spaces, buf.Cursor.x, buf.Cursor.y, buf, true)
		}
	} else {
		if isEqualLine(buf.HighlightBegine, buf.Cursor) {
			highlightLine(buf.Cursor.x, buf.HighlightBegine.x+1, buf.Cursor.y, buf, true)
		} else {
			if buf.Cursor.y == buf.HighlightBegine.y-1 {
				highlightLine(spaces, buf.HighlightBegine.x, buf.HighlightBegine.y, buf, true)
			} else {
				highlightLine(spaces, len(buf.data[buf.Cursor.y+1])+spaces, buf.Cursor.y+1, buf, true)
			}
			highlightLine(buf.Cursor.x, len(buf.data[buf.Cursor.y])+spaces, buf.Cursor.y, buf, true)
		}
	}
}

func (buf *Buffer) SetStyleDefault() {
	spaces := len(fmt.Sprintf("%d", len(buf.data))) + 1
	if isAhead(buf) {
		if buf.Cursor.y < buf.HighlightEnd.y {
			highlightLine(spaces, len(buf.data[buf.HighlightEnd.y])+spaces, buf.HighlightEnd.y, buf, false)
			highlightLine(buf.Cursor.x+1, len(buf.data[buf.Cursor.y])+spaces, buf.Cursor.y, buf, false)
		} else if buf.Cursor.x < buf.HighlightEnd.x {
			screen.SetContent(buf.HighlightEnd.x, buf.HighlightEnd.y, buf.data[buf.HighlightEnd.y][buf.HighlightEnd.x-spaces], nil, tcell.StyleDefault)
		}
	} else {
		if buf.Cursor.y > buf.HighlightEnd.y {
			highlightLine(spaces, len(buf.data[buf.HighlightEnd.y])+spaces, buf.HighlightEnd.y, buf, false)
			highlightLine(spaces, buf.Cursor.x-1, buf.Cursor.y, buf, false)
		} else if buf.Cursor.x > buf.HighlightEnd.x {
			screen.SetContent(buf.HighlightEnd.x, buf.HighlightEnd.y, buf.data[buf.HighlightEnd.y][buf.HighlightEnd.x-spaces], nil, tcell.StyleDefault)
		}
	}
}

func (buf *Buffer) ClearStyleHighlight() {
	spaces := len(fmt.Sprintf("%d", len(buf.data))) + 1
	if isAhead(buf) {
		for i := buf.HighlightBegine.y; i <= buf.HighlightEnd.y; i++ {
			highlightLine(spaces, len(buf.data[i])+spaces, i, buf, false)
		}
	} else {
		for i := buf.HighlightEnd.y; i <= buf.HighlightBegine.y; i++ {
			highlightLine(spaces, len(buf.data[i])+spaces, i, buf, false)
		}
	}
}

func (v *View) Yank() {
	v.ExitVisualMode()
}

func (v *View) ExitVisualMode() {
	v.buf.ClearStyleHighlight()
	v.buf.ClearHighlight()
	v.mode = 0
}
