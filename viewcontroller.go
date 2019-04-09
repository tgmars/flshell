package main

import (
	"time"

	"github.com/gdamore/tcell"
)

// globally defined selected line starts at 1 to account for header width.
var sl = 1

// loop ... main execution loop. Figure out wtf is happening here.
func loop(s tcell.Screen) {
	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyCtrlC:
					close(quit)
					return
				case tcell.KeyDown:
					moveSelectedLine(1, sl, 1, 1)

				case tcell.KeyUp:
					moveSelectedLine(-1, sl, 1, 1)

				case tcell.KeyLeft:

					s.Sync()
				case tcell.KeyRight:
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
				windowWidth, windowHeight = s.Size()
			}
		}
	}()
loop:
	// get good at go and learn channels, routines and selects properly.
	for {
		select {
		case <-quit:
			break loop
		// what do you do
		case <-time.After(time.Millisecond * 5):
		}
		populateView(s)
	}
	s.Fini()
}

// moveSelectedLine ... Return the currently selected line with deviation specified by amount value.
// Amount should only ever be 1 or -1. The logic in moveSelectedLine is not robust enough
// to safely handle other values.
func moveSelectedLine(amount int, selectedline int, headerheight int, footerheight int) {
	// if amount is a negative value and the currently selected line is the top (ie, headerwidth), keep it at the
	// at the index that's the same thickness as the header.
	if (amount < 0) && (selectedline == headerheight) {
		sl = headerheight
		// if the amount is positive and the currently selected line is at the max of the window, including the footer,
		// then return the max value.
	} else if (amount > 0) && (selectedline == windowHeight-footerheight-1) {
		sl = windowHeight - footerheight - 1
	} else {
		sl = selectedline + amount
	}

}

// windowString ... takes a chunk of text delimited by newlines and returns a specified portion of it
// that fits within the current size of the terminal window. A number of lines to offset
// into the original message must be provided by the selected integer.
// windowheight will usually be an integer that starts at 1 whereas selected is likely
// starting from 0. When passing variables to this function, it is recommended to increment
// selected by 1.
// The function does not return a reduced portion of the input message if it has less lines
// than windowheight.
/*func windowString(lines []string, selectedline int, headerheight int, footerheight int) []string {
	if windowHeight-headerheight-footerheight >= len(lines) {
		return lines
	}
	if selected > windowheight {
		return strings.Join(lines[selected-windowheight:selected], "\n")
	}
	return strings.Join(lines[0:windowheight], "\n")
}
*/
