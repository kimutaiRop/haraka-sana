package handlers

import (
	"haraka-sana/config"
	"haraka-sana/helpers"
	"haraka-sana/users/models"
	"haraka-sana/users/objects"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Register(c *gin.Context) {
	var createUser *objects.CreateUser

	if err := c.ShouldBindJSON(&createUser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var errors []gin.H
	if createUser.Email == "" {
		errors = append(errors, gin.H{
			"field": "fullname",
			"error": "Email is required",
		})
	}
	if createUser.Username == "" {
		errors = append(errors, gin.H{
			"field": "email",
			"error": "Username is required"})
	}
	if createUser.ConfirmPassword != createUser.Password {
		errors = append(errors, gin.H{
			"field": "password",
			"error": "password fields do not match",
		})
	}
	if len(errors) > 0 {
		c.JSON(400, gin.H{"errors": errors})
		return
	}

	var foundUser models.User
	config.DB.Where(&models.User{Email: createUser.Email}).
		Or(&models.User{Username: createUser.Username}).
		First(&foundUser)
	if foundUser.Id != 0 && foundUser.Email == createUser.Email {
		c.JSON(400, gin.H{"error": "User with email already exists"})
		return
	}
	if foundUser.Id != 0 && foundUser.Username == createUser.Username {
		c.JSON(400, gin.H{"error": "User with username already exists"})
		return
	}
	pass := helpers.HashAndSalt([]byte(createUser.Password))

	User := models.User{
		Email:     createUser.Email,
		Username:  createUser.Username,
		Active:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Password:  pass,
	}
	email_token, err := helpers.GenerateVerifyEmailToken(helpers.VerifyClaims{
		Email: createUser.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"sucess":  "account created successfully",
			"warning": "error generating email verification token"})
		return
	}

	config.DB.Create(&User)
	url := os.Getenv("FRONTEND_URL") + "/auth/verify-account/" + email_token
	templateData := struct {
		Name    string
		Link    string
		Company string
	}{
		Name:    createUser.Username,
		Link:    url,
		Company: os.Getenv("COMPANY_NAME"),
	}

	r := helpers.NewEmailRequest([]string{User.Email}, "Hello "+createUser.Username,
		"Activate your Account with Using link: "+url)
	if err := r.ParseTemplate("templates/emails/verify-account.html", templateData); err != nil {
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

func UserLogin(c *gin.Context) {

	var loginUser objects.Login

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	err := config.DB.Where(&models.User{Username: loginUser.Username}).
		Or(&models.User{Email: loginUser.Username}).
		First(&user).
		Error

	if err != nil {
		c.JSON(400, gin.H{"error": "invalid credentials"})
		return
	}

	if user.Id == 0 {
		c.JSON(400, gin.H{"error": "invalid credentials"})
		return
	}

	if !user.Active {
		c.JSON(400, gin.H{"error": "account with username not active"})
		return
	}

	if !helpers.CheckPasswordHash(loginUser.Password, user.Password) {
		c.JSON(400, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := helpers.GenerateToken(helpers.AuthClaims{
		ID:          user.Id,
		AccountType: "user",
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"warning": "error logging in, please try again or contact support"})
		return
	}
	c.JSON(200, gin.H{
		"token": token,
		"user":  user,
	})
}

func VerifyAccount(c *gin.Context) {
	var verifyAccount *objects.VerifyAccount

	err := c.ShouldBindJSON(&verifyAccount)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	claims, err := helpers.ValidateVerifyEmailToken(verifyAccount.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access: invalid token"})
		return
	}
	config.DB.Model(&models.User{}).Where(&models.User{Email: claims.Email}).
		Updates(&models.User{
			VeriedAt: time.Now(),
			Active:   true,
		})
	c.JSON(200, gin.H{"success": "account verified successfully"})
}

func RequestPasswordReset(c *gin.Context) {

	var requestReset *objects.RequestPasswordReset
	if err := c.ShouldBindJSON(&requestReset); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	getErr := config.DB.Where(&models.User{Username: requestReset.Email}).First(&user).Error

	if getErr != nil {
		c.JSON(http.StatusOK, gin.H{"success": "if account with emails is found email to set password is sent"})
		return
	}

	email_token, err := helpers.GenerateVerifyEmailToken(helpers.VerifyClaims{
		Email: user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": "if account with emails is found email to set password is sent",
		})
		return
	}
	url := os.Getenv("FRONTEND_URL") + "/auth/set-password/" + email_token
	templateData := struct {
		Name    string
		URL     string
		Company string
	}{
		Name:    user.Username,
		URL:     url,
		Company: os.Getenv("COMPANY_NAME"),
	}

	r := helpers.NewEmailRequest([]string{user.Email}, "Hello "+user.Username,
		"Set your password by clicking on this link"+url)

	if err := r.ParseTemplate("templates/set-password.html", templateData); err == nil {
		r.SendEmail()
	}

	c.JSON(200, gin.H{"success": "if account with emails is found email to set password is sent"})
}

func SetUserPassword(c *gin.Context) {
	var setPass *objects.SetPassword
	if err := c.ShouldBindJSON(&setPass); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token := setPass.Token
	claims, err := helpers.ValidateVerifyEmailToken(token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access: invalid token"})
		return
	}
	config.DB.Model(&models.User{}).Where(&models.User{Email: claims.Email}).
		Updates(map[string]interface{}{
			"password":       helpers.HashAndSalt([]byte(setPass.Password)),
			"email_verified": true,
		})
	c.JSON(200, gin.H{"success": "password set successfully"})
}
