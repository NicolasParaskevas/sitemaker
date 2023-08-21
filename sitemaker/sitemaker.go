package sitemaker

import (
	"errors"
	"fmt"
	"os"
)

func RunCommand(cmd string, args []string) error {
	switch cmd {
	case "new":
		err := createNewProject(args)
		if err != nil {
			return err
		}
	case "gen":
		err := generateProject(args)
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

func createNewProject(args []string) (err error) {

	if len(args) != 1 {
		return errors.New("new command accepts 1 argument")
	}

	projectDir := args[0]

	layouts := projectDir + "/layouts"
	components := projectDir + "/components"
	data := projectDir + "/data"

	err = os.MkdirAll(layouts, 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(components, 0755)

	if err != nil {
		return err
	}

	err = os.MkdirAll(data, 0755)

	if err != nil {
		return err
	}

	return nil
}

func generateProject(args []string) error {
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
