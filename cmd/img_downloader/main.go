package main

import (
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"img_downloader/gen/img_downloader/v1/img_downloaderv1connect"
	"img_downloader/internal/config"
	imageServer "img_downloader/internal/server/image"
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

	_ = storage.ConnectPostgresDB()

	imgDownloader := imageServer.NewImageServer(log)

	mux := http.NewServeMux()
	path, handler := img_downloaderv1connect.NewImageServiceHandler(imgDownloader)
	mux.Handle(path, handler)
	err := http.ListenAndServe(
		cfg.GRPC.Host+":"+strconv.Itoa(cfg.GRPC.Port),
		h2c.NewHandler(mux, &http2.Server{}),
	)

	if err != nil {
		fmt.Errorf("server error: %v", err)
	}

	log.Info(fmt.Sprintf("Server is running on %s:%d", cfg.GRPC.Host, cfg.GRPC.Port))
}
