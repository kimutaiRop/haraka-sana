package models

import (
	staffMode "haraka-sana/staff/models"
	"time"
)

type Batch struct {
	Id            int              `json:"id" gorm:"primary_key"`
	VehicleNumber string           `json:"vehicle_number"`
	VehicleType   string           `json:"vehicle_type"`
	StartCountry  string           `json:"start_country"`
	StopCountry   string           `json:"stop_country"`
	Status        string           `json:"status"`
	StartLocation string           `json:"start_location"`
	StopLocation  string           `json:"stop_location"`
	TrackingInfo  string           `json:"tracking_info"`
	OpenTracking  bool             `json:"open_tracking"` //customers can track
	StaffId       int              `json:"-"`
	Staff         *staffMode.Staff `json:"staff" gorm:"foreignKey:StaffId"`

	DepartureTime time.Time `json:"departure_time"`
	ArrivalTime   time.Time `json:"arrival_time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type BatchOrder struct {
	Id        int       `json:"id" gorm:"primary_key"`
	BatchId   int       `json:"batch_id"`
	Batch     *Batch    `json:"batch" gorm:"foreignKey:BatchId"`
	OrderId   int       `json:"order_id"`
	Order     *Order    `json:"order" gorm:"foreignKey:order_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
