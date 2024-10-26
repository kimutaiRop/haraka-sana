package handlers

import (
	oauthModel "haraka-sana/oauth/models"
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
	contextOrganization, _ := c.Get("organization")

	organization := contextOrganization.(oauthModel.AuthorizationToken)
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
