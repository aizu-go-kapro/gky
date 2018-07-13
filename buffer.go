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
