package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nicolasparaskevas/sitemaker/sitemaker"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("please provide a command")
	}

	command := os.Args[1] // first argument is the command
	args := os.Args[2:]   // the rest are the command arguments

	err := runCommand(command, args)

	if err != nil {
		log.Fatal(err)
	}

}

func runCommand(cmd string, args []string) error {
	switch cmd {
	case "new":

		if len(args) != 1 {
			return errors.New("new command accepts 1 argument")
		}

		projectDir := args[0]

		err := sitemaker.CreateNewProject(projectDir)
		if err != nil {
			return err
		}
	case "gen":
		if len(args) < 2 {
			return errors.New("gen command invalid arguments")
		}

		sourceDir := filepath.Clean(args[0])
		outputDir := filepath.Clean(args[1])

		err := sitemaker.GenerateProject(sourceDir, outputDir)
		if err != nil {
			return err
		}
	case "help":
		printHelp()
	default:
		return errors.New("invalid command: " + cmd)
	}

	return nil
}

func printHelp() {
	help := `Usage sitemaker [command] [argument]
Commands:
	-> gen [source dir] [output dir]
	 	Generates site on the output directory basted on the source dir
	-> new [source dir]
	 	Creates new project structure
	-> help
		Prints all available commands`

	fmt.Println(help)
}
