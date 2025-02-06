package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"haraka-sana/config"
	"haraka-sana/orders/models"
	"haraka-sana/orders/objects"
	staffModels "haraka-sana/staff/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UpdateOrderStep(c *gin.Context) {
	var orderStep objects.OrderStep
	staffContst, _ := c.Get("staff")
	staff := staffContst.(staffModels.Staff)
	err := c.ShouldBindJSON(&orderStep)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	orderEvent := models.OrderEvent{
		OrderId:   orderStep.OrderId,
		StaffId:   staff.Id,
		EventTime: time.Now(),
		Country:   orderStep.Country,
		Message:   orderStep.Message,
	}

	err = config.DB.Create(&orderEvent).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error recording the event for the order",
		})
		return
	}
	// notify the orgaization of the new order step
	go func() {
		data, err := json.Marshal(orderEvent)
		fmt.Println(string(data))
		if err != nil {
			print("Error Marshalling Data for task: ", config.ORDER_EVENTS_CHANNEL)
		}
		config.ValkeyCompat.Publish(context.TODO(), config.ORDER_EVENTS_CHANNEL, string(data))
	}()

	c.JSON(http.StatusOK, gin.H{
		"suucess": "order event recorded",
	})
}
