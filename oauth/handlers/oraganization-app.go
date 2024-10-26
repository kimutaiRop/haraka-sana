package handlers

import (
	"haraka-sana/config"
	"haraka-sana/oauth/models"
	"haraka-sana/oauth/objects"
	"haraka-sana/oauth/services"
	userModel "haraka-sana/users/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateNewApp(c *gin.Context) {
	var createApp *objects.CreateApp

	err := c.ShouldBindJSON(&createApp)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	contextUser, _ := c.Get("user")
	user := contextUser.(userModel.User)

	var orgExist models.OrganizationApplication

	config.DB.Where(&models.OrganizationApplication{UserId: user.Id, ApplicationName: createApp.ApplicationName}).
		First(&orgExist)
	if orgExist.Id != 0 {
		c.JSON(400, gin.H{"error": "you already have an application with this name"})
		return
	}
	clientId := services.GenerateRandomString(25)
	clientSecret := services.GenerateRandomString(64)
	NewApp := models.OrganizationApplication{
		ApplicationName: createApp.ApplicationName,
		Website:         createApp.Website,
		Logo:            createApp.Logo,
		RedirectURIs:    createApp.RedirectURIs,
		ClientId:        clientId,
		ClientSecret:    clientSecret,
		UserId:          user.Id,
	}

	config.DB.Create(&NewApp)

	c.JSON(http.StatusOK, gin.H{
		"client_id":    clientId,
		"clien_secret": clientSecret,
	})
}

func GetMyApps(c *gin.Context) {
	contextUser, _ := c.Get("user")
	user := contextUser.(userModel.User)

	var orgApps []models.OrganizationApplication
	config.DB.Where(&models.OrganizationApplication{UserId: user.Id}).
		Find(&orgApps)

	c.JSON(http.StatusOK, gin.H{
		"apps": orgApps,
	})
}
