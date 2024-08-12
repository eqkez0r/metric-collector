package result

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewResult(t *testing.T) {
	t.Run("result_not_nil", func(t *testing.T) {
		res := New()

		require.NotNil(t, res)
	})
}
