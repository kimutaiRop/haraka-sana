package routes

import (
	"haraka-sana/helpers/middleware"
	"haraka-sana/orders/handlers"

	"github.com/gin-gonic/gin"
)

func OrdersRoutes(basePath *gin.RouterGroup) {
	ordersRoutes := basePath.Group("/orders")
	ordersRoutes.Use(middleware.StaffJWTAuthMiddleware())

	ordersRoutes.GET("/", handlers.GetOrders).
		Use(middleware.PermissionMiddleware("view_orders"))

	oauthAccess := basePath.Group("/app-orders").
		Use(middleware.OrganizationTokenMiddleware())
	oauthAccess.GET("/", handlers.OrganizationGetOrders)
}
