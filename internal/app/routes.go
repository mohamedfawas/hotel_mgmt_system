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
	r.POST("/hotels", a.TenantMw.ResolveTenant(), a.HotelHandler.CreateHotel)
	r.POST("/hotels/:hotel_id/rooms", a.TenantMw.ResolveTenant(), a.RoomHandler.CreateRoom)
	r.GET("/hotels", a.TenantMw.ResolveTenant(), a.HotelHandler.ListHotels)
}
