package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewAgentConfig(t *testing.T) {
	t.Run("agent_config_not_nil", func(t *testing.T) {
		c, err := NewAgentConfig()
		if err != nil {
			t.Fatal(err)
		}
		require.NotNil(t, c)
	})
}
