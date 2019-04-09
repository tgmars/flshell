package main

import (
	"fmt"
	"os"
	"os/exec"
)

//Reading executer.go
// This file flows in building blocks, initially constructing the arguments that are passed.
// Then the struct that is used to execute a command
// The code that enables execution of the given struct

// argsUpdater ... Updates the slice of arguments that are passed with the command to be executed.
// if there are less than 4 arguments (ie, a path to an image, -o, and an offset (with no inode))
// Then the arguments are updated and the inode string is appended. Otherwise, just update the inode
// value to the requested one.
func argsUpdater(arguments []string, inode string) []string {
	if len(arguments) < 4 {
		arguments = append(arguments, inode)
	}
	arguments[3] = inode
	return arguments
}

// argsUpdaterRecover ... Updates the slice of arguments that are passed with the command to be executed.
// if there are less than 5 arguments (ie, a path to an image, -o, an offset, -d and no inode)
// Then the arguments are updated and the inode string is appended. Otherwise, just update the inode
// value to the requested one.
func argsUpdaterRecover(arguments []string, inode string, carveunallocated bool) []string {
	if len(arguments) < 6 {
		arguments = append(arguments, inode, "a")
	}
	if carveunallocated {
		arguments[5] = "e"
	} else {
		arguments[5] = "a"
	}
	arguments[4] = inode
	return arguments
}

// executeFLS ... Pass the command to be executed and its arguments as a slice of strings.
// fills an item object with the current items for fls, for icat and tsk_recover it writes to a file.
func executeFLS(args []string) bool {
	cmdStruct := exec.Command("fls", args...)
	cmdOutput := commandExecuter(cmdStruct)

	// Only attempt to populate the directory if the output of the fls command run contains stuff.
	if len(cmdOutput) > 0 {
		directory.populate(cmdOutput)
		return true
	}

	return false
}

// executeCarvers ... Pass the command to be executed and its arguments as a slice of strings.
// Writes a directory or file and its mft information to a file specified by itemname.
// Returns true on successful completion, false if no cases are met.
func executeCarvers(cmd string, args []string, itemname string) bool {
	switch cmd {
	case "icat":
		icatstruct := exec.Command(cmd, args...)
		writeCmdToFile(itemname, icatstruct)
		istatstruct := exec.Command(cmd, args...)
		writeCmdToFile(itemname+".mft", istatstruct)
		return true
	case "tsk_recover":
		recoverstruct := exec.Command(cmd, args...)
		writeCmdToFile(itemname, recoverstruct)
		return true
	}
	return false
}

// commandExecuter ... Executes a command on the host given a command structure and
// returns the output as a string.
func commandExecuter(cmdstruct *exec.Cmd) string {
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
