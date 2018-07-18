package main

import (
	"fmt"

	"github.com/gdamore/tcell"
)

type MoveC interface{}

func (v *View) EventHandle(ev *tcell.EventKey) error {
	switch v.mode {
	case Normal:
		if err := v.NormalEvent(ev); err != nil {
			return err
		}
	case Insert:
		if err := v.InsertEvent(ev); err != nil {
			return err
		}
	case Visual:
		spaces := len(fmt.Sprintf("%d", len(v.buf.data))) + 1
		screen.SetContent(v.buf.Cursor.x, v.buf.Cursor.y, v.buf.data[v.buf.Cursor.y][v.buf.Cursor.x-spaces], nil, tcell.StyleDefault.Reverse(true))
		if err := v.VisualEvnet(ev); err != nil {
			return err
		}
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
	case 'i':
		v.mode = Insert
	case 'v':
		v.mode = Visual
	}
	return nil
}

func (v *View) VisualEvnet(ev *tcell.EventKey) error {
	v.buf.SetHighlightBegine()
	v.buf.SetHighlightEnd()

	switch ev.Key() {
	case tcell.KeyEsc:
		v.ExitVisualMode()
	case tcell.KeyBackspace, tcell.KeyLeft, tcell.KeyRight, tcell.KeyUp, tcell.KeyDown, tcell.KeyEnter:
		v.buf.CursorMoveVisual(MoveC(ev.Key()))
	}
	switch ev.Rune() {
	case 'j', 'h', 'l', 'k':
		//選択部分の色を反転させる処理
		v.buf.CursorMoveVisual(MoveC(ev.Rune()))
	case 'y':
		v.Yank()
	}
	return nil
}

func (v *View) InsertEvent(ev *tcell.EventKey) error {
	switch ev.Key() {
	case tcell.KeyLeft, tcell.KeyRight, tcell.KeyUp, tcell.KeyDown:
		v.buf.CursorMove(MoveC(ev.Key()))
	case tcell.KeyEsc:
		// chage to Normal mode
		v.mode = Normal
	case tcell.KeyEnter:
		v.buf.Insert([]rune{'\n'})

		// update cursor
		v.buf.Cursor.x = 0
		v.buf.CursorMove(MoveC('j'))
	case tcell.KeyBackspace2:
		lastCursor := v.buf.Remove(1)
		if lastCursor != 0 {
			// update cursor
			v.buf.Cursor.x = lastCursor
			v.buf.CursorMove(MoveC('k'))
		} else {
			// update cursor
			v.buf.CursorMove(MoveC('h'))
		}
	default:
		v.buf.Insert([]rune{ev.Rune()})

		// update cursor
		v.buf.CursorMove(MoveC('l'))
	}

	v.buf.Render(0)
	return nil
}
