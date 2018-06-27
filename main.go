package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

var (
	Events                    = make(chan tcell.Event)
	screen       tcell.Screen = nil
	screenWidth  int
	screenHeight int
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
