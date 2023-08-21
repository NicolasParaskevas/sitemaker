package main

import "fmt"

func main() {
	// Usage sitemaker [command] [argument]
	// Commands:
	// -> gen [source dir] [output dir]
	// 	Generates site on the output directory basted on the source dir
	// -> new [source dir]
	// 	Creates new project structure
	// -> help
	//   Prints all available commands
	//
	// source dir must contain:
	// layouts:
	// 		- homelayout.html
	//		- detailslayout.html
	// components:
	//		- image.html
	//		- project-item.html
	//		- contact-form.html
	//
	// data:
	//  	- homepage.yml
	//  	- projcets.yml
	//  	- about.yml

	fmt.Println("Hello World!")
}
