package config

import (
	"fmt"
	"log"
	"os"

	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeycompat"
)

var (
	ValkeyClient         valkey.Client
	ValkeyCompat         valkeycompat.Cmdable
	ORDER_EVENTS_CHANNEL = "order_events_relay"
)

func InitValkey() {
	// connect to valkey
	valkeyClient, valErr := valkey.NewClient(valkey.ClientOption{InitAddress: []string{os.Getenv("VALKEY_ADDRESS")}})
	ValkeyCompat = valkeycompat.NewAdapter(valkeyClient)

	if valErr != nil {
		log.Fatal(valErr)
	}
	fmt.Println("Successfully connected to Valkey!")

	ValkeyClient = valkeyClient
}
