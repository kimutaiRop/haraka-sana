package routes

import (
	"haraka-sana/config"
	"haraka-sana/helpers/middleware"
	"haraka-sana/orders/handlers"

	"github.com/gin-gonic/gin"
)

func OrdersRoutes(basePath *gin.RouterGroup) {
	ordersRoutes := basePath.Group("/orders")
	ordersRoutes.Use(middleware.StaffJWTAuthMiddleware())

	ordersRoutes.GET("", middleware.PermissionMiddleware(config.Permissions.VIEW_ORDERS), handlers.GetOrders)

	oauthAccess := basePath.Group("/organization-orders").
		Use(middleware.OrganizationApplicationMiddleware())
	oauthAccess.GET("", handlers.OrganizationGetOrders)
	oauthAccess.POST("/create", handlers.OrganizationCreateOrder)
}
