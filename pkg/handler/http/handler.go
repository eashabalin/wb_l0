package http

import (
	"github.com/gin-gonic/gin"
	"wb_l0/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	orders := router.Group("orders")
	{
		orders.GET("", h.getAllOrders)
		orders.GET("/:id", h.getOrderByUID)
	}

	return router
}
