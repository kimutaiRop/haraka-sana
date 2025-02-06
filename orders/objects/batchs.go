package objects

type CreateBatch struct {
	VehicleNumber string `json:"vehicle_number"`
	VehicleType   string `json:"vehicle_type"`
	StartCountry  string `json:"start_country"`
	Status        string `json:"status"`
	StartLocation string `json:"start_location"`
	OpenTracking  bool   `json:"open_tracking"` //customers can track
}

type AddOrderBatch struct {
	BatchId int `json:"batch_id"`
	OrderId int `json:"order_id"`
}
