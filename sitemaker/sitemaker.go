package sitemaker

import (
	"errors"
	"fmt"
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

func createNewProject(args []string) error {

	// source dir must contain:
	//  layouts:
	//  		- homelayout.html
	// 		- detailslayout.html
	//  components:
	// 		- image.html
	// 		- project-item.html
	// 		- contact-form.html

	//  data:
	//   	- homepage.yml
	//   	- projcets.yml
	//   	- about.yml`

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
