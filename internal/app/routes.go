package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) registerRoutes() {
	r := a.Router.Group("/api/v1")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})
}
