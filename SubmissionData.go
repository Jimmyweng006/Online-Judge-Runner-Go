package main

type SubmissionData struct {
	Id        int             `json:"submissionId"`
	Language  string          `json:"language"`
	Code      string          `json:"code"`
	TestCases []TestCasesData `json:"testCases"`
}

type TestCasesData struct {
	Input          string  `json:"input"`
	ExpectedOutput string  `json:"expectedOutput"`
	Score          int     `json:"score"`
	TimeOutSeconds float64 `json:"timeOutSeconds"`
}
