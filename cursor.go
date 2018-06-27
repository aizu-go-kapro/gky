package main

import "github.com/gdamore/tcell"

// by tenntenn もうちょい抽象化したほうが良さそう。tcellに依存し過ぎな気がする。
// キーコードを使うのではなく、それを対応付けたものを定義したほうがいいかな
// 例えば、ユーザが好きにキーバインドを変える機能を追加したときに変更点を減らすにはどうすればよいか？
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

func (buf *Buffer) scrollCursor() {
	if buf.Cursor.y < 0 {
		buf.Cursor.y = 0
		buf.CursorRender()
	} else if buf.Cursor.y >= screenHeight {
		buf.Cursor.y = screenHeight - 1
		buf.CursorRender()
	}

	if buf.Cursor.x < 0 {
		buf.Cursor.x = 0
	} else if buf.Cursor.x >= len(buf.data[buf.render_y]) {
		buf.Cursor.x = len(buf.data[buf.render_y]) - 1
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
