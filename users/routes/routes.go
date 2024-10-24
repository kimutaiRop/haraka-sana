package routes

import (
	"haraka-sana/helpers/middleware"
	"haraka-sana/users/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(basePath *gin.RouterGroup) {
	oauth2 := basePath.Group("/auth")
	oauth2.POST("/sign-in", handlers.UserLogin)
	oauth2.POST("/sign-up", handlers.Register)
	oauth2.POST("/verify-account", handlers.VerifyAccount)
	oauth2.POST("/reset-password", handlers.RequestPasswordReset)
	oauth2.POST("/set-password", handlers.SetUserPassword)
}

func SessionAuth(basePath *gin.Engine) {
	basePath.GET("/login", handlers.ShowLoginPage)
	basePath.POST("/login", handlers.SessionLogin)
	successGroup := basePath.Group("/success").Use(middleware.SessionMiddleware())
	successGroup.GET("/success", handlers.ShowSuccessPage)
}
