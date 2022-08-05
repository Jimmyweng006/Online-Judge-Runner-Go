package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"
)

const JVM_INPUT_FILENAME = "input.txt"
const JVM_OUTPUT_FILENAME = "output.txt"

type JVMExecutor struct {
}

func (e *JVMExecutor) execute(executableFilename string, input string, timeOutSeconds float64) IExecutorResult {
	startTime := time.Now().UnixMilli()
	var output string

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOutSeconds*float64(time.Second)))
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
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
		fmt.Println("cmd run error")
	}

	isFinished := true
	go func() {
		// blocking here
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			isFinished = false
		}
	}()

	// read output
	reader := bufio.NewReader(stdout)
	line, err := reader.ReadString('\n')
	for err == nil {
		output += line
		line, err = reader.ReadString('\n')
	}
	// Call Wait after reaching EOF.
	isCorrupted := false
	if err := cmd.Wait(); err != nil {
		isCorrupted = true
	}

	executedTime := time.Now().UnixMilli() - startTime

	return IExecutorResult{
		isTimeOut:    !isFinished,
		isCorrupted:  isCorrupted,
		executedTime: float64(executedTime) / 1000,
		output:       output,
	}
}
