package routes

import (
	"haraka-sana/helpers/middleware"
	"haraka-sana/oauth/handlers"

	"github.com/gin-gonic/gin"
)

func OauthRoutes(basePath *gin.RouterGroup) {
	oauth2 := basePath.Group("/oauth2")

	oauth2.POST("/token", handlers.AuthorizeToken)
	oauth2.POST("/client-credentials", handlers.ClientCredentials)

	authorizeGroup := oauth2.Group("/authorize").Use(middleware.SessionMiddleware())
	authorizeGroup.GET("", handlers.AuthorizeCode)

	organizationRoute := basePath.Group("/organization").
		Use(middleware.JWTAuthMiddleware())
	organizationRoute.GET("/apps", handlers.GetMyApps)
	organizationRoute.POST("/create-app", handlers.CreateNewApp)

}
