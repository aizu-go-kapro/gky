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
	data     [][]rune
	Name     string
	path     string
	Cursor   *Location
	render_y int
}

func NewLocation(l, c int) *Location {
	return &Location{x: l, y: c}
}

// by tenntenn *os.Fileではなく、io.ReadSeekerの方がいいかも
func (buf *Buffer) Read(f *os.File) error {
	// by tenntenn バッファリングはするべきだけど、全部のデータをメモリ上に展開するとでかいファイルをひらけなくなるのでは？エディタを作ったことがないのでなんとも言えないけど
	b := make([]byte, BUFSIZE)
	// by tenntenn 全部読み込むならioutil.ReadAllを使ったほうがいい。
	n, err := f.Read(b)
	if err != nil {
		return err
	}
	b = b[:n]

	var s []rune
	// by tenntenn bufio.Scannerを使ったほうが良さそう
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
	// by tenntenn os.Openなので、新規作成ができない
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

	for i := 0; i < h; i++ {
		w := len(buf.data[from])
		if screenWidth < w {
			w = screenWidth
		}
		for j := 0; j < w; j++ {
			screen.SetContent(j, i, buf.data[from][j], nil, tcell.StyleDefault)
		}
		from++
	}
}

func (buf *Buffer) getLine() int {
	return len(buf.data)
}
