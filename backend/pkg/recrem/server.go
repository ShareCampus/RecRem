package recrem

import (
	"github.com/ShareCampus/RecRem/backend/pkg/utils/logger"
	"gorm.io/gorm"
	"net/http"
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
