package httpserver

import (
	"context"
	"github.com/Eqke/metric-collector/internal/config"
	"github.com/Eqke/metric-collector/internal/handlers"
	stor "github.com/Eqke/metric-collector/internal/storage"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

type HTTPServer struct {
	server   *http.Server
	engine   *gin.Engine
	settings *config.ServerConfig
	ctx      context.Context
}

func New(ctx context.Context, s *config.ServerConfig, storage stor.Storage) *HTTPServer {

	gin.DisableConsoleColor()
	f, err := os.Create("gin-metric.log")
	if err != nil {
		log.Println(err)
	}
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.New()

	r.GET("/", handlers.GetRootMetricsHandler(storage))
	r.GET("/value/:type/:name", handlers.GETMetricHandler(storage))
	r.POST("/update/:type/:name/:value", handlers.POSTMetricHandler(storage))

	return &HTTPServer{
		server: &http.Server{
			Addr:    s.Endpoint,
			Handler: r,
		},
		engine:   r,
		settings: s,
		ctx:      ctx,
	}
}

func (s HTTPServer) Run() {
	log.Printf("Server was started.\n Listening on: %s", s.settings.Endpoint)
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Fatalf("listen and serve: %v", err)
		}
	}()
	<-s.ctx.Done()
	s.Shutdown()
}

func (s HTTPServer) Shutdown() {
	log.Println("Server was stopped.")
	err := s.server.Shutdown(context.Background())
	if err != nil {
		log.Printf("Server shutdown failed: %v", err)
	}
}
