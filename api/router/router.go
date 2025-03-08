package router

import (
	"fmt"
	"net/http"

	"mailbox-api/api/handler"
	"mailbox-api/api/middleware"
	"mailbox-api/config"
	"mailbox-api/logger"
	"mailbox-api/service"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
	config *config.Config
	logger *logger.Logger
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

func SetupRouter(cfg *config.Config, logger *logger.Logger, mailboxService service.MailboxService) *Router {
	router := &Router{
		engine: gin.New(),
		config: cfg,
		logger: logger,
	}

	router.engine.Use(gin.Recovery())
	router.engine.Use(middleware.LoggerMiddleware(logger))

	mailboxHandler := handler.NewMailboxHandler(mailboxService, logger)

	router.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.engine.Group("/api")
	{
		api.GET("/token/ceo", func(c *gin.Context) {
			token, err := middleware.GenerateToken(cfg, middleware.RoleCEO)
			if err != nil {
				logger.Error("Failed to generate CEO token", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"token": token})
		})

		api.GET("/token/cto", func(c *gin.Context) {
			token, err := middleware.GenerateToken(cfg, middleware.RoleCTO)
			if err != nil {
				logger.Error("Failed to generate CTO token", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"token": token})
		})

		// Protected routes for both CEO and CTO (with role-based filtering)
		mailboxes := api.Group("/mailboxes")
		mailboxes.Use(middleware.AuthMiddleware(cfg, logger))
		mailboxes.Use(middleware.RoleMiddleware(middleware.RoleCEO, middleware.RoleCTO))
		{
			mailboxes.GET("", mailboxHandler.GetMailboxes)
			mailboxes.GET("/:id", mailboxHandler.GetMailbox)

			// Only CEO can recalculate metrics
			calcMetrics := mailboxes.Group("/calculate-metrics")
			calcMetrics.Use(middleware.RoleMiddleware(middleware.RoleCEO))
			{
				calcMetrics.POST("", mailboxHandler.CalculateOrgMetrics)
			}
		}
	}

	return router
}

func (r *Router) Start(port int) *http.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r.engine,
	}

	go func() {
		r.logger.Info("Starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			r.logger.Fatal("Failed to start server", "error", err)
		}
	}()

	return srv
}
