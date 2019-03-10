package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/nsf/termbox-go"
)

var current string
var cmd = "fls"
var selectedline = 0
var selectedrunes []rune
var selectedstring string
var cache []string

var imagepath = os.Args[1]
var diskoffset = os.Args[2]

//var imagepath = "/media/vboxshared/recruitment.raw"
//var diskoffset = "718848"

var args = []string{"-o", diskoffset, imagepath}
var cachecounter int = 0
var usecache bool
var maxlines int
var goingup = false

var currentDir = Folder{Item{"d/d", "1", "root", nil}, nil}

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

func checkReErr(e error) {
	if e != nil {
		writeStringToFile("RegexErrors.txt", e.Error())
	}
}

func inodeMatcher(line string) string {
	re, err := regexp.Compile(`(?m).\/.\s(\*?\s*.*):\s.*`)
	checkReErr(err)
	return re.FindStringSubmatch(line)[1]
}

func dirMatcher(line string) string {
	re, err := regexp.Compile(`(?m)(.\/.)\s\*?\s*.*:\s.*`)
	checkReErr(err)
	return re.FindStringSubmatch(line)[1]
}

func nameMatcher(line string) string {
	re, err := regexp.Compile(`(?m):\s+(.+)`)
	checkReErr(err)
	return re.FindStringSubmatch(line)[1]
}

func newlineCounter(input string) int {
	re, err := regexp.Compile(`(?m)\n`)
	checkReErr(err)
	return len(re.FindAllStringIndex(input, -1))
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

	// .Run() returns an error if anything fails when opening pipes to STDIN, STDOUT or STDERR.
	// First we execute cmdstruct.Run() and store the results in err. Then we evaluate the if statement to determine if an error occured or not.
	// Therefore if error resolves to true then .Run() has executed with an error and we need to handle that.
	// Has been replaced with .Output() to return the result from STDOUT too.
	if cmdOutput, cmdErr = cmdstruct.Output(); cmdErr != nil {
		fmt.Fprintln(os.Stderr, cmdErr)
		os.Exit(1)
	}
	return string(cmdOutput)
}

func fillCache(input string, dir Folder) {
	//Split input on newlines and assign to a slice of strings
	s := strings.Split(input, "\n")
	//For each line in the split, add it as a child to the root folder.
	for _, line := range s {
		// Only attempt to parse values out of the line if there's stuff in it.
		if line != "" {
			itemType := dirMatcher(line)
			itemInode := inodeMatcher(line)
			itemName := nameMatcher(line)
			dir.addChild(&Item{itemType, itemInode, itemName, &dir})
		}
	}
	writeStringToFile("children.txt", dir.listChildren())
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

func commandexecuter() {
	// execute new command
	cmdstruct := exec.Command(cmd, args...)
	raw := executer(cmdstruct)
	writeStringToFile("raw.txt", raw)
	fillCache(raw, currentDir)
	currentDir.sortChildrenByAlphaDescending()
	current = currentDir.listChildren()
	writeStringToFile("current.txt", current)
	maxlines = newlineCounter(current)

	/*if (len(cache) == 0) || !(usecache) {
		cmdstruct := exec.Command(cmd, args...)
		current = executer(cmdstruct)
		maxlines = newlineCounter(current)
		// use cached command
	} else if usecache && goingup {
		current = cache[cachecounter]
	}
	*/
}

func icatexecuter() {
	filename := nameMatcher(selectedstring)
	// execute new command
	cmdstruct := exec.Command(cmd, args...)
	// open the out file for writing
	writeCmdToFile(filename, cmdstruct)
	fmt.Println("\tSucessfully wrote " + filename)
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

mainloop:
	for {
		if cmd == "fls" {
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
				usecache = true
				goingup = false
			case termbox.KeyArrowDown:
				moveSelectedLine(1, maxlines)
				cmd = "fls"
				usecache = true
				goingup = false
			case termbox.KeyArrowLeft:
				cmd = "fls" //cache this
				usecache = true
				goingup = true
				if cachecounter < 1 {
					cachecounter = 0
				} else {
					cachecounter--
					selectedline = 0
					current = cache[cachecounter]

				}
			case termbox.KeyArrowRight:
				cmd = "fls"
				usecache = false
				goingup = false
				//update the command struct with the currently selected string
				if dirMatcher(selectedstring) == "d/d" {
					selectedline = 0
					args = argsupdater(args, inodeMatcher(selectedstring))
					commandexecuter()
					cache = append(cache, current)
					cachecounter++
				}
			case termbox.KeyEnter:
				usecache = false
				goingup = false
				if dirMatcher(selectedstring) == "d/d" {
					cmd = "tsk_recover"
					args = argsupdater(args, inodeMatcher(selectedstring))
					commandexecuter()
				}
				if dirMatcher(selectedstring) == "r/r" {
					cmd = "icat" //make this icat
					args = argsupdater(args, inodeMatcher(selectedstring))
					icatexecuter()
				}
			case termbox.KeyCtrlC:
				break mainloop
			}

		case termbox.EventError:
			panic(event.Err)
		}
		firstrun = false
		redrawAll()

	}

	// Execute fls initially to get inodes of the root directory.

}
