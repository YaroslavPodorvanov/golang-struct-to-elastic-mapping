# Golang struct to elastic mapping

### Examples
```golang
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
```

### Original
* [Issue](https://github.com/olivere/elastic/issues/694)
* [Mapping template](https://github.com/olivere/elastic/blob/29ee98974cf1984dfecf53ef772d721fb97cb0b9/recipes/mapping/mapping.go#L26)

### Elasticsearch Guide
* [Field data types](https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping-types.html)

### TODOs
* Elasticsearch type "binary"
* Elasticsearch type "keyword"
