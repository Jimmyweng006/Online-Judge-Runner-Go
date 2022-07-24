package main

import (
	"fmt"
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
	fmt.Println(submission)
	for submission != nil {
		judger := Judger{
			&KotlinCompiler{},
			&JVMExecutor{},
		}

		result := judger.judge(submission)
		submissionSource.setResult(submission.Id, result)
		submission = submissionSource.getNextSubmissionData()
	}

	time.Sleep(5 * time.Second)
}
