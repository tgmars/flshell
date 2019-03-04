package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

const (
	// Command to be passed to os to execute.
	cmd = "fls"
)

func main() {

	imagepath := os.Args[1]
	diskoffset := os.Args[2]

	args := []string{"-o", diskoffset, imagepath}

	// Define a cmd struct that consists of the executable, it's location and the arguments passed.
	cmdstruct := exec.Command(cmd, args...)

	// Execute fls initially to get inodes of the root directory.

	fmt.Println(executer(cmdstruct))
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter inode: ")
	for scanner.Scan() {
		//fmt.Println("Top of scan loop.")

		//Reading using Scanner is best for simple use cases.
		//Might have to move to ReadString for keys pressed.
		inode := scanner.Text()
		//fmt.Printf("%T\t%v\n", inode, inode)

		//Boy oh boy do we need some error handling here.
		fmt.Println(executer(structupdater(args, inode)))
		fmt.Println("Enter inode: ")
	}
	//Basic error handling on the scanner.
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

}

// Updates the struct that is passed to exec.Output() to include the current directory inode.
func structupdater(args []string, inode string) *exec.Cmd {
	if len(args) < 4 {
		args = append(args, inode)
		// fmt.Println("args were less than 4. Appending inode value. ", args)
	}
	args[3] = inode
	// fmt.Println("Updated args[3] with inode of:", inode)
	//Yuck, fix this.
	return exec.Command(cmd, args...)

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

//TODO
//func outputparser(output string) string {
//
//}
