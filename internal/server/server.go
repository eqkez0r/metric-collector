package httpserver

import (
	"context"
	"github.com/Eqke/metric-collector/internal/config"
	handlers2 "github.com/Eqke/metric-collector/internal/server/handlers"
	"github.com/Eqke/metric-collector/internal/server/middleware"
	stor "github.com/Eqke/metric-collector/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
)

type HTTPServer struct {
	ctx      context.Context
	server   *http.Server
	engine   *gin.Engine
	settings *config.ServerConfig
	logger   *zap.SugaredLogger
}

func New(
	ctx context.Context,
	s *config.ServerConfig,
	storage stor.Storage,
	logger *zap.SugaredLogger) *HTTPServer {

	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(middleware.Logger(logger))
	r.GET("/", handlers2.GetRootMetricsHandler(logger, storage))
	r.GET("/value/:type/:name", handlers2.GETMetricHandler(logger, storage))
	r.POST("/update/:type/:name/:value", handlers2.POSTMetricHandler(logger, storage))

	return &HTTPServer{
		server: &http.Server{
			Addr:    s.Endpoint,
			Handler: r,
		},
		engine:   r,
		settings: s,
		ctx:      ctx,
		logger:   logger,
	}
}

func (s HTTPServer) Run() {
	s.logger.Infof("Server was started.\n Listening on: %s", s.settings.Endpoint)
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Errorf("Server error: %v", err)
			os.Exit(1)
		}
	}()
	<-s.ctx.Done()
	s.Shutdown()
}

func (s HTTPServer) Shutdown() {
	s.logger.Infof("Server was stopped.")
	err := s.server.Shutdown(context.Background())
	if err != nil {
		s.logger.Errorf("Server shutdown error: %v", err)
	}
}
