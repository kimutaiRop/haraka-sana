package models

import (
	"time"
)

type User struct {
	Id int `gorm:"primary_key" json:"id"`

	Username string `json:"username"`
	Email    string `json:"email"`
	FrstName string `json:"firtname"`
	LastName string `json:"lastname"`
	Password string `json:"-"`

	VeriedAt time.Time `json:"veried_at"`
	Active   bool      `json:"active"`
	Location string    `json:"location"`
	Phone    string    `json:"phone"`
	Business string    `json:"business"`

	ForDeletion  bool      `json:"for_deletion"`
	DeletionTime time.Time `json:"deletion_time"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
