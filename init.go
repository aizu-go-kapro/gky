package main

import (
	"fmt"

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
	screen.Resize(0, 0, screenWidth, screenHeight)

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

	lineLen = len(fmt.Sprintf("%d", buf.getLine())) + 1
	buf.Cursor = NewLocation(lineLen, 0)
	buf.HighlightBegine = NewLocation(-1, -1)
	buf.HighlightEnd = NewLocation(-1, -1)

	return buf, nil
}
