package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const SUPPORTED_LANGUAGE = "kotlin"

type DatabaseSubmissionSource struct {
	db  *gorm.DB
	rdb *redis.Client
}

func getConnection(rdb *redis.Client) error {
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()

	if err != nil {
		fmt.Println("ping error, try reconnect", err.Error())
		rdb = redis.NewClient(&redis.Options{})
		pong, err = rdb.Ping(ctx).Result()
		return err
	}

	fmt.Println("ping result:", pong)
	return nil
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
	d.rdb = redis.NewClient(&redis.Options{})

	// create tables
	db.Transaction(func(tx *gorm.DB) error {
		tx.AutoMigrate(&ProblemTable{}, &TestCaseTable{}, &SubmissionTable{})

		return nil
	})
}

func (d *DatabaseSubmissionSource) getNextSubmissionData() *SubmissionData {
	if err := getConnection(d.rdb); err != nil {
		return nil
	}

	ctx := context.Background()
	var data SubmissionData
	dataString, err := d.rdb.LPop(ctx, SUPPORTED_LANGUAGE).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(dataString), &data); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return &data
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
