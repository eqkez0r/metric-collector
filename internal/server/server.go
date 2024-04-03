package httpserver

import (
	"context"
	"github.com/Eqke/metric-collector/internal/handlers"
	"github.com/Eqke/metric-collector/internal/storage/localstorage"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type HTTPServer struct {
	server   *http.Server
	engine   *gin.Engine
	settings *Settings
	ctx      context.Context
}

type Settings struct {
	Endpoint string
}

func New(ctx context.Context, s *Settings) *HTTPServer {

	store := localstorage.New()
	gin.DisableConsoleColor()
	f, err := os.Create("gin-metric.log")
	if err != nil {
		log.Println(err)
	}
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.New()

	r.GET("/", handlers.GetRootMetricsHandler(store))
	r.GET("/value/:type/:name", handlers.GETMetricHandler(store))
	r.POST("/update/:type/:name/:value", handlers.POSTMetricHandler(store))

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
	signal.NotifyContext(s.ctx, syscall.SIGINT, syscall.SIGTERM)
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
