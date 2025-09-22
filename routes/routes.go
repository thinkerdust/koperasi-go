package routes

import (
	"koperasi-go/api"
	"koperasi-go/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	apiRoute := r.Group("/api")
	{
		apiRoute.POST("/login", api.Login)
		apiRoute.POST("/register", api.Register)

		protected := apiRoute.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/user", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Protected route OK"})
			})
		}
	}
}
