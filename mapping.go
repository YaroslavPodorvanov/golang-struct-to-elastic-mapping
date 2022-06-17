package mapping

import "encoding/json"

type Index struct {
	Mappings IndexMappings `json:"mappings"`
}

type IndexMappings struct {
	Properties json.RawMessage `json:"properties"`
}
