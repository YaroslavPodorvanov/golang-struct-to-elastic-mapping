package mapping

import (
	"testing"
	"time"

	"github.com/YaroslavPodorvanov/golang-struct-to-elastic-mapping/generator"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	// empty
	{
		type Empty struct {
		}

		// language=JSON
		const expected = `{
  "mappings": {
    "properties": {}
  }
}`

		var result, err = generator.Generate(&Empty{})

		require.NoError(t, err)
		require.Equal(t, expected, string(result))
	}

	// empty
	{
		type User struct {
			ID int `json:"id"`
		}

		// language=JSON
		const expected = `{
  "mappings": {
    "properties": {
      "id": {
        "type": "integer"
      }
    }
  }
}`

		var result, err = generator.Generate(&User{})

		require.NoError(t, err)
		require.Equal(t, expected, string(result))
	}

	// tweet https://github.com/olivere/elastic/blob/29ee98974cf1984dfecf53ef772d721fb97cb0b9/recipes/mapping/mapping.go#L57
	{
		type Tweet struct {
			User     string                 `json:"user"`
			Message  string                 `json:"message"`
			Retweets int                    `json:"retweets"`
			Created  time.Time              `json:"created" es:"type:date"`
			Attrs    map[string]interface{} `json:"attributes" es:"type:object"`
		}

		// language=JSON
		const expected = `{
  "mappings": {
    "properties": {
      "user": {
        "type": "text"
      },
      "message": {
        "type": "text"
      },
      "retweets": {
        "type": "integer"
      },
      "created": {
        "type": "date"
      },
      "attributes": {
        "type": "object"
      }
    }
  }
}`

		var result, err = generator.Generate(&Tweet{})

		require.NoError(t, err)
		require.Equal(t, expected, string(result))
	}

	// tags
	{
		// language=JSON
		const expected = `{
  "mappings": {
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
      },
      "description": {
        "type": "text"
      },
      "employee_count": {
        "type": "integer",
        "index": false
      }
    }
  }
}`

		type Company struct {
			ID            int    `json:"id" es:"index:true"`
			Alias         string `json:"alias" es:"type:keyword,index:true"`
			Name          string `json:"name" es:"type:text"`
			Description   string `json:"description" es:"type:text"`
			EmployeeCount int    `json:"employee_count" es:"index:false"`
		}

		var result, err = generator.Generate(&Company{})

		require.NoError(t, err)
		require.Equal(t, expected, string(result))
	}

	{
		// language=JSON
		const expected = `{
  "mappings": {
    "properties": {
      "id": {
        "type": "integer"
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
            "type": "integer"
          },
          "alias": {
            "type": "text"
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
            "type": "text"
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
            "type": "text"
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
            "type": "text"
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
			Alias string `json:"alias"`
			Name  string `json:"name"`
		}

		type Company struct {
			ID    int    `json:"id"`
			Alias string `json:"alias"`
			Name  string `json:"name"`
		}

		type Vacancy struct {
			ID              int     `json:"id"`
			Title           string  `json:"title"`
			Description     string  `json:"description"`
			Company         Company `json:"company"`
			RequiredSkills  []Alias `json:"required_skills"`
			PreferredSkills []Alias `json:"preferred_skills"`
			DesiredSkills   []Alias `json:"desired_skills"`
		}

		var result, err = generator.Generate(&Vacancy{})

		require.NoError(t, err)
		require.Equal(t, expected, string(result))
	}
}

func BenchmarkGenerate(b *testing.B) {
type Alias struct {
	Alias string `json:"alias"`
	Name  string `json:"name"`
}

type Company struct {
	ID    int    `json:"id"`
	Alias string `json:"alias"`
	Name  string `json:"name"`
}

type Vacancy struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Company         Company `json:"company"`
	RequiredSkills  []Alias `json:"required_skills"`
	PreferredSkills []Alias `json:"preferred_skills"`
	DesiredSkills   []Alias `json:"desired_skills"`
}

	for i := 0; i < b.N; i++ {
		_, _ = generator.Generate(&Vacancy{})
	}
}
