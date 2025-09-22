package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONResponse(c *gin.Context, status int, success bool, message string, data interface{}) {
	c.JSON(status, gin.H{
		"success": success,
		"message": message,
		"data":    data,
	})
}

func Success(c *gin.Context, message string, data interface{}) {
	JSONResponse(c, http.StatusOK, true, message, data)
}

func Error(c *gin.Context, status int, message string) {
	JSONResponse(c, status, false, message, nil)
}
