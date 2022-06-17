package mapping

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	var result, err = Generate(nil)

	require.Equal(t, []byte(`{}`), result)
	require.NoError(t, err)
}

func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Generate(nil)
	}
}
