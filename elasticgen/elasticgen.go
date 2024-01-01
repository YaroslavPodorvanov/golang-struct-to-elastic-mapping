package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const defaultDestinationPath = "./internal/elasticgen"

func main() {
	source := flag.String("source", "", "Input Go source file or folder.")
	destination := flag.String("destination", "", "Output folder.")
	flag.Parse()

	var (
		structs map[string]interface{}
		err     error
	)

	// If source is not specified, use the caller's directory
	if *source == "" {
		executablePath, err := os.Executable()
		if err == nil {
			*source = filepath.Dir(executablePath)
		} else {
			log.Fatal("Unable to determine caller's directory:", err)
		}
	}

	var p = new(parser)
	if strings.HasSuffix(*source, ".go") {
		p.parserType = parserTypeFile
	} else {
		p.parserType = parserTypeFolder
	}
	p.path = *source

	structs, err = p.do()
	if err != nil {
		log.Fatal("Failed to parse Go source code:", err)
	}
	if len(structs) < 1 {
		log.Fatal("No structs found in this path")
	}

	var g = new(generator)
	if *destination == "" {
		g.destination = defaultDestinationPath
	} else {
		g.destination = *destination
	}

	if err := g.do(structs); err != nil {
		log.Fatal("Failed to generate Elasticsearch mapping:", err)
	}

	fmt.Printf("Elasticsearch mapping generated successfully and saved to %s\n", *destination)
}
