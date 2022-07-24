package main

type ProblemTable struct {
	Id int `gorm:"auto_increment;primary_key;" json:"problemId"`
	// Title       string `json:"title"`
	// Description string `json:"description"`
}
