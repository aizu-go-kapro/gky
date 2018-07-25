package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

var (
	Events       = make(chan tcell.Event)
	screen       tcell.Screen
	screenWidth  int
	screenHeight int
	lineLen      int
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: gky <filename>")
		return 1
	}

	if err := initScreen(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer screen.Fini()

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

	for ev := range Events {
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlQ {
				if !view.buf.exist {
					os.Remove(view.buf.path)
				}
				break
			} else if err := view.EventHandle(ev); err != nil {
				fmt.Fprintln(os.Stderr, err)
				break
			}
		case *tcell.EventResize:
			screenWidth, screenHeight = ev.Size()
			view.buf.Render(0)
		}
	}
	return 0
}
