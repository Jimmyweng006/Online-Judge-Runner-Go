package main

import (
	"fmt"
	"os"
	"strings"
)

type Result int

const (
	Accepted Result = iota
	WrongAnswer
)

func (r Result) String() string {
	switch r {
	case Accepted:
		return "Accepted"
	case WrongAnswer:
		return "WrongAnswer"
	default:
		return fmt.Sprintf("%d", int(r))
	}
}

type Judger struct {
	Compiler ICompiler
	Executor IExecutor
}

func (j *Judger) judge(submission *SubmissionData) Result {
	executableFilename := j.Compiler.compile(submission.Code)

	isCorrect := true
	for _, testCase := range submission.TestCases {
		output := strings.TrimSpace(j.Executor.execute(executableFilename, testCase.Input))
		expectedOutput := strings.TrimSpace(testCase.ExpectedOutput)

		if output != expectedOutput {
			isCorrect = false
			break
		}
	}

	err := os.Remove(executableFilename)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	var ans Result
	if isCorrect {
		ans = Accepted
	} else {
		ans = WrongAnswer
	}

	return ans
}
