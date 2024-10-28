package routes

import (
	"haraka-sana/helpers/middleware"
	"haraka-sana/permissions/handlers"

	"github.com/gin-gonic/gin"
)

func PermissionRoutes(basePath *gin.RouterGroup) {
	permRoutes := basePath.Group("/permissions")
	permRoutes.Use(middleware.StaffJWTAuthMiddleware())
	permRoutes.GET("", handlers.GetPermissions)
	permRoutes.GET("/positions", handlers.GetPositions)
	permRoutes.GET("/position-permissions", handlers.GetPositionPermission)
}
