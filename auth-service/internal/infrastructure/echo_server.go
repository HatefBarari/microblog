package infrastructure

import (
	"strconv"

	"github.com/HatefBarari/microblog-auth/internal/presenter"
	"github.com/HatefBarari/microblog-shared/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func StartEcho(cfg *Config, log *zap.Logger, handler *presenter.HTTPHandler) error {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info("request",
				zap.String("method", c.Request().Method),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))

	// JWT middleware
	jwtMid := auth.Middleware(cfg.Auth.AccessSecret)

	// public routes
	e.POST("/register", handler.Register)
	e.POST("/login", handler.Login)
	e.POST("/auth/refresh", handler.Refresh)
	
	// email verification routes
	e.GET("/verify", handler.VerifyEmail)
	e.POST("/forgot-password", handler.SendPasswordResetEmail)
	e.POST("/reset-password", handler.ResetPassword)

	// protected routes
	protected := e.Group("/api/v1", jwtMid)
	protected.POST("/resend-verification", handler.ResendVerificationEmail)
	protected.GET("/me", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"userID": c.Get("userID").(string)})
	})

	log.Info("starting auth server", zap.Int("port", cfg.Server.Port))
	return e.Start(":" + strconv.Itoa(cfg.Server.Port))
}