package main

import "github.com/gdamore/tcell"

type View struct {
	buf  *Buffer
	mode Mode
}

func NewView(buf *Buffer) *View {
	return &View{
		buf:  buf,
		mode: Normal,
	}
}

func (v *View) render() {
	h := len(v.buf.data)
	if screenHeight < h {
		h = screenHeight
	}

	for i := 0; i < h; i++ {
		w := len(v.buf.data[i])
		if screenWidth < w {
			w = screenWidth
		}
		for j := 0; j < w; j++ {
			screen.SetContent(j, i, v.buf.data[i][j], nil, tcell.StyleDefault)
		}
	}

	//TODO render時のcursor位置を0,0では無くす
	v.buf.setCursor()
	screen.Show()
}
