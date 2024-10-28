package objects

type OrderStep struct {
	OrderId   int    `json:"order_id"`
	Country   string `json:"country"`
	Delivered bool   `json:"delivered"`
	Message   string `json:"message"`
}
