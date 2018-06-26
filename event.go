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
		if err := v.InsertEvent(ev); err != nil {
			return err
		}
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
	case 'i':
		v.mode = 1
	}
	return nil
}

func (v *View) InsertEvent(ev *tcell.EventKey) error {
	switch ev.Key() {
	case tcell.KeyEsc:
		// chage to Normal mode
		v.mode = 0
	case tcell.KeyEnter:
		v.buf.Insert([]rune{'\n'})

		// update cursor
		v.buf.Cursor.x = 0
		v.buf.CursorMove(MoveC('j'))
	case tcell.KeyBackspace:
		// [TODO] remove the data from buffer
	default:
		v.buf.Insert([]rune{ev.Rune()})

		// update cursor
		v.buf.CursorMove(MoveC('l'))
	}

	v.buf.Render(0)
	return nil
}
