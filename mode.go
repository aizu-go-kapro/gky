package main

type Mode int

const (
	Normal Mode = iota
	Insert
	Visual
)

func (m Mode) String() string {
	switch m {
	case Normal:
		return "NORMAL"
	case Insert:
		return "INSERT"
	case Visual:
		return "VISUAL"
	}
	return "non mode"
}
