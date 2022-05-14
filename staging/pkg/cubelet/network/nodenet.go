package network

import (
	"Cubernetes/pkg/cubenetwork/weaveplugins"
	"log"
)

func InitNodeNetwork(args []string) {
	var err error
	if len(args) == 2 {
		err = weaveplugins.InitWeave()
	} else if len(args) == 3 {
		//err = weaveplugins.AddNode()
	} else {
		panic("Error: too much or little args when start cubelet;")
	}

	if err != nil {
		log.Panicf("Init weave network failed, err: %v", err.Error())
		return
	}
}
