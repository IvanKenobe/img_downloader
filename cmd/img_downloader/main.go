package main

import (
	"fmt"
	"go.uber.org/mock/gomock"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"img_downloader/gen/img_downloader/v1/img_downloaderv1connect"
	"img_downloader/internal/config"
	natsConsumer "img_downloader/internal/nats/consumer"
	natsProducer "img_downloader/internal/nats/producer"
	imageRepo "img_downloader/internal/repository/image"
	imageServer "img_downloader/internal/server/image"
	imageService "img_downloader/internal/services/image"
	"img_downloader/internal/storage"
	"img_downloader/internal/uploader"
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

	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockUploader := uploader.NewMockUploader(ctrl)

	mockUploader.EXPECT().UploadToS3(gomock.Any()).Return("mock-s3-url", nil).AnyTimes()
	mockUploader.EXPECT().UploadToSFTP(gomock.Any()).Return("mock-sftp-url", nil).AnyTimes()

	producer := natsProducer.New(&cfg.Nats, "image_urls", log)
	consumer := natsConsumer.New(&cfg.Nats, "image_urls", log, mockUploader)

	consumer.Start()

	imgRepo := imageRepo.New(db)
	imgService := imageService.New(log, imgRepo, producer)
	imgServer := imageServer.New(log, imgService, producer)

	mux := http.NewServeMux()
	path, handler := img_downloaderv1connect.NewImageServiceHandler(imgServer)
	mux.Handle(path, handler)
	err := http.ListenAndServe(
		cfg.GRPC.Host+":"+strconv.Itoa(cfg.GRPC.Port),
		h2c.NewHandler(mux, &http2.Server{}),
	)

	if err != nil {
		fmt.Errorf("server error: %v", err)
	}
}
