package models

import (
	staffMode "haraka-sana/staff/models"
	"time"
)

type Step struct {
	Id        int    `json:"id" gorm:"primary_key"`
	Name      string `json:"name"`
	StepIndex int    `json:"step_index"`
}

type OrderEvent struct {
	Id        int             `json:"id" gorm:"primary_key"`
	OrderId   int             `json:"-"`
	Order     Order           `json:"order" gorm:"foreignKey:OrderId"`
	StaffId   int             `json:"-"`
	Staff     staffMode.Staff `json:"staff" gorm:"foreignKey:StaffId"`
	StepId    int             `json:"-"`
	Step      Step            `json:"step" gorm:"foreignKey:StepId"`
	EventTime time.Time       `json:"event_time"`
}
