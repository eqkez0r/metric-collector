package poller

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestNewPoller(t *testing.T) {
	t.Run("poller_not_nil", func(t *testing.T) {
		l := zaptest.NewLogger(t).Sugar()
		poller := NewPoller(l)

		require.NotNil(t, poller)
	})
}
