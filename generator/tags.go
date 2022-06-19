package generator

import "strings"

type tags struct {
	typeName string
	index    *bool
}

func parseTags(s string) tags {
	var result tags

	var pairs = strings.Split(s, ",")
	for _, pair := range pairs {
		var key, value, _ = strings.Cut(pair, ":")
		switch key {
		case "type":
			result.typeName = value
		case "index":
			var index = value == "true"

			result.index = &index
		}
	}

	return result
}
