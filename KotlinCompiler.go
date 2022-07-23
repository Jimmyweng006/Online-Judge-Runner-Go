package main

import (
	"log"
	"os"
	"os/exec"
)

const KOTLIN_CODE_FILENAME = "_code.kt"
const KOTLIN_CODE_EXECUTABLE_FILENAME = "_code.jar"

type KotlinCompiler struct {
}

func (c *KotlinCompiler) compile(code string) string {
	f, createErr := os.Create(KOTLIN_CODE_FILENAME)
	if createErr != nil {
		log.Fatal(createErr)
	}

	defer f.Close()

	_, writeErr := f.WriteString(code)
	if writeErr != nil {
		log.Fatal(writeErr)
	}

	cmd := exec.Command(
		"kotlinc",
		KOTLIN_CODE_FILENAME,
		"-include-runtime",
		"-d",
		KOTLIN_CODE_EXECUTABLE_FILENAME,
	)
	cmd.Run()

	os.Remove(KOTLIN_CODE_FILENAME)
	return KOTLIN_CODE_EXECUTABLE_FILENAME
}
