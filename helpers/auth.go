package helpers

import (
	"fmt"
	"haraka-sana/config"
	staffModel "haraka-sana/staff/models"
	userModel "haraka-sana/users/models"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthClaims struct {
	*jwt.StandardClaims
	ID          int    `json:"id"`
	AccountType string `json:"account_type"`
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
	tokens := strings.Split(tokenString, " ")
	token, err := jwt.Parse(tokens[1], func(token *jwt.Token) (interface{}, error) {
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
		if claims["account_type"] == "user" {
			var user userModel.User
			config.DB.Where("id = ?", claims["id"]).First(&user)
			if user.Id == 0 {
				return nil, fmt.Errorf("agent not found")
			}
			// if inactive, return error
			if !user.Active {
				return nil, fmt.Errorf("account with id not active")
			}
		}

		// for staff account type make sure staff is found
		if claims["account_type"] == "staff" {
			var staff staffModel.Staff
			config.DB.Where("id = ?", claims["id"]).First(&staff)
			if staff.Id == 0 {
				return nil, fmt.Errorf("staff not found")
			}
			// if inactive, return error
			if !staff.Active {
				return nil, fmt.Errorf("account with id not active")
			}
		}
		id := claims["id"].(float64)

		return &AuthClaims{
			ID:          int(id),
			AccountType: claims["account_type"].(string),
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
		fmt.Print("email", claims["email"])
		return &VerifyClaims{
			Email: claims["email"].(string),
		}, nil
	} else {
		return nil, err
	}
}

func HashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func CheckPasswordHash(password, hash string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hash)
	bytePassword := []byte(password)

	// Check if the password's hash is equal to the hash
	// stored in the database
	err := bcrypt.CompareHashAndPassword(byteHash, bytePassword)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
