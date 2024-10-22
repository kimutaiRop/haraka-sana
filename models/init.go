package models

import (
	"fmt"
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
	loggerMode := logger.Default.LogMode(logger.Info)
	if os.Getenv("DEBUG") != "1" {
		loggerMode = logger.Default.LogMode(logger.Silent)
	}

	var database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: loggerMode,
	})
	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = database

}

func SeedDatabase() {
	fmt.Println("Seeding database...")
}
