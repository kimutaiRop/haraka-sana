package models

import (
	authModel "haraka-sana/oauth/models"
)

type Customer struct {
	Id       int    `json:"id" gorm:"primary_key"`
	FullName string `json:"full_name"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type Seller struct {
	Id       int    `json:"id" gorm:"primary_key"`
	FullName string `json:"full_name"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type Product struct {
	Id      int    `json:"id" gorm:"primary_key"`
	Size    string `json:"size"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	OrderId string `json:"order_id"`
}

type Order struct {
	Id                       int                                `json:"id" gorm:"primary_key"`
	CustomerId               int                                `json:"-"`
	OraganizationAppId       int                                `json:"-"`
	Customer                 Customer                           `json:"customer" gorm:"foreignKey:CustomerId"`
	SellerId                 int                                `json:"-"`
	Seller                   Seller                             `json:"seller" gorm:"foreignKey:SellerId"`
	OraganizationApplication authModel.OraganizationApplication `json:"oraganization_application" gorm:"foreignKey:OrganizationAppId"`
	Status                   string                             `json:"status"`
	Delivered                string                             `json:"delivered"`
}
