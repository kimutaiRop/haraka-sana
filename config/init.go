package config

import (
	"fmt"
	oauthModels "haraka-sana/oauth/models"
	ordersModels "haraka-sana/orders/models"
	permissionsModel "haraka-sana/permissions/models"
	staffModels "haraka-sana/staff/models"
	userModels "haraka-sana/users/models"
	"log"
	"os"

	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeycompat"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB

	ValkeyClient         valkey.Client
	ValkeyCompat         valkeycompat.Cmdable
	ORDER_EVENTS_CHANNEL = "order_events_relay"
)

func ConnectDatabase() {
	println("Connecting to database...")
	dsn := "host=" + os.Getenv("POSTGRES_HOST") +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=" + os.Getenv("POSTGRES_PORT") +
		" sslmode=disable TimeZone=" + os.Getenv("TIME_ZONE")
	mode := os.Getenv("APP_MODE")

	loggerMode := logger.Default.LogMode(logger.Info)
	if mode != "debug" {
		loggerMode = logger.Default.LogMode(logger.Silent)
	}

	var database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: loggerMode,
	})
	if err != nil {
		panic("Failed to connect to database!")
	}
	database.AutoMigrate(
		&oauthModels.OrganizationApplication{},
		&oauthModels.Code{},
		&oauthModels.AuthorizationToken{},

		&permissionsModel.Position{},
		&permissionsModel.Permission{},
		&permissionsModel.PositionPermission{},

		&userModels.User{},

		&staffModels.Staff{},

		&ordersModels.OrderEvent{},
		&ordersModels.Customer{},
		&ordersModels.Seller{},
		&ordersModels.Product{},
		&ordersModels.Order{},
	)
	fmt.Println("Database migrated successfully")

	DB = database

}

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

func SeedDatabase() {
	SeedPermissions()
}
