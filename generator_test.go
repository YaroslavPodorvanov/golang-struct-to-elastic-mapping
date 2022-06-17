package mapping

import (
	"testing"
	"time"

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

		var result, err = Generate(&Empty{})

		require.Equal(t, expected, string(result))
		require.NoError(t, err)
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
        "type": "integer",
        "index": false
      }
    }
  }
}`

		var result, err = Generate(&User{})

		require.Equal(t, expected, string(result))
		require.NoError(t, err)
	}

	// tweet https://github.com/olivere/elastic/blob/29ee98974cf1984dfecf53ef772d721fb97cb0b9/recipes/mapping/mapping.go#L57
	{
		type Tweet struct {
			User     string                 `json:"user"`
			Message  string                 `json:"message"`
			Retweets int                    `json:"retweets"`
			Created  time.Time              `json:"created"`
			Attrs    map[string]interface{} `json:"attributes,omitempty"`
		}

		// language=JSON
		const expected = `{
  "mappings": {
    "properties": {
      "user": {
        "type": "keyword"
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

		var result, err = Generate(&Tweet{})

		require.Equal(t, expected, string(result))
		require.NoError(t, err)
	}
}

func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		type User struct {
			ID int `json:"id"`
		}

		_, _ = Generate(&User{})
	}
}
