package mapping

import "encoding/json"

type Index struct {
	Mappings IndexMappings `json:"mappings"`
}

type IndexMappings struct {
	Properties json.RawMessage `json:"properties"`
}

type NestedProperties struct {
	Type       string          `json:"type"`
	Properties json.RawMessage `json:"properties"`
}

type Property struct {
	Type  string `json:"type"`
	Index *bool  `json:"index,omitempty"`
}
