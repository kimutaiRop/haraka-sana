package routes

import (
	"haraka-sana/helpers/middleware"
	"haraka-sana/oauth/handlers"

	"github.com/gin-gonic/gin"
)

func OauthRoutes(basePath *gin.RouterGroup) {
	oauth2 := basePath.Group("/oauth2")
	oauth2.GET("/authorize", handlers.AuthorizeCode)
	oauth2.POST("/token", handlers.AuthorizeToken)
	oauth2.POST("/client-credentials", handlers.ClientCredentials)

	organizationRoute := basePath.Group("/oraganizations").
		Use(middleware.JWTAuthMiddleware())
	organizationRoute.POST("/create-app", handlers.CreateNewApp)
}
