package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

const JVM_INPUT_FILENAME = "input.txt"
const JVM_OUTPUT_FILENAME = "output.txt"

type JVMExecutor struct {
}

func (e *JVMExecutor) execute(executableFilename string, input string) string {
	var output string

	cmd := exec.Command(
		"java",
		"-jar",
		executableFilename,
	)
	// stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, input)
	}()
	// stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(stdout)
	line, err := reader.ReadString('\n')
	for err == nil {
		output += line
		line, err = reader.ReadString('\n')
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	return output
}
