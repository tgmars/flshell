package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nsf/termbox-go"
)

const (
	// Command to be passed to os to execute.
	cmd = "fls"
)

func main() {

	imagepath := os.Args[1]
	diskoffset := os.Args[2]

	args := []string{"-o", diskoffset, imagepath}
	cmdrequired := "fls"

	// Define a cmd struct that consists of the executable, it's location and the arguments passed.

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	for {
		cmdstruct := exec.Command(cmd, args...)
		switch cmdrequired {
		case "fls":
			fmt.Println(executer(cmdstruct))
		case "icat":
			fmt.Println("Do icat stuff on current file inode.")
		case "tsk_recover":
			fmt.Println("Do tsk_recover stuff on current directory inode.")
		default:
			//do nothing because we assume no command is needed.
		}

		switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey:
			switch event.Key {
			case termbox.KeyArrowUp:
				fmt.Println("ARROW UP")
				cmdrequired = ""
			}
		case termbox.EventError:
			panic(event.Err)
		}
		termbox.Flush()

	}

	// Execute fls initially to get inodes of the root directory.

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
