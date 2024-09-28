package sitemaker

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Page struct {
	Layout  string        `xml:"layout"`
	View    string        `xml:"view"`
	Content []ContentItem `xml:"content"`
}

type ContentItem struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func RunCommand(cmd string, args []string) error {
	switch cmd {
	case "new":

		if len(args) != 1 {
			return errors.New("new command accepts 1 argument")
		}

		projectDir := args[0]

		err := createNewProject(projectDir)
		if err != nil {
			return err
		}
	case "gen":
		if len(args) < 2 {
			return errors.New("gen command invalid arguments")
		}

		sourceDir := args[0]
		outputDir := args[1]

		err := generateProject(sourceDir, outputDir)
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

func createNewProject(projectDir string) (err error) {

	layouts := projectDir + "/layouts"
	views := projectDir + "/views"
	data := projectDir + "/data"

	err = os.MkdirAll(layouts, 0755)

	if err != nil {
		return err
	}

	err = os.MkdirAll(views, 0755)

	if err != nil {
		return err
	}

	err = os.MkdirAll(data, 0755)

	if err != nil {
		return err
	}

	return nil
}

func generateProject(sourceDir, outputDir string) error {

	data, err := loadSourceFiles(sourceDir)

	if err != nil {
		return err
	}

	assets, err := loadAssetFiles(sourceDir)

	if err != nil {
		return err
	}

	for k, v := range data {
		var page Page
		err := xml.Unmarshal([]byte(v), &page)
		if err != nil {
			log.Fatalf("Error unmarshaling XML: %v", err)
		}

		// add source directory in the front
		page.Layout = sourceDir + "/" + page.Layout
		page.View = sourceDir + "/" + page.View

		// Parse templates
		tmpl, err := template.ParseFiles(page.Layout, page.View)
		if err != nil {
			log.Fatalf("Error parsing templates: %v", err)
		}

		fname := filepath.Base(k)
		if fname == "." {
			log.Fatal("Error parsing filepath")
		}

		// convert .xml to .html
		fname = strings.TrimSuffix(fname, filepath.Ext(fname)) + ".html"

		// Add output directory
		fname = outputDir + fname

		outputFile, err := os.Create(fname)
		if err != nil {
			log.Fatalf("Error creating output file: %v", err)
		}

		defer outputFile.Close()

		fmt.Println(page)

		err = tmpl.ExecuteTemplate(outputFile, fname, page)

		if err != nil {
			log.Fatalf("Error executing template: %v", err)
		}

		fmt.Println(k, v)
	}

	// TODO: move assets to output folder
	for k, v := range assets {
		fmt.Println(k, v)
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
