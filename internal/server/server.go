package httpserver

import (
	"context"
	"github.com/Eqke/metric-collector/internal/handlers"
	"github.com/Eqke/metric-collector/internal/storage/localstorage"
	"log"
	"net/http"
	"strings"
)

type HTTPServer struct {
	server   *http.Server
	settings *Settings
	ctx      context.Context
}

type Settings struct {
	Host string
	Port string
}

func (s Settings) GetAddress() string {
	return strings.Join([]string{s.Host, s.Port}, ":")
}

func New(ctx context.Context, s *Settings) *HTTPServer {
	mux := http.NewServeMux()
	store := localstorage.New()

	mux.Handle(handlers.UpdatePath, middleware(handlers.UpdateHandler{}))
	mux.Handle(handlers.GaugePath, middleware(handlers.GaugeHandler{Storage: store}))
	mux.Handle(handlers.CounterPath, middleware(handlers.CounterHandler{Storage: store}))

	return &HTTPServer{
		server: &http.Server{
			Addr:    s.GetAddress(),
			Handler: mux,
		},
		settings: s,
		ctx:      ctx,
	}
}

func (s HTTPServer) Run() {
	log.Printf("Server was started.\n Listening on: %s", s.settings.GetAddress())
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve: %v", err)
		}
	}()
	<-s.ctx.Done()
	s.Shutdown()
}

func (s HTTPServer) Shutdown() {
	log.Println("Server was stopped.")
	err := s.server.Shutdown(s.ctx)
	if err != nil {
		log.Printf("Server shutdown failed: %v", err)
	}
}
