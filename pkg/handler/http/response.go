package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func newErrorResponse(c *gin.Context, statusCode int, err string) {
	fmt.Println(err)
	c.AbortWithStatusJSON(statusCode, map[string]any{
		"message": err,
	})
}
