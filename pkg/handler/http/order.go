package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getOrderByUID(c *gin.Context) {
	uid := c.Param("id")

	order, exists := h.services.Order.GetByUID(uid)
	if !exists {
		newErrorResponse(c, http.StatusNotFound, "order not found")
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *Handler) getAllOrders(c *gin.Context) {
	orders := h.services.GetAllFromDB()
	c.JSON(http.StatusOK, orders)
}
