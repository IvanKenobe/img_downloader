package main

import (
	"context"
	"fmt"
	"go.uber.org/mock/gomock"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"img_downloader/gen/img_downloader/v1/img_downloaderv1connect"
	"img_downloader/internal/config"
	imageHandler "img_downloader/internal/image/handler"
	imageRepo "img_downloader/internal/image/repository"
	imageServer "img_downloader/internal/image/server"
	imageService "img_downloader/internal/image/service"
	natsConsumer "img_downloader/internal/nats/consumer"
	natsProducer "img_downloader/internal/nats/producer"
	"img_downloader/internal/storage"
	"img_downloader/internal/uploader"
	"log/slog"
	"net/http"
	"strconv"
)

func main() {

	ctx := context.Background()
	log := slog.Default()

	log.Info("Start")

	// Load config
	cfg := config.MustLoad()
	log.Info("Load config",
		slog.String("Env", cfg.Env),
		slog.String("Port", strconv.Itoa(cfg.GRPC.Port)))

	// Connect DB
	db := storage.ConnectPostgresDB(log)

	// Mocks for uploaders
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()
	mockS3Uploader := uploader.NewMockS3Uploader(ctrl)
	mockSFTPUploader := uploader.NewMockSFTPUploader(ctrl)
	mockS3Uploader.EXPECT().UploadToS3(gomock.Any()).Return("mock-s3-url", nil).AnyTimes()
	mockSFTPUploader.EXPECT().UploadToSFTP(gomock.Any()).Return("mock-sftp-url", nil).AnyTimes()

	// Initialize NATS Producer
	producer := natsProducer.New(&cfg.Nats, "image_urls", log)

	// Initialize IMAGE repo, service, server and handler
	imgRepo := imageRepo.NewImageRepository(db)
	imgService := imageService.NewImageService(log, imgRepo, producer)
	imgServer := imageServer.NewImageServer(log, imgService, producer)
	imgHandler := imageHandler.NewImageHandler(log, mockS3Uploader, mockSFTPUploader)

	// Start NATS Consumers
	consumer, _ := natsConsumer.New(&cfg.Nats, log)
	consumer.RegisterHandler("image_urls", imgHandler)
	consumer.Start(ctx)

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
