package routes

import (
	"haraka-sana/users/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(basePath *gin.RouterGroup) {
	oauth2 := basePath.Group("/auth")
	oauth2.GET("/register", handlers.Register)
}
