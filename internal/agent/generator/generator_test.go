package generator

import (
	"github.com/Eqke/metric-collector/internal/agent/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	t.Run("generator_not_nil", func(t *testing.T) {
		l := zaptest.NewLogger(t).Sugar()
		gen := NewGenerator(l, &config.AgentConfig{
			RateLimit: 100,
		})

		require.NotNil(t, gen)
	})
}
