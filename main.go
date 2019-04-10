package main

import (
	"bufio"
	"flag"
	"log"
	"strings"
)

//Legacy vars
var current string
var currentSplit []string
var fullcurrent string
var cmd = "fls"

//Current, in use vars
var directory = Item{"d/d", "1", "/", nil, nil}
var imagepath = flag.String("f", "", "Path to the file to analyse. Usage example: -f=/path/to/image.dd")
var diskoffset = flag.String("o", "0", "Offset of the desired partition (in sectors). Usage example: -o=2048 ")
var unallocatedRecover = flag.Bool("e", false, "When recovering a directory, carve unallocated entries too. Disabled by default. Usage to enable: -e")

// Two slices of type string used to hold the arguments for fls/icat/istat execution and the other for tsk_recover.
var executionArgs []string
var executionArgsRecover []string

func main() {
	setupMain()
	s := newView()
	loop(s)
}

// setupMain ... Initialise and setup variables before entering the main loop of program execution.
// Assign command line parameters to variables and conduct an initial run of FLS in the root directory.
// TODO: integrate initial FLS execution with interfacecontroller.go
func setupMain() {
	flag.Parse()
	executionArgs = []string{"-o", *diskoffset, *imagepath, "5"}
	executionArgsRecover = []string{"-o", *diskoffset, *imagepath, "-d"}
	if !executeFLS(executionArgs) { //operates on cache object directory
		log.Println("failed to execute FLS")
	}
	prepGlobals()
}

// lineSelector ... Given an input reader, parse a line of input, remove the newline from it and return just the value.
func lineSelector(r *bufio.Reader) string {
	selection, _ := r.ReadString('\n')
	return strings.Replace(selection, "\n", "", -1)
}

// prepGlobals ... prepare global variables for future use to speed up processing.
func prepGlobals() {
	current = directory.listChildren()          //global current text
	currentSplit = strings.Split(current, "\n") //global slice of current text.
}
