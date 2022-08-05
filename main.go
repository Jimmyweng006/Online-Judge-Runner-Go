package main

import (
	"time"
)

func main() {
	// var submissionSource ISubmissionSource = &DatabaseSubmissionSource{}
	var submissionSource = &DatabaseSubmissionSource{}
	submissionSource.Init()

	for true {
		worker(submissionSource)
		// go worker(submissionSource)
	}
}

func worker(submissionSource ISubmissionSource) {
	submission := submissionSource.getNextSubmissionData()

	for submission != nil {
		judger := Judger{
			&KotlinCompiler{},
			&JVMExecutor{},
		}

		resultState := judger.judge(submission)
		submissionSource.setResult(submission.Id, resultState.Result, resultState.ExecutedTime, resultState.TotalScore)
		submission = submissionSource.getNextSubmissionData()
	}

	time.Sleep(5 * time.Second)
}
