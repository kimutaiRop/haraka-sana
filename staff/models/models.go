package models

import (
	permissionsModel "haraka-sana/permissions/models"
	"time"
)

type Staff struct {
	Id           int       `gorm:"primary_key" json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	FrstName     string    `json:"firtname"`
	LastName     string    `json:"lastname"`
	Password     string    `json:"-"`
	Active       bool      `json:"active"`
	VeriedAt     time.Time `json:"veried_at"`
	Location     string    `json:"location"`
	ProfileImage string    `json:"profile_image"`
	Phone        string    `json:"phone"`
	IDNumber     string    `json:"id_number"`
	EmployeeId   string    `json:"employee_id"`

	Station string `json:"station"`
	Role    string `json:"role"`

	PositionID int                       `json:"-"`
	Position   permissionsModel.Position `json:"position" gorm:"foreignKey:PositionID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
