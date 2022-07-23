package main

type ICompiler interface {
	compile(code string) string
}
