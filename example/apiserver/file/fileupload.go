package main

import (
	"Cubernetes/pkg/apiserver/jobfile"
	"fmt"
)

func main() {
	jobfile.PostJobFile("12345", "/Users/shen/Desktop/123.tar.gz")
	jobfile.GetJobFile("12345", "./test.tar.gz")

	jobfile.PostJobOutput("123", "yes,hhh!!!")
	fmt.Println(jobfile.GetJobOutput("123"))
}
