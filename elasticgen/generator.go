package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	elastic_gen "github.com/YaroslavPodorvanov/golang-struct-to-elastic-mapping/generator"
)

type generator struct {
	destination string
}

func (g *generator) do(structs map[string]interface{}) error {
	if err := os.MkdirAll(g.destination, 0755); err != nil {
		return fmt.Errorf("failed to create destination folder: %w", err)
	}

	var (
		filecounter = 1
	)

	for name, s := range structs {
		result, err := elastic_gen.Generate(s)
		if err != nil {
			return fmt.Errorf("failed to generate Elasticsearch mapping: %w", err)
		}

		var (
			filename = fmt.Sprintf("%03d_%s.json", filecounter, camelToSnakeCase(name))
			filePath = filepath.Join(g.destination, filename)
		)

		if err = os.WriteFile(filePath, result, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		filecounter++
	}

	return nil
}

func camelToSnakeCase(s string) string {
	var result strings.Builder
	for i, char := range s {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(char)
	}
	return strings.ToLower(result.String())
}
