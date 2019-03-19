package main

import (
	"fmt"
	"os"
	"os/exec"
)

//defines the information to be passed to the main view
//handles underlying program execution

// execute new command
func commandexecuter() {
	cmdstruct := exec.Command(cmd, args...)
	tooloutput := executer(cmdstruct)
	directory.populate(tooloutput)
	displayexecuter()
}

//Alternative to commandexecuter() that hsould be called when a command is not required to be run.
func displayexecuter() {
	fullcurrent = directory.listChildren()
	maxlines = newlineCounter(fullcurrent)
	current = windowString(windowheight, fullcurrent, selectedline+1)

	//current = directory.listChildren()
	//maxlines = newlineCounter(current)
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
	output := string(cmdOutput)
	if output == "" && cmd == "fls" {
		forceprint = "FLS did not output anything, try another method to investigate this directory. FLShell will quit in 3 seconds."
		enteredbaddir = true
	} else {
		forceprint = ""
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
