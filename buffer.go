package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gdamore/tcell"
	homedir "github.com/mitchellh/go-homedir"
)

type Location struct {
	x int
	y int
}

type Buffer struct {
	data            [][]rune
	Name            string
	path            string
	Cursor          *Location
	HighlightBegine *Location
	HighlightEnd    *Location
	renderY         int
	exist           bool
}

func NewLocation(l, c int) *Location {
	return &Location{x: l, y: c}
}

func (buf *Buffer) Insert(data []rune) {
	x, y := buf.Cursor.x-lineLen, buf.Cursor.y

	for i := len(data) - 1; i >= 0; i-- {
		char := data[i]
		if byte(char) == '\n' {
			tmp := append([]rune(nil), buf.data[y][x:]...)
			buf.data[y] = buf.data[y][:x]
			buf.data = append(buf.data[:y+1], append([][]rune{tmp}, buf.data[y+1:]...)...)
		} else {
			buf.data[y] = append(buf.data[y][:x], append([]rune{char}, buf.data[y][x:]...)...)
		}
	}
}

func (buf *Buffer) Remove(length int) int {
	lastCursor := 0

	x, y := buf.Cursor.x-lineLen, buf.Cursor.y

	for i := 0; i < length; i++ {
		if x == 0 {
			n := len(buf.data[y-1])
			restLength := screenWidth - n

			if restLength > len(buf.data[y]) {
				lastCursor = n + 1
				// delete extra value
				if n == 1 && byte(buf.data[y-1][0]) == '\n' {
					buf.data[y-1] = buf.data[y-1][:len(buf.data[y-1])-1]
				}

				buf.data[y-1] = append(buf.data[y-1], buf.data[y]...)
				buf.data = append(buf.data[:y], buf.data[y+1:]...)
			} else {
				// insert
				buf.data[y-1] = append(buf.data[y-1], buf.data[y][:restLength]...)

				// [TODO] delete from buf.data[y-1]

				buf.data = append(buf.data[:y], buf.data[y+1:]...)
			}
		} else {
			buf.data[y] = append(buf.data[y][:x-1], buf.data[y][x:]...)
			lastCursor = 0
		}
	}

	return lastCursor
}

func (buf *Buffer) Read(f *os.File) error {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

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
	if buf.data == nil {
		buf.data = append(buf.data, []rune{' '})
	}

	return nil
}

func (buf *Buffer) fileManage(path string) error {
	var f *os.File

	namepath, err := homedir.Expand(path)
	if err != nil {
		return err
	}
	buf.path = namepath

	info, err := os.Stat(buf.path)
	if err != nil {
		f, err = os.OpenFile(buf.path, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		buf.exist = false
	} else {
		f, err = os.Open(buf.path)
		if err != nil {
			return err
		}
		buf.exist = true
	}
	defer f.Close()
	buf.Name = f.Name()

	info, _ = os.Stat(buf.path)
	if info.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}

	if err := buf.Read(f); err != nil {
		return err
	}

	return nil
}

func (buf *Buffer) Render(from int) {
	screen.Clear()
	buf.render(from)
	screen.Show()
}

func (buf *Buffer) render(from int) {
	h := len(buf.data[from:])
	if screenHeight < h {
		h = screenHeight
	}
	lineNum := from

	lineLen = len(fmt.Sprintf("%d", buf.getLine())) + 1

	for i := 0; i < h; i++ {
		for j := lineLen; j > 0; j-- {
			screen.SetContent(0+j, i, ' ', nil, tcell.StyleDefault.Background(tcell.ColorSlateGray))
		}
	}

	for i := 0; i < h; i++ {
		w := len(buf.data[from])

		l := fmt.Sprintf("%d", lineNum+1)
		for lineI, v := range l {
			screen.SetContent(0+lineI, i, rune(v), nil, tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorSlateGray))
		}

		for j := 0; j < w; j++ {
			screen.SetContent(j+lineLen, i, buf.data[from][j], nil, tcell.StyleDefault.Foreground(tcell.ColorAqua))
		}
		from++
		lineNum++
	}
}

func (buf *Buffer) getLine() int {
	return len(buf.data)
}

func (buf *Buffer) Save() error {
	var bytes []byte
	for i := range buf.data {
		for _, v := range buf.data[i] {
			bytes = append(bytes, byte(v))
		}
		bytes = append(bytes, byte('\n'))
	}
	if err := ioutil.WriteFile(buf.path, bytes, 0644); err != nil {
		return err
	}
	buf.exist = true
	return nil
}
