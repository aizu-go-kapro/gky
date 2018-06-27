package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

// by tenntenn 構造体にすることでパッケージ変数を減らせないか考えてみてくdさい
var (
	Events = make(chan tcell.Event)
	// by tenntenn 変数は明示的に初期化しなくてもゼロ値で初期化されている
	screen       tcell.Screen = nil
	screenWidth               = 0
	screenHeight              = 0
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	if len(args) != 2 {
		fmt.Fprint(os.Stderr, "Usage: gky <filename>\n")
		return 1
	}

	if err := initScreen(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer screen.Fini() // by tenntenn 略さない方がいい

	initEvent()

	b, err := initBuffer(args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	view := NewView(b)
	if err := view.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// by tenntenn 別の関数にしてreturnで抜けたほうがわかりやすいかも
loop:
	for {
		select {
		case ev := <-Events:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlQ {
					break loop
				} else if err := view.EventHandle(ev); err != nil {
					fmt.Fprintln(os.Stderr, err)
					break loop
				}
			case *tcell.EventResize:
				screenWidth, screenHeight = screen.Size()
				view.buf.Render(0)
			}
		}
	}
	return 0
}
