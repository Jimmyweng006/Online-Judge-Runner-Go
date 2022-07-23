package main

import (
	"fmt"
	"os"
)

const FILE_SUBMISSION_CODE_FILENAME = "file/code.txt"
const FILE_SUBMISSION_INPUT_FILENAME = "file/input.txt"
const FILE_SUBMISSION_OUTPUT_FILENAME = "file/output.txt"

type FileSubmissionSource struct {
	isGet bool
}

func (f *FileSubmissionSource) getNextSubmissionData() *SubmissionData {
	if f.isGet {
		return nil
	}

	codeFile, err := os.ReadFile(FILE_SUBMISSION_CODE_FILENAME)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	inputFile, err := os.ReadFile(FILE_SUBMISSION_INPUT_FILENAME)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	outputFile, err := os.ReadFile(FILE_SUBMISSION_OUTPUT_FILENAME)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	f.isGet = true
	var tempTestCases []TestCasesData
	tempTestCases = append(tempTestCases, TestCasesData{
		string(inputFile),
		string(outputFile),
		100,
		1.0,
	})
	return &SubmissionData{
		1,
		"kotlin",
		string(codeFile),
		tempTestCases,
	}
}

func (f *FileSubmissionSource) setResult(id int, result Result) {
	fmt.Printf("Submission %v: %v", id, result)
}
