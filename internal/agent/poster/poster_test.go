package poster

import (
	"github.com/Eqke/metric-collector/internal/agent/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestNewPoster(t *testing.T) {
	t.Run("poster_not_nil", func(t *testing.T) {
		l := zaptest.NewLogger(t).Sugar()
		poster := NewPoster(l, &config.AgentConfig{})
		require.NotNil(t, poster)
	})
}
