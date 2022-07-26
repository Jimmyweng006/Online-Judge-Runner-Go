package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Result int

const (
	Accepted Result = iota
	WrongAnswer
	CompileError
	RuntimeError
	TimeLimitExceeded
)

func (r Result) String() string {
	switch r {
	case Accepted:
		return "Accepted"
	case WrongAnswer:
		return "WrongAnswer"
	case CompileError:
		return "CompileError"
	case RuntimeError:
		return "RuntimeError"
	case TimeLimitExceeded:
		return "TimeLimitExceeded"
	default:
		return fmt.Sprintf("%d", int(r))
	}
}

const NO_EXECUTED_TIME = -1.0
const NO_SCORE = 0

type ResultState struct {
	Result       Result
	ExecutedTime float64
	TotalScore   int
}

type Judger struct {
	Compiler ICompiler
	Executor IExecutor
}

func (j *Judger) judge(submission *SubmissionData) ResultState {
	// compile code
	executableFilename := j.Compiler.compile(submission.Code)

	if &executableFilename == nil {
		return ResultState{
			Result:       CompileError,
			ExecutedTime: NO_EXECUTED_TIME,
			TotalScore:   NO_SCORE,
		}
	} else if _, err := os.Stat(executableFilename); err != nil {
		return ResultState{
			Result:       CompileError,
			ExecutedTime: NO_EXECUTED_TIME,
			TotalScore:   NO_SCORE,
		}
	}

	// execute code
	ans := j.execute(executableFilename, submission.TestCases)
	time.Sleep(2 * time.Second)
	err := os.Remove(executableFilename)
	if err != nil {
		fmt.Println(err)
	}

	return ans
}

func (j *Judger) execute(executableFilename string, testCases []TestCasesData) ResultState {
	isCorrect := true
	totalExecutedTime := 0.0
	totalScore := 0

	for _, testCase := range testCases {
		result := j.Executor.execute(executableFilename, testCase.Input, testCase.TimeOutSeconds)

		if &result == nil {
			return ResultState{
				Result:       RuntimeError,
				ExecutedTime: NO_EXECUTED_TIME,
				TotalScore:   NO_SCORE,
			}
		}

		if result.isTimeOut {
			return ResultState{
				Result:       TimeLimitExceeded,
				ExecutedTime: NO_EXECUTED_TIME,
				TotalScore:   NO_SCORE,
			}
		}

		if result.isCorrupted {
			return ResultState{
				Result:       RuntimeError,
				ExecutedTime: NO_EXECUTED_TIME,
				TotalScore:   NO_SCORE,
			}
		}

		output := strings.TrimSpace(result.output)
		expectedOutput := strings.TrimSpace(testCase.ExpectedOutput)

		totalExecutedTime += result.executedTime

		if output == expectedOutput {
			totalScore += testCase.Score
		} else {
			isCorrect = false
		}
	}

	if isCorrect {
		return ResultState{
			Result:       Accepted,
			ExecutedTime: totalExecutedTime,
			TotalScore:   totalScore,
		}
	} else {
		return ResultState{
			Result:       WrongAnswer,
			ExecutedTime: totalExecutedTime,
			TotalScore:   totalScore,
		}
	}
}
