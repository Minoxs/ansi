package ansi

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func Test_isSequenceStart(t *testing.T) {
	t.Run(
		"CheckValid", func(t *testing.T) {
			var b = []byte{27, 91}
			if !isSequenceStart(b) {
				t.Fatal("Sequence Should be valid")
			}
		},
	)
	t.Run(
		"CheckInvalid", func(t *testing.T) {
			var b = []byte{91, 27, 91}
			if isSequenceStart(b) {
				t.Fatal("Sequence Should be invalid")
			}
		},
	)
}

func Test_ScanCodes(t *testing.T) {
	t.Run(
		"Normal Text", func(t *testing.T) {
			var b = []byte{72, 101, 108, 108, 111, 91, 44, 32, 87, 111, 114, 108, 100, 33}
			var advance, token, err = ScanCodes(b, false)
			if err != nil {
				t.Fatal("Errored while scanning normal text")
			}

			if advance != len(b) {
				t.Fatal("Scan isn't properly consuming text")
			}

			if string(token) != string(b) {
				t.Fatal("Error parsing normal text")
			}
		},
	)

	t.Run(
		"Text with sequence", func(t *testing.T) {
			var b = []byte{
				EscapeCode, StartCode, byte('3'), byte('1'), EndCode,
				72, 101, 108, 108, 111, 91, 44, 32, 87, 111, 114, 108, 100, 33,
			}

			var advance, token, err = ScanCodes(b, false)
			if err != nil {
				t.Fatal("Errored while scanning text with sequence")
			}

			if advance != 5 {
				t.Fatal("Scan isn't properly consuming sequence")
			}

			if string(token) != string(b[0:5]) {
				t.Fatal("Error parsing sequence")
			}

			b = b[advance:]
			advance, token, err = ScanCodes(b, false)
			if err != nil {
				t.Fatal("Errored while scanning normal text")
			}

			if advance != len(b) {
				t.Fatal("Scan isn't properly consuming text")
			}

			if string(token) != string(b) {
				t.Fatal("Error parsing normal text")
			}
		},
	)
}

func Test_IsSequence(t *testing.T) {
	if IsSequence([]byte{StartCode, byte('3'), byte('1'), EndCode}) {
		t.Fatal("Sequence isn't being escaped")
	}

	if !IsSequence([]byte{EscapeCode, StartCode, byte('3'), byte('1'), EndCode}) {
		t.Fatal("Sequence isn't being properly detected")
	}

	if IsSequence([]byte{EscapeCode, StartCode, byte('3'), byte('1')}) {
		t.Fatal("Sequence end not being checked")
	}
}

func Test_ParseColourCodes(t *testing.T) {
	t.Run(
		"Simple", func(t *testing.T) {
			var (
				expected = []Colour{NORMAL}
				actual   = MustParseColourCodes([]byte{EscapeCode, StartCode, byte('0'), EndCode})
			)

			for i := 0; i < len(expected); i++ {
				if expected[i] != actual[i] {
					t.Fatal("Colour code parsed incorrectly")
				}
			}
		},
	)

	t.Run(
		"Separated", func(t *testing.T) {
			var (
				expected = []Colour{NORMAL, RED}
				actual   = MustParseColourCodes(
					[]byte{
						EscapeCode, StartCode, byte('0'), Separator, byte('3'), byte('1'), EndCode,
					},
				)
			)

			for i := 0; i < len(expected); i++ {
				if expected[i] != actual[i] {
					t.Fatal("Colour code parsed incorrectly")
				}
			}
		},
	)

	t.Run(
		"Complex", func(t *testing.T) {
			var (
				expected = []Colour{
					NORMAL, RED,
					BLUE + BackgroundOffset, RED + HighIntensityOffset,
					BLUE + BackgroundOffset + HighIntensityOffset,
				}
				actual = MustParseColourCodes(
					[]byte{
						EscapeCode, StartCode,
						byte('0'), Separator,
						byte('3'), byte('1'), Separator,
						byte('4'), byte('4'), Separator,
						byte('9'), byte('1'), Separator,
						byte('1'), byte('0'), byte('4'), EndCode,
					},
				)
			)

			for i := 0; i < len(expected); i++ {
				if expected[i] != actual[i] {
					t.Fatal("Colour code parsed incorrectly")
				}
			}
		},
	)
}

