package routes

import (
	"haraka-sana/config"
	"haraka-sana/helpers/middleware"
	"haraka-sana/staff/handlers"

	"github.com/gin-gonic/gin"
)

func StaffRoutes(basePath *gin.RouterGroup) {

	staffRousources := basePath.Group("/staff")
	staffRousources.Use(middleware.StaffJWTAuthMiddleware())
	staffRousources.GET("",
		middleware.PermissionMiddleware(config.Permissions.VIEW_STAFF),
		handlers.GetStaff)

	staffRousources.POST("/create",
		middleware.PermissionMiddleware(config.Permissions.CREATE_STAFF),
		handlers.CreateNewStaff)
	staffRousources.POST("/reset-password",
		middleware.PermissionMiddleware(config.Permissions.EDIT_STAFF),
		handlers.StaffRequestPasswordReset)
	staffRousources.POST("/update-status",
		middleware.PermissionMiddleware(config.Permissions.EDIT_STAFF),
		handlers.UpdateStaffActiveStatus)

	staffAuth := basePath.Group("/staff/auth")
	staffAuth.POST("/create-admin", handlers.CreateAdmin)
	staffAuth.POST("/login", handlers.StaffLogin)
	staffAuth.POST("/set-password", handlers.SetPassword)
}
