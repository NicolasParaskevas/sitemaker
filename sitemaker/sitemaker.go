package sitemaker

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Page struct {
	Layout  string `xml:"layout"`
	View    string `xml:"view"`
	Content map[string]interface{}
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

	layouts := filepath.Join(projectDir, "layouts")
	views := filepath.Join(projectDir, "views")
	data := filepath.Join(projectDir, "data")

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

		// Unmarshal the layout and view sections (static)
		err := xml.Unmarshal([]byte(v), &page)
		if err != nil {
			log.Printf("Error unmarshaling XML: %v", err)
			continue
		}

		// Now handle <content> part and unmarshal it into a map
		page.Content = make(map[string]interface{})
		decoder := xml.NewDecoder(strings.NewReader(v))
		var currentElement string

		for {
			// Read tokens from the XML
			token, err := decoder.Token()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("Error decoding XML: %v", err)
				return err
			}

			switch elem := token.(type) {
			case xml.StartElement:
				// Check if it's a <content> element
				if elem.Name.Local == "content" {
					// Read into a map dynamically
					for {
						innerToken, _ := decoder.Token()
						if innerToken == nil {
							break
						}

						switch innerElem := innerToken.(type) {
						case xml.StartElement:
							currentElement = innerElem.Name.Local
						case xml.CharData:
							contentValue := strings.TrimSpace(string(innerElem))

							if contentValue == "" {
								continue
							}

							if existingValue, ok := page.Content[currentElement]; ok {
								switch v := existingValue.(type) {
								case []string:
									// Append to the existing slice
									page.Content[currentElement] = append(v, contentValue)
								case string:
									// Convert single string to a slice of strings
									page.Content[currentElement] = []string{v, contentValue}
								}
							} else {
								// Set as a single string
								page.Content[currentElement] = contentValue
							}
						}
					}
				}
			}
		}

		// Add source directory in the front
		page.Layout = filepath.Join(sourceDir, page.Layout)
		page.View = filepath.Join(sourceDir, page.View)

		// Parse templates
		tmpl, err := template.ParseFiles(page.Layout, page.View)
		if err != nil {
			log.Printf("Error parsing templates: %v", err)
			continue
		}

		fname := filepath.Base(k)
		if fname == "." {
			log.Println("Error parsing filepath")
			continue
		}

		// Convert .xml to .html
		fname = strings.TrimSuffix(fname, filepath.Ext(fname)) + ".html"

		// Add output directory
		outputFilePath := filepath.Join(outputDir, fname)

		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			log.Printf("Error creating output file: %v", err)
			continue
		}
		defer outputFile.Close()

		err = tmpl.ExecuteTemplate(outputFile, filepath.Base(page.Layout), page.Content)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			continue
		}

		fmt.Println("Generated:", outputFilePath)
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
