package ansi

type Colour byte

const (
	BackgroundOffset    = 10
	HighIntensityOffset = 60
)

const (
	NORMAL Colour = 0
	BOLD   Colour = 1
	ITALIC Colour = 2
	BLACK  Colour = 30
	RED    Colour = 31
	GREEN  Colour = 32
	YELLOW Colour = 33
	BLUE   Colour = 34
	PURPLE Colour = 35
	CYAN   Colour = 36
	WHITE  Colour = 37
)

// RGBA8 returns the RGB representation of the ansi colours.
// Return value is in the sensible way any human being with a brain would use.
// Colour scheme used is Windows 10's Campbell colour scheme as described in https://en.wikipedia.org/wiki/ANSI_escape_code.
func (c Colour) RGBA8() (r, g, b, a uint32) {
	if c.IsHighIntensity() {
		return c.RGBA8Bright()
	} else {
		return c.RGBA8Normal()
	}
}

// RGBA8Normal returns the RGBA8 color with normal intensity.
// Colour scheme used is Windows 10's Campbell colour scheme as described in https://en.wikipedia.org/wiki/ANSI_escape_code.
func (c Colour) RGBA8Normal() (r, g, b, a uint32) {
	a = 255

	switch c.Normalize() {
	case BLACK:
		return 12, 12, 12, a
	case RED:
		return 197, 15, 31, a
	case GREEN:
		return 19, 161, 14, a
	case YELLOW:
		return 193, 156, 0, a
	case BLUE:
		return 0, 55, 218, a
	case PURPLE:
		return 136, 23, 152, a
	case CYAN:
		return 58, 150, 221, a
	case WHITE:
		return 204, 204, 204, a
	default:
		return
	}
}

// RGBA8Bright returns the RGBA color with high intensity.
// Colour scheme used is Windows 10's Campbell colour scheme as described in https://en.wikipedia.org/wiki/ANSI_escape_code.
func (c Colour) RGBA8Bright() (r, g, b, a uint32) {
	a = 255

	switch c.Normalize() {
	case BLACK:
		return 118, 118, 118, a
	case RED:
		return 231, 72, 86, a
	case GREEN:
		return 22, 198, 12, a
	case YELLOW:
		return 249, 241, 165, a
	case BLUE:
		return 59, 120, 255, a
	case PURPLE:
		return 180, 0, 158, a
	case CYAN:
		return 97, 214, 214, a
	case WHITE:
		return 242, 242, 242, a
	default:
		return
	}
}

// RGBA implements the color.Color interface.
// Returns value in the same stupid way color.RGBA does.
// Colour scheme used is Windows 10's Campbell colour scheme as described in https://en.wikipedia.org/wiki/ANSI_escape_code.
func (c Colour) RGBA() (r, g, b, a uint32) {
	r, g, b, a = c.RGBA8()
	r |= r << 8
	g |= g << 8
	b |= b << 8
	a |= a << 8
	return
}

// Normalize returns the Colour in the BLACK to WHITE range.
func (c Colour) Normalize() Colour {
	if c > WHITE {
		return BLACK + (c % 10)
	}
	return c
}

// IsText returns if the Colour is a foreground colour.
func (c Colour) IsText() bool {
	if c.IsHighIntensity() {
		c -= HighIntensityOffset
	}
	return WHITE >= c && c >= BLACK
}

// IsBackground returns if the Colour is a background colour.
func (c Colour) IsBackground() bool {
	if c.IsHighIntensity() {
		c -= HighIntensityOffset
	}
	return WHITE+BackgroundOffset >= c && c >= BLACK+BackgroundOffset
}

// IsHighIntensity returns if the Colour is supposed to be high intensity.
func (c Colour) IsHighIntensity() bool {
	return c >= BLACK+HighIntensityOffset
}
