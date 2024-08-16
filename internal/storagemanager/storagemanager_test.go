package storagemanager

import (
	"context"
	"github.com/Eqke/metric-collector/internal/server/config"
	"testing"

	"github.com/Eqke/metric-collector/internal/storage/localstorage"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestGetStorage(t *testing.T) {

	t.Run("getting_local_database", func(t *testing.T) {
		l := zaptest.NewLogger(t).Sugar()
		cfg := &config.ServerConfig{
			FileStoragePath: "./test_path.json",
		}

		store, err := GetStorage(context.Background(), l, cfg)
		if err != nil {
			t.Fatal(err)
		}
		defer store.Close()
		_, ok := store.(*localstorage.LocalStorage)
		require.Equal(t, true, ok)
	})
}
