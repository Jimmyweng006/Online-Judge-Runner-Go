package main

func main() {
	var submissionSource ISubmissionSource
	submissionSource = &FileSubmissionSource{
		false,
	}
	submission := submissionSource.getNextSubmissionData()

	for submission != nil {
		judger := Judger{
			&KotlinCompiler{},
			&JVMExecutor{},
		}

		result := judger.judge(submission)
		submissionSource.setResult(submission.Id, result)
		submission = submissionSource.getNextSubmissionData()
	}
}
