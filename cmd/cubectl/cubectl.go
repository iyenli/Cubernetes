package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

func fileError(filename string) {
	fmt.Printf("FATAL: CANNOT open %s\n", filename)
}

func parseError(filename string) {
	fmt.Printf("FATAL: CANNOT parse %s\n", filename)
}

func executeError(ret int) {
	fmt.Printf("FATAL: execution error: %v\n", ret)
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

func newObj(config Config) int {
	var ret = -1
	switch config.getKind() {
	case "Pod":
		var podConfig = config.(PodConfig)
		fmt.Println(podConfig)
		ret = 0
	}
	return ret
}

func delObj(config Config) int {
	var ret = -1
	return ret
}

func getObj(config Config) int {
	var ret = -1
	return ret
}

func apply(config Config) int {
	var ret = -1
	return ret
}

func main() {
	if !parseArgs() {
		help()
		return
	}

	if inFile == "" {
	}
	configFile, err := ioutil.ReadFile(inFile)
	if err != nil {
		fileError(inFile)
	}

	var config = parseConfig(configFile)
	if config == nil {
		parseError(inFile)
		return
	}

	var ret = -1
	switch argCmd {
	case "create":
		ret = newObj(config)
	case "delete":
		ret = delObj(config)
	case "get":
		ret = getObj(config)
	case "apply":
		ret = apply(config)
	}

	if ret < 0 {
		executeError(ret)
		return
	}
}