type benchScan struct {
	name string
	test []byte
}

var benches = []benchScan{
	{
		name: "Simple",
		// Red Hello world
		// With reset at the end
		test: []byte(
			// Red
			string([]byte{EscapeCode, StartCode, byte('3'), byte('1'), EndCode}) +
				// Text
				"Hello, World!" +
				// Reset
				string([]byte{EscapeCode, StartCode, byte('0'), EndCode}) +
				// NewLine
				"\n",
		),
	},
	{
		name: "Complex",
		test: []byte{
			10, 32, 32, 27, 91, 48, 59, 49, 59, 51, 52, 59, 57, 52, 109, 95, 95, 95, 27, 91, 48, 109, 32, 32, 27, 91,
			48, 59, 49, 59, 51, 52, 59, 57, 52, 109, 95, 95, 95, 95, 27, 91, 48, 109, 32, 32, 27, 91, 48, 59, 49, 59,
			51, 52, 59, 57, 52, 109, 95, 27, 91, 48, 109, 32, 32, 32, 27, 91, 48, 59, 51, 52, 109, 95, 95, 95, 95, 27,
			91, 48, 109, 32, 32, 32, 27, 91, 48, 59, 51, 52, 109, 95, 95, 95, 95, 27, 91, 48, 109, 32, 10, 32, 27, 91,
			48, 59, 49, 59, 51, 52, 59, 57, 52, 109, 47, 27, 91, 48, 109, 32, 27, 91, 48, 59, 49, 59, 51, 52, 59, 57,
			52, 109, 95, 27, 91, 48, 109, 32, 27, 91, 48, 59, 49, 59, 51, 52, 59, 57, 52, 109, 92, 124, 27, 91, 48, 109,
			32, 32, 27, 91, 48, 59, 51, 52, 109, 95, 27, 91, 48, 109, 32, 27, 91, 48, 59, 51, 52, 109, 92, 40, 95, 41,
			27, 91, 48, 109, 32, 27, 91, 48, 59, 51, 52, 109, 124, 27, 91, 48, 109, 32, 32, 27, 91, 48, 59, 51, 52, 109,
			95, 27, 91, 48, 109, 32, 27, 91, 48, 59, 51, 52, 109, 92, 27, 91, 48, 109, 32, 27, 91, 48, 59, 51, 52, 109,
			47, 27, 91, 48, 109, 32, 27, 91, 48, 59, 51, 55, 109, 95, 95, 95, 124, 27, 91, 48, 109, 10, 27, 91, 48, 59,
			51, 52, 109, 124, 27, 91, 48, 109, 32, 27, 91, 48, 59, 51, 52, 109, 124, 27, 91, 48, 109, 32, 27, 91, 48,
			59, 51, 52, 109, 124, 27, 91, 48, 109, 32, 27, 91, 48, 59, 51, 52, 109, 124, 27, 91, 48, 109, 32, 27, 91,
			48, 59, 51, 52, 109, 124, 95, 41, 27, 91, 48, 109, 32, 27, 91, 48, 59, 51, 52, 109, 124, 27, 91, 48, 109,
			32, 27, 91, 48, 59, 51, 52, 109, 124, 27, 91, 48, 109, 32, 27, 91, 48, 59, 51, 55, 109, 124, 27, 91, 48,
			109, 32, 27, 91, 48, 59, 51, 55, 109, 124, 95, 41, 27, 91, 48, 109, 32, 27, 91, 48, 59, 51, 55, 109, 124,
			27, 91, 48, 109, 32, 27, 91, 48, 59, 51, 55, 109, 124, 27, 91, 48, 109, 32, 32, 32, 32, 10, 27, 91, 48, 59,
			51, 52, 109, 124, 27, 91, 48, 109, 32, 27, 91, 48, 59, 51, 52, 109, 124, 95, 124, 27, 91, 48, 109, 32, 27,
			91, 48, 59, 51, 52, 109, 124, 27, 91, 48, 109, 32, 32, 27, 91, 48, 59, 51, 55, 109, 95, 95, 47, 124, 27, 91,
			48, 109, 32, 27, 91, 48, 59, 51, 55, 109, 124, 27, 91, 48, 109, 32, 27, 91, 48, 59, 51, 55, 109, 124, 27,
			91, 48, 109, 32, 32, 27, 91, 48, 59, 51, 55, 109, 95, 95, 47, 124, 27, 91, 48, 109, 32, 27, 91, 48, 59, 49,
			59, 51, 48, 59, 57, 48, 109, 124, 95, 95, 95, 27, 91, 48, 109, 32, 10, 32, 27, 91, 48, 59, 51, 55, 109, 92,
			95, 95, 95, 47, 124, 95, 124, 27, 91, 48, 109, 32, 32, 32, 27, 91, 48, 59, 51, 55, 109, 124, 95, 124, 27,
			91, 48, 109, 32, 27, 91, 48, 59, 49, 59, 51, 48, 59, 57, 48, 109, 124, 95, 124, 27, 91, 48, 109, 32, 32, 32,
			32, 27, 91, 48, 59, 49, 59, 51, 48, 59, 57, 48, 109, 92, 95, 95, 95, 95, 124, 27, 91, 48, 109, 10, 32, 32,
			32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32,
			10, 87, 101, 108, 99, 111, 109, 101, 32, 116, 111, 32, 27, 91, 48, 59, 57, 49, 109, 65, 114, 109, 98, 105,
			97, 110, 32, 50, 51, 46, 48, 50, 46, 50, 32, 66, 117, 108, 108, 115, 101, 121, 101, 27, 91, 48, 109, 32,
			119, 105, 116, 104, 32, 27, 91, 48, 59, 57, 49, 109, 76, 105, 110, 117, 120, 32, 54, 46, 49, 46, 54, 51, 45,
			99, 117, 114, 114, 101, 110, 116, 45, 115, 117, 110, 120, 105, 27, 91, 48, 109, 10, 10, 27, 91, 48, 59, 57,
			49, 109, 78, 111, 32, 101, 110, 100, 45, 117, 115, 101, 114, 32, 115, 117, 112, 112, 111, 114, 116, 58, 32,
			27, 91, 48, 109, 117, 110, 116, 101, 115, 116, 101, 100, 32, 97, 117, 116, 111, 109, 97, 116, 101, 100, 32,
			98, 117, 105, 108, 100, 10, 10, 83, 121, 115, 116, 101, 109, 32, 108, 111, 97, 100, 58, 32, 32, 27, 91, 48,
			59, 57, 50, 109, 32, 50, 37, 27, 91, 48, 109, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 9, 85, 112, 32,
			116, 105, 109, 101, 58, 32, 32, 32, 32, 32, 32, 32, 27, 91, 57, 50, 109, 50, 32, 100, 97, 121, 115, 32, 49,
			50, 58, 52, 52, 27, 91, 48, 109, 9, 10, 77, 101, 109, 111, 114, 121, 32, 117, 115, 97, 103, 101, 58, 32, 27,
			91, 48, 59, 57, 50, 109, 32, 56, 37, 27, 91, 48, 109, 32, 111, 102, 32, 57, 57, 55, 77, 32, 32, 32, 9, 73,
			80, 58, 9, 32, 32, 32, 32, 32, 32, 32, 27, 91, 57, 50, 109, 49, 50, 55, 46, 48, 46, 48, 46, 49,
			27, 91, 48, 109, 10, 67, 80, 85, 32, 116, 101, 109, 112, 58, 32, 32, 32, 32, 32, 27, 91, 48, 59, 57, 50,
			109, 32, 51, 55, 194, 176, 67, 27, 91, 48, 109, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 9, 85, 115, 97,
			103, 101, 32, 111, 102, 32, 47, 58, 32, 32, 32, 27, 91, 48, 59, 57, 50, 109, 32, 49, 53, 37, 27, 91, 48,
			109, 32, 111, 102, 32, 50, 57, 71, 32, 32, 32, 32, 9, 10, 82, 88, 32, 116, 111, 100, 97, 121, 58, 32, 32,
			32, 32, 32, 32, 27, 91, 57, 50, 109, 49, 48, 46, 49, 32, 77, 105, 66, 27, 91, 48, 109, 32, 32, 9, 10, 10,
			91, 27, 91, 51, 49, 109, 32, 71, 101, 110, 101, 114, 97, 108, 32, 115, 121, 115, 116, 101, 109, 32, 99, 111,
			110, 102, 105, 103, 117, 114, 97, 116, 105, 111, 110, 32, 40, 98, 101, 116, 97, 41, 27, 91, 48, 109, 58, 32,
			27, 91, 49, 109, 97, 114, 109, 98, 105, 97, 110, 45, 99, 111, 110, 102, 105, 103, 27, 91, 48, 109, 32, 93,
			10, 10,
		},
	},
}

