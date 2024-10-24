package handlers

import (
	oauthModel "haraka-sana/oauth/models"
	staffModels "haraka-sana/staff/models"

	"github.com/gin-gonic/gin"
)

func GetOrders(c *gin.Context) {

	contextStaff, _ := c.Get("staff")

	staff := contextStaff.(staffModels.Staff)

	if staff.Active {

	}
	return
}

func OrganizationGetOrders(c *gin.Context) {
	contextOrganization, _ := c.Get("organization")

	organization := contextOrganization.(oauthModel.AuthorizationToken)
	if organization.Id != 1 {

	}
}
