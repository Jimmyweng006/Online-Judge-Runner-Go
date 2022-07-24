package main

type TestCaseTable struct {
	Id             int    `gorm:"auto_increment;primary_key;" json:"testcaseId"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expectedOutput"`
	// Comment        string  `json:"comment"`
	Score          int     `json:"score"`
	TimeOutSeconds float64 `json:"timeOutSeconds"`

	ProblemId int `gorm:"foreignKey:ProblemId" json:"problemId"`
}
