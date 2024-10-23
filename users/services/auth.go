package services

import (
	"fmt"
	"haraka-sana/config"
	"haraka-sana/users/models"
	"log"
	"os"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthClaims struct {
	*jwt.StandardClaims
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type VerifyClaims struct {
	*jwt.StandardClaims
	Email string `json:"email"`
}

func GenerateToken(claims AuthClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func ValidateToken(tokenString string) (*AuthClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// for agent account type make sure agent is found
		var user models.User
		config.DB.Where("id = ?", claims["id"]).First(&user)
		if user.Id == 0 {
			return nil, fmt.Errorf("agent not found")
		}
		// if inactive, return error
		if !user.Active {
			return nil, fmt.Errorf("account with id not active")
		}

		id := claims["id"].(float64)

		return &AuthClaims{
			ID:    int(id),
			Email: user.Email,
		}, nil
	} else {
		return nil, err
	}
}

func GenerateVerifyEmailToken(claims VerifyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func ValidateVerifyEmailToken(tokenString string) (*VerifyClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &VerifyClaims{
			Email: claims["email"].(string),
		}, nil
	} else {
		return nil, err
	}
}

func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func CheckPasswordHash(password, hash string) bool {
	byteHash := []byte(hash)
	bytePassword := []byte(password)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePassword)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
