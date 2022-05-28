package main

import (
	"Cubernetes/pkg/apiserver/objfile"
	"fmt"
)

func main() {
	objfile.PostJobFile("12345", "/Users/shen/Desktop/123.tar.gz")
	objfile.GetJobFile("12345", "./test.tar.gz")

	objfile.PostJobOutput("123", "yes,hhh!!!")
	fmt.Println(objfile.GetJobOutput("123"))
}
