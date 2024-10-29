package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"haraka-sana/config"
	"haraka-sana/orders/models"
	"net/http"

	"github.com/valkey-io/valkey-go"
)

func RelayOrderEvents() {
	fmt.Println("Subscribed to channel:", config.ORDER_EVENTS_CHANNEL)

	err := config.ValkeyClient.Receive(context.Background(),
		config.ValkeyClient.B().Subscribe().Channel(config.ORDER_EVENTS_CHANNEL).Build(),
		func(msg valkey.PubSubMessage) {
			var orderEvent models.OrderEvent
			if err := json.Unmarshal([]byte(msg.Message), &orderEvent); err != nil {
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
			body := []byte(msg.Message)

			r, err := http.NewRequest("POST", org_url, bytes.NewBuffer(body))
			if err != nil {
				fmt.Println("Error sending order event to organization")
			}
			r.Header.Add("Content-Type", "application/json")
			client := &http.Client{}
			res, err := client.Do(r)
			if err != nil {
				panic(err)
			}

			defer res.Body.Close()

			resStatus := res.StatusCode

			if resStatus != http.StatusOK {
				//TODO:impliment saving retry for task
			}
		})
	if err != nil {
		fmt.Println("Error subscribing to channel:", err)
	}
}
