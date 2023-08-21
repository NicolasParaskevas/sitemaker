package main

import (
	"log"
	"os"

	"github.com/nicolasparaskevas/sitemaker/sitemaker"
)

func main() {

	command := os.Args[1] // first argument is the command
	args := os.Args[2:]   // the rest are the command arguments

	err := sitemaker.RunCommand(command, args)

	if err != nil {
		log.Fatal(err)
	}

}
