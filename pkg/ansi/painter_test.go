package ansi

import "testing"

func TestPainter_NewPainterFromSequence(t *testing.T) {
	var p = NewPainterFromSequence(
		[]byte{
			EscapeCode, StartCode,
			byte('1'), Separator,
			byte('2'), Separator,
			byte('3'), byte('1'), Separator,
			byte('4'), byte('2'), Separator,
			byte('9'), byte('1'),
			EndCode,
		},
	)

	if !p.Bold {
		t.Fatal("Failed to detect BOLD")
	}

	if !p.Italic {
		t.Fatal("Failed to detect ITALIC")
	}

	if p.TextColor().(Colour) != RED+HighIntensityOffset {
		t.Fatal("Failed to detect foreground colour")
	}

	if p.BackgroundColor().(Colour) != GREEN+BackgroundOffset {
		t.Fatal("Failed to detect background colour")
	}
}
