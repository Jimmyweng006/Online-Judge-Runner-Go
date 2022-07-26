package main

type IExecutor interface {
	execute(executableFilename string, input string, timeOutSeconds float64) IExecutorResult
}

type IExecutorResult struct {
	isTimeOut    bool
	isCorrupted  bool
	executedTime float64
	output       string
}
