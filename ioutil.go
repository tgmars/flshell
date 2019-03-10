package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeCmdToFile(filename string, cmdstruct *exec.Cmd) {
	outfile, err := os.Create("./" + filename)
	check(err)

	defer outfile.Close()

	stdoutPipe, err := cmdstruct.StdoutPipe()
	check(err)

	writer := bufio.NewWriter(outfile)
	defer writer.Flush()

	err = cmdstruct.Start()
	check(err)

	go io.Copy(writer, stdoutPipe)
	cmdstruct.Wait()
}

func writeStringToFile(filename string, input string) {
	f, err := os.Create("./" + filename)
	check(err)

	defer f.Close()

	_, err = f.WriteString(input)
	check(err)

	f.Sync()

}
