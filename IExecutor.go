package main

type IExecutor interface {
	execute(executableFilename string, input string) string
}
