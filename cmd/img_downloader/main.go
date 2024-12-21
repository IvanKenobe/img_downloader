package main

import (
	"img_downloader/internal/config"
	"img_downloader/internal/storage"
	"log/slog"
	"strconv"
)

func main() {
	log := slog.Default()

	log.Info("Start")

	cfg := config.MustLoad()

	log.Info("Load config",
		slog.String("Env", cfg.Env),
		slog.String("Port", strconv.Itoa(cfg.GRPC.Port)))

	_ = storage.ConnectPostgresDB()
}
