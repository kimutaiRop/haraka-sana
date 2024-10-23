package handlers

import (
	"haraka-sana/config"
	"haraka-sana/helpers"
	"haraka-sana/users/models"
	"haraka-sana/users/objects"
	"haraka-sana/users/services"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Register(c *gin.Context) {
	var createUser objects.CreateUser

	if err := c.ShouldBindJSON(createUser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var errors []gin.H
	if createUser.Email == "" {
		errors = append(errors, gin.H{
			"field": "fullname",
			"error": "username is required"})
	}
	if createUser.Username == "" {
		errors = append(errors, gin.H{
			"field": "email",
			"error": "email is required"})
	}
	if createUser.ConfirmPassword != createUser.Password {
		errors = append(errors, gin.H{
			"field": "password",
			"error": "password fields do not match"})
	}
	if len(errors) > 0 {
		c.JSON(400, gin.H{"errors": errors})
		return
	}

	var foundUser models.User
	config.DB.Where("email = ?", createUser.Email).First(&foundUser)
	if foundUser.Id != 0 {
		c.JSON(400, gin.H{"error": "User with email already exists"})
		return
	}
	config.DB.Where("username = ?", createUser.Username).First(&foundUser)
	if foundUser.Id != 0 {
		c.JSON(400, gin.H{"error": "User with username already exists"})
		return
	}
	pass := services.HashAndSalt([]byte(createUser.Password))

	User := models.User{
		Email:     createUser.Email,
		Username:  createUser.Username,
		Active:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Password:  pass,
	}
	email_token, err := services.GenerateVerifyEmailToken(services.VerifyClaims{
		Email: createUser.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"sucess":  "account created successfully",
			"warning": "error generating email verification token"})
		return
	}

	config.DB.Create(&User)
	templateData := struct {
		Name    string
		Link    string
		Company string
	}{
		Name:    createUser.Username,
		Link:    os.Getenv("FRONTEND_URL") + "/auth/set-password/" + email_token,
		Company: os.Getenv("COMPANY_NAME"),
	}

	r := helpers.NewRequest([]string{User.Email}, "Hello "+createUser.Username, "Activate your Account with Using link: "+email_token)
	if err := r.ParseTemplate("templates/emails/set-pass-email.html", templateData); err != nil {
		c.JSON(500, gin.H{
			"warning": "error sending email",
			"success": "account created successfully",
		})
		return
	}
	sent, err := r.SendEmail()
	if err != nil || !sent {
		c.JSON(500, gin.H{
			"warning": "error sending email",
			"success": "account created successfully",
		})
	}
	c.JSON(200, gin.H{
		"success": "account created successfully, email to set password sent",
	})

}
