package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	applicationjobs "github.com/saral-gupta7/recode/backend/internal/application/jobs"
	httpmiddleware "github.com/saral-gupta7/recode/backend/internal/httpapi/middleware"
)

type healthResponse struct {
	Status string `json:"status"`
}

type Dependencies struct {
	Jobs           *applicationjobs.Service
	MaxUploadBytes int64
	Readiness      func(context.Context) error
}

func NewRouter(logger *slog.Logger, dependencies ...Dependencies) *gin.Engine {
	router := gin.New()

	router.Use(httpmiddleware.RequestID())
	router.Use(httpmiddleware.RequestLogger(logger))
	router.Use(httpmiddleware.Recovery(logger))

	router.GET("/health/live", handleLiveness)
	if len(dependencies) > 0 && dependencies[0].Readiness != nil {
		router.GET("/health/ready", handleReadiness(dependencies[0].Readiness))
	}

	if len(dependencies) > 0 && dependencies[0].Jobs != nil {
		handler := jobsHandler{
			service:  dependencies[0].Jobs,
			maxBytes: dependencies[0].MaxUploadBytes,
		}
		api := router.Group("/api/v1")
		api.POST("/jobs", handler.create)
		api.GET("/jobs/:jobID", handler.get)
		api.POST("/jobs/:jobID/cancel", handler.cancel)
		api.GET("/jobs/:jobID/result", handler.result)
		api.DELETE("/jobs/:jobID", handler.delete)
	}

	return router
}

func handleReadiness(check func(context.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := check(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, healthResponse{Status: "unavailable"})
			return
		}
		c.JSON(http.StatusOK, healthResponse{Status: "ok"})
	}
}

func handleLiveness(c *gin.Context) {
	c.JSON(http.StatusOK, healthResponse{
		Status: "ok",
	})
}
