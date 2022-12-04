package asciiui

import (
	"golang.org/x/term"
	"os"
)

func NewWindow() (*Element, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, err
	}

	return &Element{width}, nil
}

type Element struct {
	Width int
}

type Renderable interface {
	Render() string
}
