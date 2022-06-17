package mapping

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	// empty
	{
		type Empty struct {
		}

		// language=JSON
		const expect = `{
  "mappings": {
    "properties": {}
  }
}`

		var result, err = Generate(&Empty{})

		require.Equal(t, expect, string(result))
		require.NoError(t, err)
	}

	// empty
	{
		type User struct {
			ID int `json:"id"`
		}

		// language=JSON
		const expect = `{
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

		require.Equal(t, expect, string(result))
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
