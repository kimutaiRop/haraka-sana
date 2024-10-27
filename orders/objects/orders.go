package objects

import "errors"

type Customer struct {
	FullName string `json:"full_name"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type Seller struct {
	FullName string `json:"full_name"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type Product struct {
	Size  string `json:"size"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type CreateOrder struct {
	Customer      Customer `json:"customer"`
	Seller        Seller   `json:"seller"`
	Product       Product  `json:"product"`
	SellerOrderId string   `json:"order_id"`
}

func (o *CreateOrder) Validate() error {
	if o.Customer.FullName == "" {
		return errors.New("customer full name is required")
	}
	if o.Customer.Country == "" {
		return errors.New("customer country is required")
	}
	if o.Customer.City == "" {
		return errors.New("customer city is required")
	}
	if o.Customer.Address == "" {
		return errors.New("customer address is required")
	}
	if o.Customer.Phone == "" {
		return errors.New("customer phone is required")
	}
	if o.Customer.Email == "" {
		return errors.New("customer email is required")
	}
	if o.Seller.FullName == "" {
		return errors.New("seller full name is required")
	}
	if o.Seller.Country == "" {
		return errors.New("seller country is required")
	}
	if o.Seller.City == "" {
		return errors.New("seller city is required")
	}
	if o.Seller.Address == "" {
		return errors.New("seller address is required")
	}
	if o.Seller.Phone == "" {
		return errors.New("seller phone is required")
	}
	if o.Seller.Email == "" {
		return errors.New("seller email is required")
	}
	if o.Product.Name == "" {
		return errors.New("product name is required")
	}
	if o.Product.Size == "" {
		return errors.New("product size is required")
	}
	if o.Product.Image == "" {
		return errors.New("product image is required")
	}
	if o.SellerOrderId == "" {
		return errors.New("seller order ID is required")
	}
	return nil
}
