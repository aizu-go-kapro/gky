package main

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

// by tenntenn 初期化処理だとしても、極力パッケージ変数の初期化はしない
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
	// by tenntenn この-2はなんだろう？
	screen.Resize(0, 0, screenWidth, screenHeight-2)

	return nil
}

func initEvent() {
	// by tenntenn このゴールーチンを終わらせるための仕組みがない
	// See: https://qiita.com/tenntenn/items/dd6041d630af7feeec52
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
	buf := new(Buffer) // by tenntenn var buf Bufferでも使える
	if err := buf.fileManage(path); err != nil {
		return nil, err
	}
	buf.Cursor = NewLocation(0, 0)

	return buf, nil // by tenntenn var buf Bufferとした場合は&buf
}
