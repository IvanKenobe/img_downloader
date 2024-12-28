package handler

import (
	"context"
	"github.com/nats-io/nats.go"
	"img_downloader/internal/uploader"
	"img_downloader/pkg/utils"
	"log/slog"
)

type ImageHandler struct {
	log          *slog.Logger
	s3Uploader   uploader.S3Uploader
	sftpUploader uploader.SFTPUploader
}

func NewImageHandler(log *slog.Logger, s3Uploader uploader.S3Uploader, sftpUploader uploader.SFTPUploader) *ImageHandler {
	return &ImageHandler{
		log:          log,
		s3Uploader:   s3Uploader,
		sftpUploader: sftpUploader,
	}
}

func (h *ImageHandler) Process(ctx context.Context, msg *nats.Msg) error {
	const op = "ImageHandler.ProcessImageMessage"
	log := h.log.With(slog.String("op", op))
	url := string(msg.Data)

	_, err := utils.DownloadImage(url)
	if err != nil {
		log.Error("Failed to download image",
			slog.String("url", url),
			slog.String("error", err.Error()))
		return err
	}

	s3URL, err := h.s3Uploader.UploadToS3(url)
	if err != nil {
		log.Error("Failed to upload image to S3",
			slog.String("url", url),
			slog.String("error", err.Error()))
		return err
	}

	sftpURL, err := h.sftpUploader.UploadToSFTP(url)
	if err != nil {
		log.Error("Failed to upload image to SFTP",
			slog.String("url", url),
			slog.String("error", err.Error()))
		return err
	}

	log.Info("Image processed successfully",
		slog.String("originalURL", url),
		slog.String("s3URL", s3URL),
		slog.String("sftpURL", sftpURL),
	)

	return nil
}
