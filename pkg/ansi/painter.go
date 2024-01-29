package ansi

import (
	"image/color"
)

// Painter is a struct that contains a text (foreground) and background colours.
// Can be created from an ANSI sequence and has helper methods to read sequence.
// Nil values means the colour should be ignored and a default used.
type Painter struct {
	text color.Color
	back color.Color

	Bold   bool
	Italic bool
}

// NewPainter creates a painter with the given foreground and background colors.
func NewPainter(text color.Color, back color.Color) *Painter {
	return &Painter{
		text: text,
		back: back,
	}
}

// DefaultPainter is a helper function that returns a painter with default values.
func DefaultPainter() *Painter {
	return &Painter{}
}

// NewPainterFromSequence creates an ANSI sequence, as parsed by ansi.ScanCodes.
// Panics if sequence is invalid.
func NewPainterFromSequence(sequence []byte) (p *Painter) {
	p = DefaultPainter()

	var codes = MustParseColourCodes(sequence)
	for _, code := range codes {
		switch code {
		case NORMAL:
			continue
		case BOLD:
			p.Bold = true
		case ITALIC:
			p.Italic = true
		default:
			switch {
			case code.IsText():
				p.text = code
			case code.IsBackground():
				p.back = code
			}
		}
	}

	return
}

// TextColor returns the underlying foreground color.
func (p *Painter) TextColor() color.Color {
	return p.text
}

// BackgroundColor returns the underlying background color.
func (p *Painter) BackgroundColor() color.Color {
	return p.back
}
