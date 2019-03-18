package main

import "strings"

func windowString(height int, message string, selected int) string {
	//fmt.Printf("\t\t\tHeight: %v\tMaxlines: %v", height, maxlines)
	//take message and return a number of lines that equals window height.
	//return a specific x line section of those lines, provided an offset of lines into the string

	//height is a value that starts at 1, ie 1 line minimum returns one.
	// selected starts at 0, need to compensate for this.
	selected++

	lines := strings.Split(message, "\n")
	if height >= len(lines) {
		return message
	}
	if selected > height {
		scrollline = height
		return strings.Join(lines[selected-height:selected], "\n")
	} else {
		//	writeStringToFile("lines.txt", strings.Join(lines, "\n"))
		return strings.Join(lines[0:height], "\n")
	}

}
