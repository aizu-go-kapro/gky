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
