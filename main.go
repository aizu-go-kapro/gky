package main

import (
	"fmt"
	"math"
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

type Cursor struct {
	x int32
	y int32
}

type FileManager struct {
	namepath string
	file     *os.File
	bytes    []byte
	size     int
}

type Window struct {
	width  int
	height int
	full   int
	screen tcell.Screen
}

type Editor struct {
	cursor Cursor
	mode   Mode
	fm     FileManager
	win    Window
}

const (
	BUFSIZE = math.MaxInt32
)

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
	if err := e.fm.Read(); err != nil {
		return err
	}

	return nil
}

func (f *FileManager) Read() error {
	buf := make([]byte, BUFSIZE)
	n, err := f.file.Read(buf)
	if err != nil {
		return err
	}
	f.size = n
	f.bytes = buf[:n]

	return nil
}

func (w *Window) Show(bytes []byte) {
	var max_win int
	w.screen.Clear()
	count := 0
	h := 0
	max_win = w.full
	if w.full > len(bytes) {
		max_win = len(bytes)
	}

	for i := 0; i < max_win; i++ {
		if w.width == count || bytes[i] == byte(10) {
			h++
			count = 0
			w.screen.SetContent(count, h, rune(bytes[i]), nil, tcell.StyleDefault)
			continue
		}
		w.screen.SetContent(count, h, rune(bytes[i]), nil, tcell.StyleDefault)
		count++
	}
	w.screen.Show()
}

//TODO
func (e *Editor) Init() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	e.win.screen = s
	if err := e.win.screen.Init(); err != nil {
		return err
	}

	w, h := s.Size()
	e.win.width = w - 1
	e.win.height = h - 1
	e.win.full = e.win.width * e.win.height

	e.win.Show(e.fm.bytes)
	return nil
}

func NewEditor() *Editor {
	return &Editor{
		cursor: Cursor{x: 0, y: 0},
		mode:   Normal,
		fm:     FileManager{},
		win:    Window{},
	}
}

func main() {
	encoding.Register()
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

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
	defer e.win.screen.Fini()

	quit := make(chan struct{})
	go func() {
		for {
			ev := e.win.screen.PollEvent()
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
