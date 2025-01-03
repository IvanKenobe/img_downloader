package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

func ConnectPostgresDB(log *slog.Logger) *gorm.DB {

	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=img_downloader_db sslmode=disable timezone=Europe/Kyiv"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	log.Info("Database connected")

	return db
}
