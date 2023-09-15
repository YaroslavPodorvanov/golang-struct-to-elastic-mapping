package main

import (
	"testing"

	"github.com/YaroslavPodorvanov/golang-struct-to-elastic-mapping/generator"

	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	// language=JSON
	const expected = `{
  "mappings": {
    "properties": {
      "id": {
        "type": "integer",
        "index": true
      },
      "title": {
        "type": "text"
      },
      "description": {
        "type": "text"
      },
      "company": {
        "type": "nested",
        "properties": {
          "id": {
            "type": "integer",
            "index": true
          },
          "alias": {
            "type": "keyword",
            "index": true
          },
          "name": {
            "type": "text"
          }
        }
      },
      "required_skills": {
        "type": "nested",
        "properties": {
          "alias": {
            "type": "keyword",
            "index": true
          },
          "name": {
            "type": "text"
          }
        }
      },
      "preferred_skills": {
        "type": "nested",
        "properties": {
          "alias": {
            "type": "keyword",
            "index": true
          },
          "name": {
            "type": "text"
          }
        }
      },
      "desired_skills": {
        "type": "nested",
        "properties": {
          "alias": {
            "type": "keyword",
            "index": true
          },
          "name": {
            "type": "text"
          }
        }
      }
    }
  }
}`

	type Alias struct {
		Alias string `json:"alias" es:"type:keyword,index:true"`
		Name  string `json:"name" es:"type:text"`
	}

	type Company struct {
		ID    int    `json:"id" es:"index:true"`
		Alias string `json:"alias" es:"type:keyword,index:true"`
		Name  string `json:"name" es:"type:text"`
	}

	type Vacancy struct {
		ID              int     `json:"id" es:"index:true"`
		Title           string  `json:"title" es:"type:text"`
		Description     string  `json:"description" es:"type:text"`
		Company         *Company `json:"company"`
		RequiredSkills  []Alias `json:"required_skills"`
		PreferredSkills []Alias `json:"preferred_skills"`
		DesiredSkills   []Alias `json:"desired_skills"`
	}

	var result, err = generator.Generate(&Vacancy{})

	require.NoError(t, err)
	require.Equal(t, expected, string(result))
}
