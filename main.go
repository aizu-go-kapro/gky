package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	homedir "github.com/mitchellh/go-homedir"
)

type Mode int

const (
	Normal Mode = iota
	Insert
	Visual
)

type Cursor struct {
	x     int
	y     int
	max_x int
	max_y int
}

type FileManager struct {
	width    int
	height   int
	namepath string
	file     *os.File
}

type Editor struct {
	cursor Cursor
	mode   Mode
	fm     FileManager
	screen tcell.Screen
}

func (e *Editor) Open(filename string) error {
	name, err := homedir.Expand(filename)
	if err != nil {
		return err
	}
	e.fm.namepath = name

	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := os.Stat(name)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("%s is a directory", name)
	}
	e.fm.file = f

	return nil
}

//TODO init screen init
func (e *Editor) Init() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	e.screen = s

	if err := e.screen.Init(); err != nil {
		return err
	}
	return nil
}

func NewEditor() *Editor {
	return &Editor{
		cursor: Cursor{x: 0, y: 0},
		mode:   Normal,
		fm:     FileManager{},
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: command <filename>\n")
		os.Exit(1)
	}

	e := NewEditor()
	if err := e.Open(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if err := e.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer e.screen.Fini()

	quit := make(chan struct{})
	go func() {
		for {
			ev := e.screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape {
					close(quit)
				} else if ev.Key() == tcell.KeyCtrlC {
					close(quit)
				}
			default:
			}
		}
	}()

loop:
	for {
		select {
		case <-quit:
			break loop
		}
	}
}
