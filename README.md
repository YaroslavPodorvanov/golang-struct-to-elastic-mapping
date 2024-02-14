# Golang struct to elastic mapping

## Terms of use<sup>[[?]](https://github.com/Tyrrrz/.github/blob/master/docs/why-so-political.md)</sup>

By using this project or its source code, for any purpose and in any shape or form, you grant your **implicit agreement** to all the following statements:

- You **condemn Russia and its military aggression against Ukraine**
- You **recognize that Russia is an occupant that unlawfully invaded a sovereign state**
- You **support Ukraine's territorial integrity, including its claims over temporarily occupied territories of Crimea and Donbas**
- You **reject false narratives perpetuated by Russian state propaganda**

To learn more about the war and how you can help, [click here](https://tyrrrz.me/ukraine). Glory to Ukraine! 🇺🇦

### Examples
```golang
package main

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
	ID              int      `json:"id" es:"index:true"`
	Title           string   `json:"title" es:"type:text"`
	Description     string   `json:"description" es:"type:text"`
	Company         *Company `json:"company"`
	RequiredSkills  []Alias  `json:"required_skills"`
	PreferredSkills []Alias  `json:"preferred_skills"`
	DesiredSkills   []Alias  `json:"desired_skills"`
}
```
```json
{
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
}
```
```golang
package main

import (
	"testing"

	"github.com/YaroslavPodorvanov/golang-struct-to-elastic-mapping/generator"

	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	// language=JSON
	const expected = `...`

	var result, err = generator.Generate(&Vacancy{})

	require.NoError(t, err)
	require.Equal(t, expected, string(result))
}
```

### Original
* [Issue](https://github.com/olivere/elastic/issues/694)
* [Mapping template](https://github.com/olivere/elastic/blob/29ee98974cf1984dfecf53ef772d721fb97cb0b9/recipes/mapping/mapping.go#L26)

### Elasticsearch Guide
* [Field data types](https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping-types.html)
