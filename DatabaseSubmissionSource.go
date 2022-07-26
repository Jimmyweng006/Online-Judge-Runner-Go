package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseSubmissionSource struct {
	db *gorm.DB
}

func initDatabase() (db *gorm.DB, err error) {
	dsn := "host=localhost user=postgres password=123456789 " +
		"dbname=onlinejudge-go port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	a, err := db.DB()
	a.SetMaxIdleConns(10)
	a.SetMaxOpenConns(100)

	return db, err
}

func (d *DatabaseSubmissionSource) Init() {
	db, err := initDatabase()
	if err != nil {
		fmt.Println(err)
		return
	}
	d.db = db

	// create tables
	db.Transaction(func(tx *gorm.DB) error {
		tx.AutoMigrate(&ProblemTable{}, &TestCaseTable{}, &SubmissionTable{})

		return nil
	})
}

func (d *DatabaseSubmissionSource) getNextSubmissionData() *SubmissionData {
	var submissionData *SubmissionData = nil
	var submission SubmissionTable

	d.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where(&SubmissionTable{Result: "-"}).First(&submission)

		if result.RowsAffected != 0 {
			// find bug on getProblemByID
			// rows, err := tx.Model(&TestCaseTable{ProblemId: submission.ProblemId}).Rows()
			rows, err := tx.Model(&TestCaseTable{}).Where("problem_id = ?", submission.ProblemId).Rows()
			defer rows.Close()
			if err != nil {
				return err
			}

			var testcases []TestCasesData
			for rows.Next() {
				var testcase TestCaseTable
				tx.ScanRows(rows, &testcase)

				temp := TestCasesData{
					Input:          testcase.Input,
					ExpectedOutput: testcase.ExpectedOutput,
					Score:          testcase.Score,
					TimeOutSeconds: testcase.TimeOutSeconds,
				}
				testcases = append(testcases, temp)
			}

			submissionData = &SubmissionData{
				Id:        submission.Id,
				Language:  submission.Language,
				Code:      submission.Code,
				TestCases: testcases,
			}
		}

		return nil
	})

	return submissionData
}

func (d *DatabaseSubmissionSource) setResult(id int, result Result, executedTime float64, score int) {
	d.db.Transaction(func(tx *gorm.DB) error {
		tx.Model(&SubmissionTable{Id: id}).Updates(SubmissionTable{
			Result:       result.String() + fmt.Sprintf("(%d)", score),
			ExecutedTime: executedTime,
		})

		return nil
	})

	fmt.Printf("Submission %v: %v - Score: %v (%v)\n", id, result, score, executedTime)
}
