package generator

import (
	"encoding/json"

	"github.com/YaroslavPodorvanov/golang-struct-to-elastic-mapping/mapping"
)

func wrap(properties []byte) ([]byte, error) {
	return json.MarshalIndent(&mapping.Index{
		Mappings: mapping.IndexMappings{
			Properties: properties,
		},
	}, "", "  ")
}
