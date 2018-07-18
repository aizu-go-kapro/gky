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

//TODO: もっときれいに書け
func (buf *Buffer) SetStyleHighlight() {
	//行番号と空白分のスペース
	spaces := len(fmt.Sprintf("%d", len(buf.data))) + 1
	if isAhead(buf) {
		if buf.Cursor.y > buf.HighlightEnd.y {
			if buf.Cursor.y == buf.HighlightBegine.y+1 {
				for i := buf.HighlightEnd.x; i < len(buf.data[buf.HighlightEnd.y])+spaces; i++ {
					screen.SetContent(i, buf.HighlightEnd.y, buf.data[buf.HighlightEnd.y][i-spaces], nil, tcell.StyleDefault.Reverse(true))
				}
			} else {
				for i := spaces; i < len(buf.data[buf.HighlightEnd.y])+spaces; i++ {
					screen.SetContent(i, buf.HighlightEnd.y, buf.data[buf.HighlightEnd.y][i-spaces], nil, tcell.StyleDefault.Reverse(true))
				}
			}
			for j := spaces; j <= buf.Cursor.x; j++ {
				screen.SetContent(j, buf.Cursor.y, buf.data[buf.Cursor.y][j-spaces], nil, tcell.StyleDefault.Reverse(true))
			}
		} else {
			screen.SetContent(buf.Cursor.x, buf.Cursor.y, buf.data[buf.Cursor.y][buf.Cursor.x-spaces], nil, tcell.StyleDefault.Reverse(true))
		}
	} else {
		if buf.Cursor.y < buf.HighlightEnd.y {
			if buf.Cursor.y == buf.HighlightBegine.y-1 {
				for i := buf.HighlightEnd.x; i >= spaces; i-- {
					screen.SetContent(i, buf.HighlightEnd.y, buf.data[buf.HighlightEnd.y][i-spaces], nil, tcell.StyleDefault.Reverse(true))
				}
			} else {
				for i := spaces; i < len(buf.data[buf.HighlightEnd.y])+spaces; i++ {
					screen.SetContent(i, buf.HighlightEnd.y, buf.data[buf.HighlightEnd.y][i-spaces], nil, tcell.StyleDefault.Reverse(true))
				}
			}
			for j := len(buf.data[buf.Cursor.y]) + spaces - 1; j >= buf.Cursor.x; j-- {
				screen.SetContent(j, buf.Cursor.y, buf.data[buf.Cursor.y][j-spaces], nil, tcell.StyleDefault.Reverse(true))
			}
		} else {
			screen.SetContent(buf.Cursor.x, buf.Cursor.y, buf.data[buf.Cursor.y][buf.Cursor.x-spaces], nil, tcell.StyleDefault.Reverse(true))
		}
	}
}

func (buf *Buffer) SetStyleDefault() {
	spaces := len(fmt.Sprintf("%d", len(buf.data))) + 1
	if isAhead(buf) {
		if buf.Cursor.y < buf.HighlightEnd.y {
			for i := spaces; i < len(buf.data[buf.HighlightEnd.y])+spaces; i++ {
				screen.SetContent(i, buf.HighlightEnd.y, buf.data[buf.HighlightEnd.y][i-spaces], nil, tcell.StyleDefault)
			}
			for j := len(buf.data[buf.Cursor.y]) + spaces - 1; j > buf.Cursor.x; j-- {
				screen.SetContent(j, buf.Cursor.y, buf.data[buf.Cursor.y][j-spaces], nil, tcell.StyleDefault)
			}
		} else if buf.Cursor.x < buf.HighlightEnd.x {
			screen.SetContent(buf.HighlightEnd.x, buf.HighlightEnd.y, buf.data[buf.HighlightEnd.y][buf.HighlightEnd.x-spaces], nil, tcell.StyleDefault)
		}
	} else {
		if buf.Cursor.y > buf.HighlightEnd.y {
			for i := spaces; i < len(buf.data[buf.HighlightEnd.y])+spaces; i++ {
				screen.SetContent(i, buf.HighlightEnd.y, buf.data[buf.HighlightEnd.y][i-spaces], nil, tcell.StyleDefault)
			}
			for j := spaces; j <= buf.Cursor.x-1; j++ {
				screen.SetContent(j, buf.Cursor.y, buf.data[buf.Cursor.y][j-spaces], nil, tcell.StyleDefault)
			}
		} else if buf.Cursor.x > buf.HighlightEnd.x {
			screen.SetContent(buf.HighlightEnd.x, buf.HighlightEnd.y, buf.data[buf.HighlightEnd.y][buf.HighlightEnd.x-spaces], nil, tcell.StyleDefault)
		}
	}
}

func (buf *Buffer) ClearStyleHighlight() {
	spaces := len(fmt.Sprintf("%d", len(buf.data))) + 1
	for i := buf.HighlightBegine.y; i <= buf.HighlightEnd.y; i++ {
		for j := 0 + spaces; j < len(buf.data[i])+spaces; j++ {
			screen.SetContent(j, i, buf.data[i][j-spaces], nil, tcell.StyleDefault)
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
