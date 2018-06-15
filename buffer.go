package main

import (
	"fmt"
	"math"
	"os"

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
	data   [][]byte
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

	var s []byte
	for _, v := range b {
		if v == byte(10) { // 改行のバイト
			s = append(s, v)
			buf.data = append(buf.data, s)
			s = nil
			continue
		}
		s = append(s, v)
	}

	return nil
}

func (buf *Buffer) Open() error {
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

	if err := buf.Open(); err != nil {
		return err
	}

	return nil
}

func initBuffer(path string) error {
	buf := new(Buffer)
	if err := buf.fileManage(path); err != nil {
		return nil
	}

	return nil
}
