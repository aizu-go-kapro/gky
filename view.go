package main

type View struct {
	buf  *Buffer
	mode Mode
}

func NewView(buf *Buffer) *View {
	return &View{
		buf:  buf,
		mode: Normal,
	}
}

func (v *View) Init() error {
	v.buf.Render(0)
	v.buf.setCursor()

	return nil
}
