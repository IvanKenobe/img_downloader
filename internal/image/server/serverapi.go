package server

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	imgdownloaderv1 "img_downloader/gen/img_downloader/v1"
	"img_downloader/internal/image/service"
	natsProducer "img_downloader/internal/nats/producer"
	"log/slog"
)

type ImageServer struct {
	log          *slog.Logger
	service      *service.ImageService
	natsProducer *natsProducer.Producer
}

func NewImageServer(log *slog.Logger, service *service.ImageService, natsProducer *natsProducer.Producer) *ImageServer {
	return &ImageServer{
		log:          log,
		service:      service,
		natsProducer: natsProducer,
	}
}

func (s *ImageServer) DownloadImages(
	ctx context.Context,
	req *connect.Request[imgdownloaderv1.DownloadImagesRequest],
) (*connect.Response[imgdownloaderv1.DownloadImagesResponse], error) {
	if err := ctx.Err(); err != nil {
		s.log.Error("Context error", slog.String("error", err.Error()))
		return nil, connect.NewError(connect.CodeCanceled, err)
	}

	if err := s.service.ValidateURLs(ctx, req.Msg.Urls); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	filteredURLs, err := s.service.FilterNewURLs(ctx, req.Msg.Urls)
	if err != nil {
		s.log.Info("Failed to check urls", "urls", req.Msg.Urls)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	newURLsCount, publishErrors := s.service.PublishURLsToNATS(ctx, filteredURLs)

	if len(publishErrors) > 0 {
		for _, err := range publishErrors {
			s.log.Error("Failed to publish URL to NATS", slog.String("error", err.Error()))
		}
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to publish URLs to NATS"))
	}

	existingURLsCount := len(req.Msg.Urls) - len(filteredURLs)

	return connect.NewResponse(&imgdownloaderv1.DownloadImagesResponse{
		ExistingUrls: int32(existingURLsCount),
		NewUrls:      newURLsCount,
	}), nil
}
