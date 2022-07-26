package main

type ISubmissionSource interface {
	getNextSubmissionData() *SubmissionData
	setResult(id int, result Result, executedTime float64, score int)
}
