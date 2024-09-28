package httpserver

import (
	"context"
	"crypto/rsa"
	"github.com/Eqke/metric-collector/internal/server/config"
	"github.com/Eqke/metric-collector/internal/server/httpserver/handlers"
	middleware2 "github.com/Eqke/metric-collector/internal/server/httpserver/middleware"
	stor "github.com/Eqke/metric-collector/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
	"sync"
)

type HTTPServer struct {
	server  *http.Server
	engine  *gin.Engine
	logger  *zap.SugaredLogger
	storage stor.Storage
	host    string
}

func New(
	set *config.ServerConfig,
	storage stor.Storage,
	l *zap.SugaredLogger,
	key *rsa.PrivateKey,
) *HTTPServer {
	logger := l.Named("http-server")
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	rounter := gin.New()
	rounter.RedirectFixedPath = true

	logger.Infof("Server initing with %s storage", storage.Type())

	//usage middleware
	rounter.Use(
		middleware2.Logger(logger),
		middleware2.SubnetTrust(logger, set.TrustedSubnet),
		middleware2.Hash(logger, set.HashKey),
		middleware2.Decrypt(logger, key),
		middleware2.Gzip(logger),
	)

	rounter.GET("/", handlers.GetRootMetricsHandler(logger, storage))
	rounter.GET("/value/:type/:name/", handlers.GETMetricHandler(logger, storage))
	rounter.GET("/ping/", handlers.Ping(logger, storage))
	rounter.POST("/value/", handlers.GetMetricJSONHandler(logger, storage))
	rounter.POST("/update/:type/:name/:value/", handlers.POSTMetricHandler(logger, storage))
	rounter.POST("/update/", handlers.POSTMetricJSONHandler(logger, storage))
	rounter.POST("/updates/", handlers.PostMetricUpdates(logger, storage))

	//pproff tools api
	profiler := rounter.Group("/debug/pprof")
	{
		profiler.GET("/", gin.WrapF(pprof.Index))
		profiler.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		profiler.GET("/profile", gin.WrapF(pprof.Profile))
		profiler.POST("/symbol", gin.WrapF(pprof.Symbol))
		profiler.GET("/symbol", gin.WrapF(pprof.Symbol))
		profiler.GET("/trace", gin.WrapF(pprof.Trace))
		profiler.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		profiler.GET("/block", gin.WrapH(pprof.Handler("block")))
		profiler.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		profiler.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		profiler.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		profiler.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}

	return &HTTPServer{
		server: &http.Server{
			Addr:    set.Host,
			Handler: rounter,
		},
		engine:  rounter,
		logger:  logger,
		storage: storage,
		host:    set.Host,
	}
}

func (s *HTTPServer) Run(ctx context.Context, wg *sync.WaitGroup) {
	s.logger.Infof("Server was started. Listening on: %s", s.host)

	go func() {

		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Errorf("Server error: %v", err)
		}
		wg.Done()
	}()

	<-ctx.Done()
	s.Shutdown(ctx)
}

func (s *HTTPServer) Shutdown(ctx context.Context) {
	s.logger.Info("Server was stopped.")
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Errorf("Server shutdown error: %v", err)
	}
}
