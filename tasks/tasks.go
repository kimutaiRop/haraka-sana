package tasks

import (
	"context"
	"fmt"
	"haraka-sana/config"
	ordersTaks "haraka-sana/orders/tasks"

	"github.com/valkey-io/valkey-go"
)

func ListenEvents() {
	fmt.Println("Subscribed to channel:", config.ORDER_EVENTS_CHANNEL)

	err := config.ValkeyClient.Receive(context.Background(),
		config.ValkeyClient.B().Subscribe().Channel(config.ORDER_EVENTS_CHANNEL).Build(),
		func(msg valkey.PubSubMessage) {
			if msg.Channel == config.ORDER_EVENTS_CHANNEL {
				fmt.Println(msg.Message)
				ordersTaks.RelayOrderEvents(msg.Message)
			}
		})
	if err != nil {
		fmt.Println("Error subscribing to channel:", err)
	}
}
