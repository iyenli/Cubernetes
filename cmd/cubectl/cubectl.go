package main

import (
	"flag"
	"fmt"
	"os"
)

var argCmd = ""
var argType = ""
var argName = ""
var inFile string
var outFile string

func help() {
	fmt.Println("FATAL: wrong arguments")
	fmt.Println("cubectl [command] [type] [name] [flags]")
	fmt.Println("or")
	fmt.Println("cubectl [command] -f [filename] -o [output filename]")
}

func parseArgs() bool {
	if len(os.Args) < 4 {
		return false
	}

	var flagsStart = len(os.Args)
	for index, arg := range os.Args {
		if arg[0] == '-' {
			flagsStart = index
			break
		}
	}

	var args = os.Args[1:flagsStart]
	os.Args = os.Args[flagsStart-1:]
	flag.StringVar(&inFile, "f", "", "input filename")
	flag.StringVar(&outFile, "o", "", "output filename")
	flag.Parse()

	if len(args) == 1 && inFile != "" {
		argCmd = args[0]
		return true
	}

	if len(args) == 3 && inFile == "" {
		argCmd = args[0]
		argType = args[1]
		argName = args[2]
		return true
	}

	return false
}

func newObj(t string, name string, inFile string) {
	if inFile == "" {
	}
}

func delObj(t string, name string, inFile string) {

}

func getObj(t string, name string, inFile string) {

}

func apply(inFile string) {

}

func main() {
	if !parseArgs() {
		help()
		return
	}

	switch argCmd {
	case "create":
		newObj(argType, argName, inFile)
	case "delete":
		delObj(argType, argName, inFile)
	case "get":
		getObj(argType, argName, inFile)
	case "apply":
		apply(inFile)
	}
}
