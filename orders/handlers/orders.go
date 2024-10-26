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
	var orders []models.Order

	var totalCount int64

	if len(filters.Filter) == 0 {
		config.DB.
			Offset(filters.Offest).
			Limit(filters.PageSize).
			Preload("Customer").
			Preload("Customer").
			Preload("Product").
			Order(filters.OrderBy).
			Find(&orders)

		config.DB.
			Model(&models.Order{}).
			Select("orders.id").
			Count(&totalCount)
	} else {
		config.DB.
			Where(clause.Where{Exprs: filters.Filter}).
			Preload("Customer").
			Preload("Customer").
			Preload("Product").
			Order(filters.OrderBy).
			Find(&orders)

		config.DB.
			Where(clause.Where{Exprs: filters.Filter}).
			Model(&models.Order{}).
			Select("orders.id").
			Count(&totalCount)
	}
	totalPages := 0

	if int(totalCount)%filters.PageSize == 0 {
		totalPages = int(totalCount) / filters.PageSize
	} else {
		totalPages = (int(totalCount) / filters.PageSize) + 1
	}

	pageInfo := gin.H{
		"page":        filters.Page,
		"page_size":   filters.PageSize,
		"total_count": totalCount,
		"total_pages": totalPages,
	}
	c.JSON(http.StatusOK, gin.H{
		"orders":    orders,
		"page_info": pageInfo,
	})
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
	m_filter = append(m_filter, clause.Eq{Column: "oraganization_app_id", Value: organization.Id})
	filters.Filter = m_filter
	var orders []models.Order

	var totalCount int64

	if len(filters.Filter) == 0 {
		config.DB.
			Offset(filters.Offest).
			Limit(filters.PageSize).
			Preload("Customer").
			Preload("Customer").
			Preload("Product").
			Order(filters.OrderBy).
			Find(&orders)

		config.DB.
			Model(&models.Order{}).
			Select("orders.id").
			Count(&totalCount)
	} else {
		config.DB.
			Where(clause.Where{Exprs: filters.Filter}).
			Preload("Customer").
			Preload("Customer").
			Preload("Product").
			Order(filters.OrderBy).
			Find(&orders)

		config.DB.
			Where(clause.Where{Exprs: filters.Filter}).
			Model(&models.Order{}).
			Select("orders.id").
			Count(&totalCount)
	}
	totalPages := 0

	if int(totalCount)%filters.PageSize == 0 {
		totalPages = int(totalCount) / filters.PageSize
	} else {
		totalPages = (int(totalCount) / filters.PageSize) + 1
	}

	pageInfo := gin.H{
		"page":        filters.Page,
		"page_size":   filters.PageSize,
		"total_count": totalCount,
		"total_pages": totalPages,
	}
	c.JSON(http.StatusOK, gin.H{
		"orders":    orders,
		"page_info": pageInfo,
	})
}
