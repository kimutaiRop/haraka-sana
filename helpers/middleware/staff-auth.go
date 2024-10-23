package middleware

import (
	"fmt"
	"haraka-sana/config"
	"haraka-sana/helpers"
	permissionsModel "haraka-sana/permissions/models"
	staffModel "haraka-sana/staff/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StaffJWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("Authorization")
		token, err := helpers.ValidateToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Auth token is invalid",
			})
			c.Abort() // Stop further processing
			return
		}
		if token.AccountType != "staff" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Account type",
			})
			c.Abort() // Stop further processing
			return
		}
		var staff staffModel.Staff
		getErr := config.DB.Where(&staffModel.Staff{Id: token.ID}).First(&staff).Error

		if getErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Auth token is invalid",
			})
			c.Abort() // Stop further processing
			return
		}
		if !staff.Active {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Account is not active",
			})
			c.Abort() // Stop further processing
			return
		}

		c.Set("staff", staff)
		c.Next() // Continue to the next middleware or handler if no errors
	}
}

func PermissionMiddleware(permission string) gin.HandlerFunc {

	return func(c *gin.Context) {
		contextStaff, _ := c.Get("staff")

		staff := contextStaff.(staffModel.Staff)

		hasPemission := StaffHasPermission(permission, staff.Id)
		if !hasPemission {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Insufficient permission",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func StaffHasPermission(permission string, staffID int) bool {
	// Get staff
	var staff staffModel.Staff
	config.DB.Where(&staffModel.Staff{Id: staffID}).First(&staff)

	var permissionID int
	result := config.DB.Model(&permissionsModel.Permission{}).Where("name = ?", permission).Select("id").Scan(&permissionID)
	if result.Error != nil {
		return false
	}
	var staffPermission permissionsModel.PositionPermission
	fmt.Println(staff.PositionID, permissionID)
	config.DB.Where("position_id = ? AND permission_id = ? AND active = true", staff.PositionID, permissionID).First(&staffPermission)

	return staffPermission.ID != 0
}
