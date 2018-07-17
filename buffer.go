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
	data     [][]rune
	Name     string
	path     string
	Cursor   *Location
	render_y int
}

func NewLocation(l, c int) *Location {
	return &Location{x: l, y: c}
}

func (buf *Buffer) Insert(data []rune) {
	x, y := buf.Cursor.x, buf.Cursor.y

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

	x, y := buf.Cursor.x, buf.Cursor.y

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
	line_num = from

	for i := 0; i < h; i++ {
		w := len(buf.data[from])

		line := len(fmt.Sprintf("%d", buf.getLine()))
		line += 1

		l := fmt.Sprintf("%d", line_num+1)
		for line_i, v := range l {
			screen.SetContent(0+line_i, i, rune(v), nil, tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorSlateGray))
			//TODO: ↓ fix
			if len(fmt.Sprintf("%d", buf.getLine())) == 3 {
				switch line_i {
				case 0:
					screen.SetContent(0+line_i+1, i, ' ', nil, tcell.StyleDefault.Background(tcell.ColorSlateGray))
					screen.SetContent(0+line_i+2, i, ' ', nil, tcell.StyleDefault.Background(tcell.ColorSlateGray))
					screen.SetContent(0+line_i+3, i, ' ', nil, tcell.StyleDefault.Background(tcell.ColorSlateGray))
				case 1:
					screen.SetContent(0+line_i+1, i, ' ', nil, tcell.StyleDefault.Background(tcell.ColorSlateGray))
					screen.SetContent(0+line_i+2, i, ' ', nil, tcell.StyleDefault.Background(tcell.ColorSlateGray))
				case 2:
					screen.SetContent(0+line_i+1, i, ' ', nil, tcell.StyleDefault.Background(tcell.ColorSlateGray))
				}
			}
		}

		for j := 0; j < w; j++ {
			screen.SetContent(j+line, i, buf.data[from][j], nil, tcell.StyleDefault.Foreground(tcell.ColorAqua))
		}
		from++
		line_num++
	}
}

func (buf *Buffer) getLine() int {
	return len(buf.data)
}
