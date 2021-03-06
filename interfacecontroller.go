// interfacecontroller.go contains functions to manage the execution of programs interacting with the OS
// it handles the conditions under which certain calls to functions within executer.go are made and returns
// messages avaiable for use in view.go

package main

import "strings"

// messageRetrieve ... Pass the users selection as an inode for fls or command+inode for carving and a slice of the current
// directory to retrieve a single element of the slice.
func messageRetrieve(selection string, input []string) (results string, log string) {

	//Pull the inode and name from the selection.
	// TODO Test this on selection without a slice to determine how robust the regex is.
	selectionInode := inodeMatcher(selection)
	selectionName := nameMatcher(selection)
	selectedType := dirMatcher(selection)

	// Does not matter that this is not robust logic, input will be locked down by line selection in tcell version.
	// References to global executionArgs and the unallocatedRecover command line parameter
	// Return a success message if executeCarvers completes successfully, else return a failure message to the user.
	if strings.Contains(selectedType, "d") {
		if executeCarvers("tsk_recover", argsUpdaterRecover(executionArgsRecover, selectionInode, *unallocatedRecover), selectionName) {
			return "", "Successfully carved " + selectionName
		}
		return "", "Failed to carve " + selectionName
	}

	if strings.Contains(selectedType, "r") {
		if executeCarvers("icat", argsUpdater(executionArgs, selectionInode), selectionName) {
			return "", "Successfully carved " + selectionName
		}
		return "", "Failed to carve " + selectionName
	}
	// If there's nothing in the selection string providing direction to carve,
	// then iterate through input and return a line that has an inode matching to the
	// selectionInode.
	for _, inputLine := range input {
		if selectionInode == inodeMatcher(inputLine) {
			if !executeFLS(executionArgs) {
				return inputLine, ""
			}
			return "", "FLS failed to execute on directory: " + selectionName
		}
	}

	return "", "Unable to find a line in the FLS input slice that matches the selected input string."
}
