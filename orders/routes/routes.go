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
	oauthAccess.POST("/create", handlers.OrganizationCreateOrder)
	oauthAccess.GET("", handlers.OrganizationGetOrders)
	oauthAccess.GET("/status/:order_id", handlers.OrganizationTrackOrder)

	trackingRoutes := basePath.Group("/tracking")
	trackingRoutes.Use(middleware.StaffJWTAuthMiddleware())
	trackingRoutes.POST("/record-step",
		middleware.PermissionMiddleware(config.Permissions.CREATE_STEP),
		handlers.UpdateOrderStep)
}
