package handlers

import (
	"haraka-sana/config"
	"haraka-sana/orders/models"
	"haraka-sana/orders/objects"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTrackingStep(c *gin.Context) {
	var createStep objects.CreateTrackingStep

	err := c.ShouldBindJSON(&createStep)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newStage := models.Step{
		Name:           createStep.Name,
		StreamLocation: createStep.StreamLocation,
	}

	err = config.DB.Create(&newStage).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating new step",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "New Step added to delivery chain",
	})

}
