package httpserver

import (
	"context"
	"github.com/Eqke/metric-collector/internal/config"
	h "github.com/Eqke/metric-collector/internal/server/handlers"
	"github.com/Eqke/metric-collector/internal/server/middleware"
	stor "github.com/Eqke/metric-collector/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sync"
	"time"
)

type HTTPServer struct {
	ctx      context.Context
	server   *http.Server
	engine   *gin.Engine
	settings *config.ServerConfig
	logger   *zap.SugaredLogger
	wg       sync.WaitGroup
	storage  stor.Storage
}

func New(
	ctx context.Context,
	s *config.ServerConfig,
	storage stor.Storage,
	logger *zap.SugaredLogger,
	conn *pgx.Conn) *HTTPServer {

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(middleware.Logger(logger), middleware.Gzip(logger))
	r.GET("/", h.GetRootMetricsHandler(logger, storage))
	r.GET("/value/:type/:name", h.GETMetricHandler(logger, storage))
	r.GET("/ping", h.Ping(logger, conn))
	r.POST("/update/:type/:name/:value", h.POSTMetricHandler(logger, storage))
	r.POST("/update", h.POSTMetricJSONHandler(logger, storage))
	r.POST("/value", h.GetMetricJSONHandler(logger, storage))

	return &HTTPServer{
		server: &http.Server{
			Addr:    s.Endpoint,
			Handler: r,
		},
		engine:   r,
		settings: s,
		ctx:      ctx,
		logger:   logger,
		wg:       sync.WaitGroup{},
		storage:  storage,
	}
}

func (s *HTTPServer) restoreProcess() {
	defer s.wg.Done()
	s.logger.Info("Restore was started")
	t := time.NewTicker(time.Duration(s.settings.StoreInterval) * time.Second)
	defer t.Stop()
	for {
		select {
		case <-s.ctx.Done():
			{
				s.logger.Info("Restore was finished")
				return
			}
		case <-t.C:
			{
				s.logger.Info("Restored...")
				s.restore()
				s.logger.Info("Restore was finished")
			}
		}
	}
}

func (s *HTTPServer) restore() {
	func() {
		b, err := s.storage.ToJSON()
		if err != nil {
			s.logger.Errorf("Restore error getting storage: %v", err)
			return
		}
		f, err := os.OpenFile(s.settings.FileStoragePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			s.logger.Errorf("Restore error opening file: %v", err)
			return
		}
		defer f.Close()
		_, err = f.Write(b)
		if err != nil {
			s.logger.Errorf("Restore error writing file: %v", err)
			return
		}
	}()
}

func (s *HTTPServer) Run() {
	s.logger.Infof("Server was started. Listening on: %s", s.settings.Endpoint)

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Errorf("Server error: %v", err)
			os.Exit(1)
		}
	}()

	if s.settings.Restore {
		s.wg.Add(1)
		go s.restoreProcess()
	}

	<-s.ctx.Done()
	s.wg.Wait()
}

func (s *HTTPServer) Shutdown() {
	s.logger.Infof("Server was stopped.")
	err := s.server.Shutdown(context.Background())
	if err != nil {
		s.logger.Errorf("Server shutdown error: %v", err)
	}
}
