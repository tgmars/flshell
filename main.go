package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nsf/termbox-go"
)

var current string
var fullcurrent string
var cmd = "fls"
var selectedline = 0
var selectedrunes []rune
var selectedstring string
var cache []string

var imagepath = os.Args[1]
var diskoffset = os.Args[2]

var windowheight int

//var imagepath = "/media/vboxshared/recruitment.raw"
//var diskoffset = "718848"

var args = []string{"-o", diskoffset, imagepath}
var dirlevel int
var usecache bool
var maxlines int
var goingup bool

var directory = Item{"d/d", "1", "root", nil, nil}

//Print header func

//Print footer func

//Print, if statement is to catch the fls output and present it as it does in tool output.
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {

	currentline := 0
	for _, c := range msg {
		if c == '\n' {
			y++
			x = -1
			currentline++
		}
		if currentline == selectedline {
			termbox.SetCell(x, y, c, fg, termbox.ColorGreen)
			selectedrunes = append(selectedrunes, c)
		} else {
			termbox.SetCell(x, y, c, fg, bg)
		}
		x++
	}
	if selectedrunes != nil {
		selectedstring = string(selectedrunes)
	}
	selectedrunes = nil
}

//Move the select line var by an amount of lines
// TODO all top end error catching (currently only catches negatives)
func moveSelectedLine(amount int, maxlines int) {
	if (amount < 0) && (selectedline == 0) {
		selectedline = 0
	} else if (amount > 0) && (selectedline == maxlines-1) {
		selectedline = maxlines - 1
	} else {
		selectedline += amount
	}
}

func redrawAll() {
	const defaultcolour = termbox.ColorDefault
	termbox.Clear(defaultcolour, defaultcolour)
	tbprint(0, 0, defaultcolour, defaultcolour, current)

	termbox.Flush()
}

// Executes a command on the host and prints the output as a string.
func executer(cmdstruct *exec.Cmd) string {

	// Define vars that will be used to store output and error of running the command.
	var (
		cmdOutput []byte
		cmdErr    error
	)
	if cmdOutput, cmdErr = cmdstruct.Output(); cmdErr != nil {
		fmt.Fprintln(os.Stderr, cmdErr)
		os.Exit(1)
	}
	return string(cmdOutput)
}

// Updates the struct that is passed to exec.Output() to include the current directory inode.
func argsupdater(arguments []string, inode string) []string {
	if len(arguments) < 4 {
		arguments = append(arguments, inode)
		// fmt.Println("args were less than 4. Appending inode value. ", args)
	}

	arguments[3] = inode
	// fmt.Println("Updated args[3] with inode of:", inode)
	//Yuck, fix this.
	return arguments

}

// execute new command
func commandexecuter() {
	cmdstruct := exec.Command(cmd, args...)
	directory.populate(executer(cmdstruct))
	displayexecuter()
}

//Alternative to commandexecuter() that hsould be called when a command is not required to be run.
func displayexecuter() {
	fullcurrent = directory.listChildren()
	maxlines = newlineCounter(fullcurrent)
	current = windowString(windowheight, fullcurrent)

	//current = directory.listChildren()
	//maxlines = newlineCounter(current)
}

func windowString(height int, message string) string {
	fmt.Printf("\t\t\tHeight: %v\tMaxlines: %v", height, maxlines)
	//take message and return a number of lines that equals window height.
	//return a specific x line section of those lines, provided an offset of lines into the string

	lines := strings.Split(message, "\n")
	if height >= len(lines) {
		fmt.Println("\t\t\t\t len(lines)=" + string(len(lines)))
		return message
	} else {
		writeStringToFile("lines.txt", strings.Join(lines, "\n"))
		return strings.Join(lines[0:height], "\n")
	}
}

func icatexecuter() {
	filename := nameMatcher(selectedstring)
	// execute new command
	cmdstruct := exec.Command(cmd, args...)
	// open the out file for writing
	writeCmdToFile(filename, cmdstruct)
	fmt.Print("\t\t Wrote " + filename)
}

func istatexecuter() {
	filename := nameMatcher(selectedstring) + ".mft"
	// execute new command
	cmdstruct := exec.Command(cmd, args...)
	// open the out file for writing
	writeCmdToFile(filename, cmdstruct)
	fmt.Print("\t\t Wrote " + filename + "\n")
}

func main() {

	// Define a cmd struct that consists of the executable, it's location and the arguments passed.

	//Initialise the termbox.
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	//Set input to standard inputescape' mode.
	termbox.SetInputMode(termbox.InputEsc)
	firstrun := true
	_, windowheight = termbox.Size()
	redrawAll()

mainloop:
	for {
		if cmd == "fls" && !directory.hasChildren() {
			commandexecuter()
		}
		if firstrun {
			redrawAll()
		}

		switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey:
			switch event.Key {

			case termbox.KeyArrowUp:
				moveSelectedLine(-1, maxlines)

				cmd = "fls"

			case termbox.KeyArrowDown: // on Arrow Down
				moveSelectedLine(1, maxlines)
				cmd = "fls"

			case termbox.KeyArrowLeft: //on Arrow Left
				cmd = "fls" //cache this
				if dirlevel < 1 {
					dirlevel = 0
				} else {
					dirlevel--
					selectedline = 0
					directory = *directory.goUp(directory)
					displayexecuter()
				}

			case termbox.KeyArrowRight:
				cmd = "fls"
				ftype := dirMatcher(selectedstring)
				inode := inodeMatcher(selectedstring)
				//update the command struct with the currently selected string
				if ftype == "d/d" {
					selectedline = 0
					dirlevel++
					directory = *directory.goDown(directory, ftype, inode)
					args = argsupdater(args, inodeMatcher(selectedstring))
					commandexecuter()

				}
			case termbox.KeyEnter:
				cmd = "fls"
				ftype := dirMatcher(selectedstring)
				inode := inodeMatcher(selectedstring)
				//update the command struct with the currently selected string
				if ftype == "d/d" {
					selectedline = 0
					dirlevel++
					directory = *directory.goDown(directory, ftype, inode)
					args = argsupdater(args, inodeMatcher(selectedstring))
					commandexecuter()
				}

			case termbox.KeyTab:
				//add a second event to confirm if you wish to output files
				if dirMatcher(selectedstring) == "d/d" {
					cmd = "tsk_recover"
					args = argsupdater(args, inodeMatcher(selectedstring))
					commandexecuter()
				}
				if dirMatcher(selectedstring) == "r/r" {
					cmd = "icat" //make this icat
					args = argsupdater(args, inodeMatcher(selectedstring))
					icatexecuter()
					cmd = "istat" //make this icat
					args = argsupdater(args, inodeMatcher(selectedstring))
					istatexecuter()
				}
			case termbox.KeyCtrlC:
				break mainloop
			}

		case termbox.EventError:
			panic(event.Err)
		}

		/*
			_
			if y <= maxlines {
			} else {
				current = fullcurrent
			}*/

		firstrun = false
		redrawAll()
		_, windowheight = termbox.Size()

		//fmt.Printf("current dir\t%+v\n", &currentDir)

	}

	// Execute fls initially to get inodes of the root directory.

}
