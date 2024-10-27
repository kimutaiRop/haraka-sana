package handlers

import (
	"haraka-sana/config"
	oauthModel "haraka-sana/oauth/models"
	"haraka-sana/orders/models"
	"haraka-sana/orders/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetOrders(c *gin.Context) {
	filters := services.OrderFilters(c)
	m_filter := filters.Filter
	filters.Filter = m_filter

	c.JSON(http.StatusOK, filters.GetOrders())
}

func OrganizationGetOrders(c *gin.Context) {
	contextOrganization, _ := c.Get("orgazation_app")

	organization := contextOrganization.(oauthModel.OrganizationApplication)
	if organization.Id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authorized",
		})
		return
	}
	filters := services.OrderFilters(c)
	m_filter := filters.Filter
	m_filter = append(m_filter, clause.Eq{
		Column: "organization_app_id",
		Value:  organization.Id},
	)
	filters.Filter = m_filter
	c.JSON(http.StatusOK, filters.GetOrders())
}

func OrganizationCreateOrder(c *gin.Context) {
	contextOrganization, _ := c.Get("orgazation_app")

	organizationApp := contextOrganization.(oauthModel.OrganizationApplication)
	if organizationApp.Id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authorized",
		})
		return
	}

	var orderInfo models.Order

	err := c.ShouldBindJSON(&orderInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error persring information, check fiedls",
		})
		return
	}

	config.DB.Create(&orderInfo)
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": "Order created successfully, user order id to track order progress",
		"order":   orderInfo,
	})
}