func BenchmarkParseColourCodes(b *testing.B) {
	for _, bench := range benches {
		b.Run(
			bench.name, func(b *testing.B) {
				b.ReportAllocs()

				for i := 0; i < b.N; i++ {
					var data = bench.test
					for len(data) > 0 {
						var advance, token, _ = ScanCodes(data, false)
						data = data[advance:]

						if advance == 0 {
							panic("Avoid infinite loop")
						}

						if IsSequence(token) {
							_ = MustParseColourCodes(token)
						}
					}
				}
			},
		)
	}
}

func BenchmarkParseColourCodesWithScanner(b *testing.B) {
	for _, bench := range benches {
		b.Run(
			bench.name, func(b *testing.B) {
				b.ReportAllocs()

				for i := 0; i < b.N; i++ {
					var (
						data    = bench.test
						reader  = bytes.NewReader(data)
						scanner = bufio.NewScanner(reader)
					)
					scanner.Split(ScanCodes)

					for scanner.Scan() {
						var token = scanner.Bytes()
						if IsSequence(token) {
							_ = MustParseColourCodes(token)
						}
					}
				}
			},
		)
	}
}

func ExampleScanCodes() {
	// Red Hello, World!
	var modifier = []byte{EscapeCode, StartCode, byte('3'), byte('1'), EndCode}
	var sequence = []byte(string(modifier) + "Hello, World!")

	// Create scanner
	var reader = bytes.NewReader(sequence)
	var scanner = bufio.NewScanner(reader)
	// Set ansi.ScanCodes as the split func
	scanner.Split(ScanCodes)

	// Enter scanner loop
	for scanner.Scan() {
		// Grab token
		var token = scanner.Bytes()
		// Check if it's an ANSI sequence
		if IsSequence(token) {
			// Parse the colours from the sequence
			var colours = MustParseColourCodes(token)
			// Do something (probably style some terminal)
			fmt.Println(colours)
		} else {
			// Normal text, so just print it
			fmt.Println(string(token))
		}
	}
	// Output:
	// [RED]
	// Hello, World!
}

func ExampleIsSequence() {
	var sequence = []byte{EscapeCode, StartCode, byte('0'), EndCode}
	fmt.Println(IsSequence(sequence))
	// Output:
	// true
}
