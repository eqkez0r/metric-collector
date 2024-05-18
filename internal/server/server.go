package httpserver

import (
	"context"
	"github.com/Eqke/metric-collector/internal/config"
	h "github.com/Eqke/metric-collector/internal/server/handlers"
	"github.com/Eqke/metric-collector/internal/server/middleware"
	stor "github.com/Eqke/metric-collector/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sync"
	"time"
)

type HTTPServer struct {
	server   *http.Server
	engine   *gin.Engine
	settings *config.ServerConfig
	logger   *zap.SugaredLogger
	wg       sync.WaitGroup
	storage  stor.Storage
	conn     *pgxpool.Pool
}

func New(
	ctx context.Context,
	s *config.ServerConfig,
	storage stor.Storage,
	logger *zap.SugaredLogger) (*HTTPServer, error) {

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	var conn *pgxpool.Pool = nil
	var err error

	if s.DatabaseDSN != "" {
		conn, err = pgxpool.New(ctx, s.DatabaseDSN)
		if err != nil {
			logger.Infof("Database connection error: %v", err)
			return nil, err
		}
	}

	logger.Infof("Server initing with %s storage", storage.Type())

	r.Use(middleware.Logger(logger), middleware.Hash(logger, s.HashKey), middleware.Gzip(logger))
	r.GET("/", h.GetRootMetricsHandler(logger, storage))
	r.GET("/value/:type/:name", h.GETMetricHandler(logger, storage))
	r.GET("/ping", h.Ping(logger, conn))
	r.POST("/update/:type/:name/:value", h.POSTMetricHandler(logger, storage))
	r.POST("/update", h.POSTMetricJSONHandler(logger, storage))
	r.POST("/value", h.GetMetricJSONHandler(logger, storage))
	r.POST("/updates", h.PostMetricUpdates(logger, storage))

	return &HTTPServer{
		server: &http.Server{
			Addr:    s.Endpoint,
			Handler: r,
		},
		engine:   r,
		settings: s,
		logger:   logger,
		wg:       sync.WaitGroup{},
		storage:  storage,
		conn:     conn,
	}, nil
}

func (s *HTTPServer) restoreProcess(ctx context.Context) {
	defer s.wg.Done()
	s.logger.Info("Restore was started")
	t := time.NewTicker(time.Duration(s.settings.StoreInterval) * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			{
				s.logger.Info("Restore was finished")
				return
			}
		case <-t.C:
			{
				s.logger.Info("Restored...")
				s.restore(ctx)
				s.logger.Info("Restore was finished")
			}
		}
	}
}

func (s *HTTPServer) restore(ctx context.Context) {
	func() {
		b, err := s.storage.ToJSON(ctx)
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

func (s *HTTPServer) Run(ctx context.Context) {
	s.logger.Infof("Server was started. Listening on: %s", s.settings.Endpoint)

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Errorf("Server error: %v", err)
		}
	}()

	if s.settings.Restore && s.settings.DatabaseDSN == "" {
		s.wg.Add(1)
		go s.restoreProcess(ctx)
	}

	<-ctx.Done()
	s.wg.Wait()
	s.Shutdown(ctx)
}

func (s *HTTPServer) Shutdown(ctx context.Context) {
	time.Sleep(time.Second * 5)
	s.logger.Infof("Server was stopped.")
	if s.conn != nil {
		s.conn.Close()
	}
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Errorf("Server shutdown error: %v", err)
	}
}
