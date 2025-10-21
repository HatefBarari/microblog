package infrastructure

import (
	"net/http"

	"github.com/HatefBarari/microblog-media/internal/presenter"
	"github.com/HatefBarari/microblog-media/internal/usecase"
	"github.com/HatefBarari/microblog-shared/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoServer struct {
	server *echo.Echo
	config *Config
}

func NewEchoServer(config *Config) *EchoServer {
	e := echo.New()
	
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	
	// Request ID middleware
	e.Use(middleware.RequestID())
	
	// Body limit middleware
	e.Use(middleware.BodyLimit("50MB"))
	
	return &EchoServer{
		server: e,
		config: config,
	}
}

func (s *EchoServer) SetupRoutes(handler *presenter.HTTPHandler) {
	// Health check
	s.server.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API routes with authentication
	api := s.server.Group("/api/v1")
	api.Use(auth.JWTMiddleware("access_secret")) // You'll need to pass the secret from config
	
	// Media routes
	api.POST("/media/upload", handler.Upload)
	api.GET("/media", handler.List)
	api.GET("/media/:id", handler.GetByID)
	api.DELETE("/media/:id", handler.Delete)
	
	// Serve media files
	s.server.GET("/media/:filename", handler.Serve)
}

func (s *EchoServer) Start() error {
	addr := s.config.Server.Host + ":" + s.config.Server.Port
	return s.server.Start(addr)
}

func (s *EchoServer) Shutdown() error {
	return s.server.Close()
}
