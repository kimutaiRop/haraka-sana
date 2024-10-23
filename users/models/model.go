package models

import (
	permissionsModel "haraka-sana/permissions/models"
	"time"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FrstName string `json:"firtname"`
	LastName string `json:"lastname"`
	Password string `json:"-"`

	VeriedAt     time.Time `json:"veried_at"`
	IsAdmin      bool      `json:"is_admin"`
	Active       bool      `json:"active"`
	Location     string    `json:"location"`
	DOB          string    `json:"dob"`
	Gender       string    `json:"gender"`
	JobTitle     string    `json:"jobtitle"`
	ProfileImage string    `json:"profile_image"`
	Phone        string    `json:"phone"`
	IDNumber     string    `json:"id_number"`

	ForDeletion  bool      `json:"for_deletion"`
	DeletionTime time.Time `json:"deletion_time"`

	PositionID int                       `json:"-"`
	Position   permissionsModel.Position `json:"position" gorm:"foreignKey:PositionID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
