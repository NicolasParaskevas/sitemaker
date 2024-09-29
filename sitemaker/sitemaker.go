package sitemaker

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
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

func CreateNewProject(projectDir string) (err error) {

	layouts := filepath.Join(projectDir, "layouts")
	views := filepath.Join(projectDir, "views")
	data := filepath.Join(projectDir, "data")
	assets := filepath.Join(projectDir, "assets")

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

	err = os.MkdirAll(assets, 0755)

	if err != nil {
		return err
	}

	return nil
}

func GenerateProject(sourceDir, outputDir string) error {
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

	// Create asset files
	for k, v := range assets {
		assetDir := strings.Replace(k, sourceDir, outputDir, 1)
		fmt.Println("Generated:", assetDir)
		makeDirs(assetDir, 0755)
		err = ioutil.WriteFile(assetDir, v, 0755)
		if err != nil {
			log.Printf("Error writing asset file: %v", err)
		}
	}

	return nil
}
