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
	return nil
}

func generateProject(args []string) error {
	return nil
}

func printHelp() {
	fmt.Println("Help")
}
