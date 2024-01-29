package ansi

import "testing"

func TestColour_Normalize(t *testing.T) {
	for testColour := WHITE; testColour <= BLACK; testColour++ {
		var c = testColour + BackgroundOffset + HighIntensityOffset
		if c.Normalize() != testColour {
			t.Fatal("Colour normalization is broken")
		}
	}
}

func TestColour_IsText(t *testing.T) {
	for testColour := WHITE; testColour <= BLACK; testColour++ {
		if !testColour.IsText() {
			t.Fatal("Failed to detect colour from foreground")
		}
	}
	for testColour := WHITE + BackgroundOffset; testColour <= BLACK+BackgroundOffset; testColour++ {
		if testColour.IsText() {
			t.Fatal("Failed to detect colour from background")
		}
	}
	for testColour := WHITE + HighIntensityOffset; testColour <= BLACK+HighIntensityOffset; testColour++ {
		if !testColour.IsText() {
			t.Fatal("Failed to detect colour from foreground high intensity")
		}
	}
	for testColour := WHITE + BackgroundOffset + HighIntensityOffset; testColour <= BLACK+BackgroundOffset+HighIntensityOffset; testColour++ {
		if testColour.IsText() {
			t.Fatal("Failed to detect colour from background high intensity")
		}
	}
}

func TestColour_IsBackground(t *testing.T) {
	for testColour := WHITE; testColour <= BLACK; testColour++ {
		if testColour.IsBackground() {
			t.Fatal("Failed to detect colour from foreground")
		}
	}
	for testColour := WHITE + BackgroundOffset; testColour <= BLACK+BackgroundOffset; testColour++ {
		if !testColour.IsBackground() {
			t.Fatal("Failed to detect colour from background")
		}
	}
	for testColour := WHITE + HighIntensityOffset; testColour <= BLACK+HighIntensityOffset; testColour++ {
		if testColour.IsBackground() {
			t.Fatal("Failed to detect colour from foreground high intensity")
		}
	}
	for testColour := WHITE + BackgroundOffset + HighIntensityOffset; testColour <= BLACK+BackgroundOffset+HighIntensityOffset; testColour++ {
		if !testColour.IsBackground() {
			t.Fatal("Failed to detect colour from background high intensity")
		}
	}
}

func TestColour_IsHighIntensity(t *testing.T) {
	for testColour := WHITE; testColour <= BLACK; testColour++ {
		if testColour.IsHighIntensity() {
			t.Fatal("Failed to detect colour from foreground")
		}
	}
	for testColour := WHITE + BackgroundOffset; testColour <= BLACK+BackgroundOffset; testColour++ {
		if testColour.IsHighIntensity() {
			t.Fatal("Failed to detect colour from background")
		}
	}
	for testColour := WHITE + HighIntensityOffset; testColour <= BLACK+HighIntensityOffset; testColour++ {
		if !testColour.IsHighIntensity() {
			t.Fatal("Failed to detect colour from foreground high intensity")
		}
	}
	for testColour := WHITE + BackgroundOffset + HighIntensityOffset; testColour <= BLACK+BackgroundOffset+HighIntensityOffset; testColour++ {
		if !testColour.IsHighIntensity() {
			t.Fatal("Failed to detect colour from background high intensity")
		}
	}
}
