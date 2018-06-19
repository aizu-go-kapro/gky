package main

import (
	"fmt"
	"math"
	"os"

	"github.com/gdamore/tcell"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	BUFSIZE = math.MaxInt32
)

type Location struct {
	x int
	y int
}

type Buffer struct {
	data   [][]rune
	Name   string
	path   string
	Cursor *Location
}

func (buf *Buffer) Read(f *os.File) error {
	b := make([]byte, BUFSIZE)
	n, err := f.Read(b)
	if err != nil {
		return err
	}
	b = b[:n]

	var s []rune
	for _, v := range b {
		if v == byte(10) { // 改行のバイト
			s = append(s, rune(v))
			buf.data = append(buf.data, s)
			s = nil
			continue
		}
		s = append(s, rune(v))
	}

	return nil
}

func (buf *Buffer) open() error {
	f, err := os.Open(buf.path)
	if err != nil {
		return err
	}
	defer f.Close()

	buf.Name = f.Name()

	if err := buf.Read(f); err != nil {
		return err
	}

	return nil
}

func (buf *Buffer) fileManage(path string) error {
	namepath, err := homedir.Expand(path)
	if err != nil {
		return err
	}
	buf.path = namepath

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}

	if err := buf.open(); err != nil {
		return err
	}

	return nil
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
	screen.Show()
	return nil
}

func initBuffer(path string) (*Buffer, error) {
	buf := new(Buffer)
	if err := buf.fileManage(path); err != nil {
		return nil, err
	}

	return buf, nil
}
