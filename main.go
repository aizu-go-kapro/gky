package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

var (
	Events                    = make(chan tcell.Event)
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
		fmt.Fprint(os.Stderr, err)
		return 1
	}
	defer screen.Fini()

	initEvent()

	if err := initBuffer(args[1]); err != nil {
		fmt.Fprint(os.Stderr, err)
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
				}
			}
		}
	}
	return 0
}
