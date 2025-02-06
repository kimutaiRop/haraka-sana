package models

import (
	staffMode "haraka-sana/staff/models"
	"time"
)

type OrderEvent struct {
	Id        int              `json:"id" gorm:"primary_key"`
	OrderId   int              `json:"order_id"`
	Order     *Order           `json:"order" gorm:"foreignKey:OrderId"`
	StaffId   int              `json:"-"`
	Staff     *staffMode.Staff `json:"staff,omitempty" gorm:"foreignKey:StaffId"`
	EventTime time.Time        `json:"event_time"`
	Message   string           `json:"message"`
	Country   string           `json:"country"`
	Delivered bool             `json:"-"`
}

type LiveCoordinate struct {
	Id        int       `json:"id" gorm:"primary_key"`
	BatchId   int       `json:"batch_id"`
	Batch     Batch     `json:"batch" gorm:"foreignKey:BatchId"`
	Latitude  float64   `json:"latitude"`  // Latitude coordinate
	Longitude float64   `json:"longitude"` // Longitude coordinate
	Timestamp time.Time `json:"timestamp"` // Time of location update
}
