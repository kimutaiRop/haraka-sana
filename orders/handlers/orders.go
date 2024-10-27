package handlers

import (
	"fmt"
	"haraka-sana/config"
	oauthModel "haraka-sana/oauth/models"
	"haraka-sana/orders/models"
	"haraka-sana/orders/objects"
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

	var orderInfo objects.CreateOrder

	err := c.ShouldBindJSON(&orderInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error persring information, check fiedls",
		})
		return
	}

	if err := orderInfo.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check order id not exist
	var orderExist models.Order

	config.DB.Where(&models.Order{
		SellerOrderId:     orderInfo.SellerOrderId,
		OrganizationAppId: organizationApp.Id,
	}).First(&orderExist)

	if orderExist.Id != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You already created order with same order id",
		})
		return
	}

	config.DB.Create(&orderInfo.Product)
	product := orderInfo.Product
	seller := orderInfo.Seller
	customer := orderInfo.Customer
	newProduct := models.Product{
		Size:  product.Size,
		Name:  product.Name,
		Image: product.Image,
	}
	err = config.DB.Create(&newProduct).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error creating order, check product fields and try again",
		})
		return
	}
	newSeller := models.Seller{
		Address:  seller.Address,
		Email:    seller.Email,
		Phone:    seller.Phone,
		FullName: seller.FullName,
		Country:  seller.Country,
		City:     seller.City,
	}
	err = config.DB.Create(&newSeller).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error creating order, check seller fields and try again",
		})
		return
	}

	newCustomer := models.Customer{
		FullName: customer.FullName,
		Country:  customer.Country,
		City:     customer.City,
		Address:  customer.Address,
		Phone:    customer.Phone,
		Email:    customer.Email,
	}
	err = config.DB.Create(&newCustomer).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error creating order, check customer fields and try again",
		})
		return
	}

	newOrder := models.Order{
		CustomerId:        newCustomer.Id,
		ProductId:         newProduct.Id,
		SellerId:          newSeller.Id,
		SellerOrderId:     orderInfo.SellerOrderId,
		OrganizationAppId: organizationApp.Id,
	}
	fmt.Println(newOrder)
	err = config.DB.Create(&newOrder).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "error creating order, check fields and try again",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "Order created successfully, use order id to track order progress",
	})
}
