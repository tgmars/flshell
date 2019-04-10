package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

var windowHeight int
var windowWidth int
var style = tcell.StyleDefault
var selectedStyle = tcell.StyleDefault.Background(tcell.ColorGreen)
var footerStyle = tcell.StyleDefault.Reverse(true)

// newView ... Create a new view with ASCII encoding as fallback, default style
// initialise windowHeight and windowWidth variables.
func newView() tcell.Screen {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	s.Clear()
	windowWidth, windowHeight = s.Size()
	return s
}

func populateView(s tcell.Screen) {
	s.Clear()
	drawHeader(s)
	drawBody(s, currentSplit, sl)
	drawFooter(s, "Message: "+"TOOL OUTPUT MESSAGES HERE")
	s.Show()
}

// drawHeader ... Prints current directory, the directory name, to the top line of
// the terminal window using a call to putln() with the footerStyle
func drawHeader(s tcell.Screen) {
	content := "Current directory: " + directory.Name
	colourRow(s, footerStyle, 0)
	putln(s, footerStyle, content, 0)
}

// drawFooter ... Prints version number, input image file and supplied content to the bottom line of
// the terminal window using a call to putln() with the footerStyle
func drawFooter(s tcell.Screen, content string) {
	content = "FLShell v2.0 | Image File: " + *imagepath + " | " + content
	colourRow(s, footerStyle, windowHeight-1)
	putln(s, footerStyle, content, windowHeight-1)
}

// drawBody ... Draws body of FLS output
// Draws line by line using default style unless the line index is the 'selectedline',
// in which case the selectedStyle is applied to the line.
// selectedline must be greater than 0 and less than windowHeight-1.
func drawBody(s tcell.Screen, content []string, selectedline int) {
	for i, line := range content {
		// currentline = i+1 to account for header line taking index 0 on the screen.
		currentline := i + 1
		if selectedline == currentline {
			colourRow(s, selectedStyle, currentline)
			putln(s, selectedStyle, line, currentline)
		} else {
			putln(s, style, line, currentline)
		}

	}
}

// colourRow ... Set a row on the screen with the specified style, ideally with a specific background colour.
func colourRow(s tcell.Screen, style tcell.Style, row int) {
	for x := 0; x < windowWidth; x++ {
		s.SetCell(x, row, style, ' ')
	}
}

// putln ... Stolen straight from https://github.com/gdamore/tcell/blob/master/_demos/unicode.go
// provide a screen and the string to print and it will write on the current 'row' (if implemented).
func putln(s tcell.Screen, style tcell.Style, str string, row int) {
	puts(s, style, 0, row, str)
	row++
}

// puts ... Stolen straight from https://github.com/gdamore/tcell/blob/master/_demos/unicode.go
// Writes to a screen at a given position provided a screen and position.
// Is unicode tolerant
func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	zwj := false
	for _, r := range str {
		if r == '\u200d' {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
			deferred = append(deferred, r)
			zwj = true
			continue
		}
		if zwj {
			deferred = append(deferred, r)
			zwj = false
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}
