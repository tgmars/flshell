package main

import "strings"

// windowString ... takes a chunk of text delimited by newlines and returns a specified portion of it
// that fits within the current size of the terminal window. A number of lines to offset
// into the original message must be provided by the selected integer.
// windowheight will usually be an integer that starts at 1 whereas selected is likely
// starting from 0. When passing variables to this function, it is recommended to increment
// selected by 1.
// The function does not return a reduced portion of the input message if it has less lines
// than windowheight.
func windowString(windowheight int, message string, selected int) string {
	lines := strings.Split(message, "\n")
	if windowheight >= len(lines) {
		return message
	}
	if selected > windowheight {
		return strings.Join(lines[selected-windowheight:selected], "\n")
	}
	return strings.Join(lines[0:windowheight], "\n")
}
