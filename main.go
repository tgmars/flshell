package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/nsf/termbox-go"
)

var current string
var cmd = "fls"
var selectedline = 0
var selectedrunes []rune
var selectedstring string
var cached [10]string

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

func regexMatcher(line string) string {
	re := regexp.MustCompile(`(?m)[rd]\/[rd]\s(\d+-\d+-\d+):`)
	return re.FindStringSubmatch(line)[1]
}

//Move the select line var by an amount of lines
// TODO all top end error catching (currently only catches negatives)
func moveSelectedLine(amount int) {
	if (amount < 0) && (selectedline == 0) {
		selectedline = 0
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

// Updates the struct that is passed to exec.Output() to include the current directory inode.
func argsupdater(args []string, inode string) []string {
	if len(args) < 4 {
		args = append(args, inode)
		// fmt.Println("args were less than 4. Appending inode value. ", args)
	}
	args[3] = inode
	// fmt.Println("Updated args[3] with inode of:", inode)
	//Yuck, fix this.
	return args

}

var imagepath = os.Args[1]
var diskoffset = os.Args[2]
var args = []string{"-o", diskoffset, imagepath}
var cachecounter int = 0
var usecache bool

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
	goingup := false

mainloop:
	for {
		if (len(cached[0]) == 0) || !(usecache) {
			cmdstruct := exec.Command(cmd, args...)
			current = executer(cmdstruct)
			cached[cachecounter] = current
			cachecounter++
		} else if usecache && goingup {
			current = cached[cachecounter]
		}

		switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey:
			switch event.Key {
			case termbox.KeyArrowUp:
				moveSelectedLine(-1)
				cmd = "fls"
				usecache = true
				goingup = false
			case termbox.KeyArrowDown:
				moveSelectedLine(1)
				cmd = "fls"
				usecache = true
				goingup = false
			case termbox.KeyArrowLeft:
				selectedline = 0
				cmd = "fls" //cache this
				usecache = true
				goingup = true
				cachecounter--
			case termbox.KeyArrowRight:
				selectedline = 0
				cmd = "fls"
				usecache = false
				goingup = false
				//update the command struct with the currently selected string
				args = argsupdater(args, regexMatcher(selectedstring))
			case termbox.KeyEnter:
				current = "Enter"
				cmd = "fls" //make this icat
			case termbox.KeyCtrlC:
				break mainloop
			}

		case termbox.EventError:
			panic(event.Err)
		}
		redrawAll()

	}

	// Execute fls initially to get inodes of the root directory.

}
