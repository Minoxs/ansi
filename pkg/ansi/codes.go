package ansi

import (
	"errors"
	"io"
	"strconv"
	"strings"
)

const (
	EscapeCode = 27
	StartCode  = byte('[') // 91
	EndCode    = byte('m') // 109
	Separator  = ';'       // 59
)

// ErrByteOverflow indicates that the parsed value is under 0 or above 255.
var ErrByteOverflow = errors.New("invalid byte value")

// isSequenceStart is a helper function to check if the data slice starts at an ANSI escape sequence
func isSequenceStart(data []byte) bool {
	return len(data) >= 2 && data[0] == EscapeCode && data[1] == StartCode
}

// ScanCodes will tokenize a string into ANSI colour sequences and normal text.
// Sequences can be checked with ansi.IsSequence.
// This function implements bufio.SplitFunc.
func ScanCodes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		err = io.EOF
		return
	}

	advance = len(data)

	if isSequenceStart(data) {
		for i := 0; i < advance; i++ {
			if data[i] == EndCode {
				advance = i + 1
				break
			}
		}
	} else {
		for i := 0; i < advance; i++ {
			if isSequenceStart(data[i:]) {
				advance = i
				break
			}
		}
	}

	token = data[:advance]
	return
}

// IsSequence checks wheter the given slice is a complete ANSI colour sequence.
// It is meant to be used along with ansi.ScanCodes.
func IsSequence(b []byte) bool {
	return isSequenceStart(b) && b[len(b)-1] == EndCode
}

// ParseColourCodes will parse Colour codes from an ANSI sequence.
// It is meant to be used along with ansi.ScanCodes and ansi.IsSequence.
func ParseColourCodes(sequence []byte) (res []Colour, err error) {
	var strCodes = strings.Split(string(sequence[2:len(sequence)-1]), string(Separator))
	res = make([]Colour, len(strCodes))

	for i := 0; i < len(res); i++ {
		var n int

		n, err = strconv.Atoi(strCodes[i])
		if err != nil {
			return
		}

		if n > 255 || n < 0 {
			err = ErrByteOverflow
			return
		}

		res[i] = Colour(n)
	}

	return
}

// MustParseColourCodes will parse Colour codes from an ANSI sequence.
// It is meant to be used along with ansi.ScanCodes and ansi.IsSequence.
// Panics if sequence is not valid.
func MustParseColourCodes(sequence []byte) []Colour {
	var res, err = ParseColourCodes(sequence)
	if err != nil {
		panic(err)
	}
	return res
}
