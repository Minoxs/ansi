# ansi

An ANSI color sequence parser for Go.
Compatible with color.Color and bufio.SplitFunc.

The color palette used for the RGBA() functions is Windows 10's Campbell, the default terminal palette, because I find them pleasing.
If for some reason that is a problem, either fork the repo and change it or create an issue and I can make it customizable in some way.

RGB values were taken from [Wikipedia](https://en.wikipedia.org/wiki/ANSI_escape_code)
and checked in [Microsoft's official documentation](https://learn.microsoft.com/en-us/windows/terminal/customize-settings/color-schemes).
