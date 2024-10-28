package handlers

import (
	"haraka-sana/config"
	"haraka-sana/permissions/models"
	"haraka-sana/permissions/objects"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPositions(c *gin.Context) {
	var positions []models.Position

	config.DB.Find(&positions)
	c.JSON(http.StatusOK, gin.H{
		"positions": positions,
	})
}

func GetPermissions(c *gin.Context) {
	var permissons []models.Permission
	config.DB.Find(&permissons)

	c.JSON(http.StatusOK, gin.H{
		"permissons": permissons,
	})
}

func GetPositionPermission(c *gin.Context) {
	var positionPermissions []models.PositionPermission

	config.DB.Find(&positionPermissions)

	c.JSON(http.StatusOK, gin.H{
		"position_Permissions": positionPermissions,
	})
}

func CreatePosition(c *gin.Context) {

	var createPosition objects.CreatePosition

	err := c.ShouldBindJSON(&createPosition)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if createPosition.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Position Name Required",
		})
		return
	}
	newPosition := models.Position{
		Name: createPosition.Name,
	}
	config.DB.Where(&models.Position{
		Name: newPosition.Name,
	}).FirstOrCreate(&newPosition)

	c.JSON(http.StatusOK, gin.H{
		"sucess":   "position created",
		"position": newPosition,
	})
}

func CreatePositionPermission(c *gin.Context) {

	var positionPermission objects.CreatePositionPermission

	newPermissions := models.PositionPermission{
		PositionID:   positionPermission.Position,
		PermissionID: positionPermission.Permission,
	}
	config.DB.Where(&models.PositionPermission{
		PositionID:   newPermissions.PositionID,
		PermissionID: positionPermission.Permission,
	}).FirstOrCreate(&newPermissions)

	newPermissions.Active = positionPermission.Active
	config.DB.Save(newPermissions)

	c.JSON(http.StatusOK, gin.H{
		"sucess":              "position created",
		"position_permission": newPermissions,
	})
}
