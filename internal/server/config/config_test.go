package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewServerConfig(t *testing.T) {
	c, err := NewServerConfig()
	if err != nil {
		t.Fatal(err)
	}

	require.NotNil(t, c)
}
