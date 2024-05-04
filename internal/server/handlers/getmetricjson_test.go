package handlers

import (
	"github.com/Eqke/metric-collector/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestGetMetricJSONHandler(t *testing.T) {
	type args struct {
		logger *zap.SugaredLogger
		s      storage.Storage
	}
	tests := []struct {
		name string
		args args
		want gin.HandlerFunc
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetMetricJSONHandler(tt.args.logger, tt.args.s), "GetMetricJSONHandler(%v, %v)", tt.args.logger, tt.args.s)
		})
	}
}
