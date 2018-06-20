package main

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

func initScreen() error {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}

	if err = screen.Init(); err != nil {
		return err
	}

	encoding.Register()
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen.SetStyle(tcell.StyleDefault)
	screen.Clear()

	screenWidth, screenHeight = screen.Size()

	return nil
}

func initEvent() {
	go func() {
		for {
			if screen == nil {
				break
			}
			Events <- screen.PollEvent()
		}
	}()
}

func initBuffer(path string) (*Buffer, error) {
	buf := new(Buffer)
	if err := buf.fileManage(path); err != nil {
		return nil, err
	}
	buf.Cursor = NewLocation(0, 0)

	return buf, nil
}

func (buf *Buffer) initView() error {
	h := len(buf.data)
	if screenHeight < h {
		h = screenHeight
	}

	for i := 0; i < h; i++ {
		w := len(buf.data[i])
		if screenWidth < w {
			w = screenWidth
		}
		for j := 0; j < w; j++ {
			screen.SetContent(j, i, buf.data[i][j], nil, tcell.StyleDefault)
		}
	}

	screen.ShowCursor(buf.Cursor.x, buf.Cursor.y)
	screen.Show()
	return nil
}
