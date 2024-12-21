package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgresDB() *gorm.DB {

	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=img_downloader_db sslmode=disable timezone=Europe/Kyiv"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connected to database")

	return db
}
