package main

import (
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

var current string
var fullcurrent string
var cmd = "fls"
var selectedline = 0
var selectedrunes []rune
var selectedstring string
var cache []string
var forceprint string

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
var enteredbaddir bool

var directory = Item{"d/d", "1", "root", nil, nil}

//Print header func

//Print footer func

//Print, if statement is to catch the fls output and present it as it does in tool output.
func tbprint(x, y int, fg, bg termbox.Attribute, msg string, sl int, height int) {
	//offset used to record the offset that should be applied to a selectedline so that it references a line in the printable range.
	if sl+1 > height {
		offset := sl - height
		sl = sl - offset - 1
	}
	currentline := 0
	if forceprint != "" {
		msg = forceprint
	}
	for _, c := range msg {
		if c == '\n' {
			y++
			x = -1
			currentline++
		}
		//this logic won't work out
		//fmt.Printf("\t\t\t currentline: %v \t\t\t  selectedline: %v\n", currentline, sl)
		if currentline == sl {
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
	tbprint(0, 0, defaultcolour, defaultcolour, current, selectedline, windowheight)

	termbox.Flush()
}

func goUp() {
	if dirlevel < 1 {
		dirlevel = 0
	} else {
		dirlevel--
		selectedline = 0
		directory = *directory.goUp(directory)
		displayexecuter()
	}

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
		if enteredbaddir {
			time.Sleep(time.Second * 3)
			break mainloop
		}

		current = windowString(windowheight, fullcurrent, selectedline+1)
		if firstrun {
			redrawAll()
		}

		switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey:
			switch event.Key {

			case termbox.KeyArrowUp:
				moveSelectedLine(-1, maxlines)
				cmd = "fls"
				current = windowString(windowheight, fullcurrent, selectedline+1)

			case termbox.KeyArrowDown: // on Arrow Down
				moveSelectedLine(1, maxlines)
				cmd = "fls"
				current = windowString(windowheight, fullcurrent, selectedline+1)

				//displayexecuter() //move these out eventually, unnessarily slow.

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
					//args = argsupdater(args, inodeMatcher(selectedstring))
					//commandexecuter()
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

		firstrun = false
		redrawAll()
		_, windowheight = termbox.Size()
	}
}
