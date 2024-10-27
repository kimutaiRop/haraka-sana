package models

import (
	permissionsModel "haraka-sana/permissions/models"
	"time"
)

type Staff struct {
	Id           int       `gorm:"primary_key" json:"id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"firtname"`
	LastName     string    `json:"lastname"`
	Password     string    `json:"-"`
	Active       bool      `json:"active"`
	VerifiedAt   time.Time `json:"verified_at"`
	Country      string    `json:"country"`
	City         string    `json:"city"`
	ProfileImage string    `json:"profile_image"`
	Phone        string    `json:"phone"`
	IDNumber     string    `json:"id_number"`
	EmployeeId   string    `json:"employee_id"`

	PositionID int                       `json:"-"`
	Position   permissionsModel.Position `json:"position" gorm:"foreignKey:PositionID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
