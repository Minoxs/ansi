package ansi

func (c Colour) String() (res string) {
	switch c.Normalize() {
	case NORMAL:
		return "NORMAL"
	case BOLD:
		return "BOLD"
	case ITALIC:
		return "ITALIC"
	case BLACK:
		res = "BLACK"
	case RED:
		res = "RED"
	case GREEN:
		res = "GREEN"
	case YELLOW:
		res = "YELLOW"
	case BLUE:
		res = "BLUE"
	case PURPLE:
		res = "PURPLE"
	case CYAN:
		res = "CYAN"
	case WHITE:
		res = "WHITE"
	}

	if c.IsBackground() {
		res += "_BACKGROUND"
	}

	if c.IsHighIntensity() {
		res += "_BRIGHT"
	}

	return
}
