package models

type Permission struct {
	ID        int    `json:"id" gorm:"primary_key"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Position struct {
	ID        int    `json:"id" gorm:"primary_key"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PositionPermission struct {
	ID           int        `json:"id" gorm:"primary_key"`
	PositionID   int        `json:"position_id"`
	PermissionID int        `json:"permission_id"`
	Active       bool       `json:"active"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
	Permission   Permission `json:"-" gorm:"foreignKey:PermissionID"`
	Position     Position   `json:"-" gorm:"foreignKey:PositionID"`
}
