package models

import (
	authModel "haraka-sana/oauth/models"
	"time"
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
	Id    int    `json:"id" gorm:"primary_key"`
	Size  string `json:"size"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Order struct {
	Id                      int                                `json:"id" gorm:"primary_key"`
	SellerOrderId           string                             `json:"order_id"`
	CustomerId              int                                `json:"-"`
	Customer                *Customer                          `json:"customer" gorm:"foreignKey:CustomerId"`
	SellerId                int                                `json:"-"`
	Seller                  *Seller                            `json:"seller" gorm:"foreignKey:SellerId"`
	ProductId               int                                `json:"-"`
	Product                 *Product                           `json:"product" gorm:"foreignKey:ProductId"`
	OrganizationAppId       int                                `json:"-" gorm:"index"`
	OrganizationApplication *authModel.OrganizationApplication `json:"organization_application,omitempty" gorm:"foreignKey:OrganizationAppId"`
	Status                  string                             `json:"status"`
	Delivered               bool                               `json:"delivered"`
	DeliveredAt             time.Time                          `json:"delivered_at"`
	CreatedAt               time.Time                          `json:"created_at"`
}
