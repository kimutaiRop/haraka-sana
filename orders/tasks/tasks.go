package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"haraka-sana/config"
	"haraka-sana/orders/models"
	"net/http"
)

func RelayOrderEvents(Message string) {
	var orderEvent models.OrderEvent
	if err := json.Unmarshal([]byte(Message), &orderEvent); err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return
	}
	var order models.Order

	config.DB.
		Preload("OrganizationApplication").
		Where(&models.Order{
			Id: orderEvent.OrderId,
		}).First(&order)

	if order.Id == 0 {
		fmt.Println("Error getting order to relay to organization")
		return
	}

	org_url := order.OrganizationApplication.EventsCallbackUrl
	orderEvent.Order = &order
	data, err := json.Marshal(orderEvent)
	fmt.Println(string(data))
	r, err := http.NewRequest("POST", org_url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error sending order event to organization")
	}
	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		//TODO:impliment saving retry for task
		return
	}

	defer res.Body.Close()

	resStatus := res.StatusCode

	if resStatus != http.StatusOK {
		//TODO:impliment saving retry for task
	}
}
