package recrem

import (
	"net/http"

	"github.com/ShareCampus/RecRem/backend/pkg/recrem/middleware"
	"github.com/ShareCampus/RecRem/backend/pkg/utils/logger"
	"github.com/gorilla/handlers"
	"gorm.io/gorm"
)

const (
	DefaultHTTPListenPort = "8080"
)

type server struct {
	listenPort string
	httpServer *http.Server
	dbClient   *gorm.DB
}

func NewServer(opts ...func(*server) error) *server {
	var server = new(server)
	for _, opt := range opts {
		if err := opt(server); err != nil {
			logger.Fatalf("Failed to apply option: %v", err)
		}
	}
	if err := server.componentInspection(); err != nil {
		logger.Fatalf("Failed to inspect server components: %v", err)
	}
	server.httpServer = &http.Server{
		Addr: ":" + server.listenPort,
	}
	server.bindHandlers()
	return server
}

func (s *server) componentInspection() error {
	if s.listenPort == "" {
		logger.Infof("HTTP listen port is not set, using default port %s", DefaultHTTPListenPort)
		s.listenPort = DefaultHTTPListenPort
	}
	return nil
}

func (s *server) bindHandlers() {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /hello", s.hello)
	authedMux := middleware.Auth(apiMux)
	logMux := handlers.LoggingHandler(logger.NewLog(), authedMux)
	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/health", health)

	// strip prefix /api from all endpoints to avoid all handler in authHandler get /api prefix when get request url
	mainMux.Handle("/api/", http.StripPrefix("/api", logMux))

	s.httpServer.Handler = middleware.AllowCors(mainMux)

}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status": "OK"}`))
	if err != nil {
		panic(err)
	}
}

func ListenPort(listenPort string) func(*server) error {
	return func(s *server) error {
		s.listenPort = listenPort
		return nil
	}
}

func UseDB(dbClient *gorm.DB) func(*server) error {
	return func(s *server) error {
		s.dbClient = dbClient
		return nil
	}
}
