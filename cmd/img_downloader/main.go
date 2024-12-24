package main

import (
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"img_downloader/gen/img_downloader/v1/img_downloaderv1connect"
	"img_downloader/internal/config"
	"img_downloader/internal/nats"
	imageRepo "img_downloader/internal/repository/image"
	imageServer "img_downloader/internal/server/image"
	imageService "img_downloader/internal/services/image"
	"img_downloader/internal/storage"
	"log/slog"
	"net/http"
	"strconv"
)

func main() {
	log := slog.Default()

	log.Info("Start")

	cfg := config.MustLoad()

	log.Info("Load config",
		slog.String("Env", cfg.Env),
		slog.String("Port", strconv.Itoa(cfg.GRPC.Port)))

	db := storage.ConnectPostgresDB(log)

	natsProducer, err := nats.New(&cfg.Nats, "image_urls", log)
	if err != nil {
		log.Error("Failed to create nats producer", err)
	}

	log.Info("Created nats producer on", slog.String("port", strconv.Itoa(cfg.Nats.Port)))

	imgRepo := imageRepo.New(db)
	imgService := imageService.New(log, imgRepo)
	imgServer := imageServer.New(log, imgService, natsProducer)

	mux := http.NewServeMux()
	path, handler := img_downloaderv1connect.NewImageServiceHandler(imgServer)
	mux.Handle(path, handler)
	err = http.ListenAndServe(
		cfg.GRPC.Host+":"+strconv.Itoa(cfg.GRPC.Port),
		h2c.NewHandler(mux, &http2.Server{}),
	)

	if err != nil {
		fmt.Errorf("server error: %v", err)
	}
}
