package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	homedir "github.com/mitchellh/go-homedir"
)

type Mode int

const (
	Normal Mode = iota
	Insert
	Visual
)

type FileInfo struct {
	namepath string
	file     *os.File
	contents []string
}

type Cursor struct {
	x     int
	y     int
	max_x int
	max_y int
}

type Manager struct {
	width  int
	height int
	file   FileInfo
}

type Editor struct {
	cursor  Cursor
	mode    Mode
	manager Manager
	screen  tcell.Screen
}

func (e *Editor) Open(filename string) error {
	name, err := homedir.Expand(filename)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var str []string
	for scanner.Scan() {
		str = append(str, scanner.Text()+"\n")
	}
	if err := scanner.Err(); err != nil {
		return errors.New(fmt.Sprintf("scanner err:", err))
	}

	e.manager.file = FileInfo{
		namepath: name,
		file:     file,
		contents: str,
	}

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
		cursor:  Cursor{x: 0, y: 0},
		mode:    Normal,
		manager: Manager{},
	}
}

func main() {
	encoding.Register()
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	e := NewEditor()

	if err := e.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer e.screen.Fini()

	if len(os.Args) != 2 {
		fmt.Printf("Usage: command <filename>\n")
		os.Exit(1)
	}

	if err := e.Open(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

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
