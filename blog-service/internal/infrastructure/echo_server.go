package infrastructure

import (
	"strconv"

	"github.com/HatefBarari/microblog-blog/internal/presenter"
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

	// JWT Middleware – secret بعداً از auth-service می‌آید (حالا خالی)
	jwtMid := auth.Middleware("")

	// Public
	e.GET("/articles", handler.ListArticles)
	e.GET("/articles/:slug", handler.GetArticleBySlug)
	e.GET("/categories/tree", handler.ListCategoryTree)

	// Protected (نیاز به JWT)
	articleGroup := e.Group("/articles", jwtMid)
	articleGroup.POST("", handler.CreateArticle)
	articleGroup.PUT("/:id", handler.UpdateArticle)
	articleGroup.DELETE("/:id", handler.DeleteArticle)
	articleGroup.GET("/:id/comments", handler.ListComments)
	articleGroup.POST("/:id/comments", handler.CreateComment)
	articleGroup.POST("/:id/rating", handler.RateArticle)
	articleGroup.DELETE("/:id/rating", handler.DeleteRating)

	// Manager/Admin only
	managerGroup := e.Group("/categories", jwtMid)
	managerGroup.POST("", handler.CreateCategory)

	log.Info("starting blog server", zap.Int("port", cfg.Server.Port))
	return e.Start(":" + strconv.Itoa(cfg.Server.Port))
}