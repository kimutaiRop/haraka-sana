package handlers

import (
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
