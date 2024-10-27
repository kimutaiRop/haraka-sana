package routes

import (
	"haraka-sana/config"
	"haraka-sana/helpers/middleware"
	"haraka-sana/staff/handlers"

	"github.com/gin-gonic/gin"
)

func StaffRoutes(basePath *gin.RouterGroup) {
	staffAuth := basePath.GET("/staff-auth")
	staffAuth.POST("/login", handlers.StaffLogin)
	staffAuth.POST("/set-password", handlers.SetPassword)

	staffRousources := basePath.GET("/staff")
	staffRousources.Use(middleware.StaffJWTAuthMiddleware())
	staffRousources.POST("/",
		middleware.PermissionMiddleware(config.Permissions.VIEW_STAFF),
		handlers.GetStaff)
	staffRousources.POST("/create",
		middleware.PermissionMiddleware(config.Permissions.CREATE_STAFF),
		handlers.CreateStaff)
	staffRousources.POST("/rest-password",
		middleware.PermissionMiddleware(config.Permissions.EDIT_STAFF),
		handlers.StaffRequestPasswordReset)
	staffRousources.POST("/update-status",
		middleware.PermissionMiddleware(config.Permissions.EDIT_STAFF),
		handlers.UpdateStaffActiveStatus)
}
