package config

import (
	"fmt"
	oauthModels "haraka-sana/oauth/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

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
		&oauthModels.Organization{},
		&oauthModels.Code{},
		&oauthModels.AuthorizationToken{},
	)
	fmt.Println("Database migrated successfully")

	DB = database

}

func SeedDatabase() {
	fmt.Println("Seeding database...")
}