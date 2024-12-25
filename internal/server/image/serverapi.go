package imageServer

import (
	"connectrpc.com/connect"
	"context"
	imgdownloaderv1 "img_downloader/gen/img_downloader/v1"
	natsProducer "img_downloader/internal/nats/producer"
	"img_downloader/internal/services/image"
	"log/slog"
)

type Server struct {
	log          *slog.Logger
	service      *imageService.Service
	natsProducer *natsProducer.Producer
}

func New(log *slog.Logger, service *imageService.Service, natsProducer *natsProducer.Producer) *Server {
	return &Server{
		log:          log,
		service:      service,
		natsProducer: natsProducer,
	}
}

func (s *Server) DownloadImages(
	ctx context.Context,
	req *connect.Request[imgdownloaderv1.DownloadImagesRequest],
) (*connect.Response[imgdownloaderv1.DownloadImagesResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := s.service.ValidateURLs(ctx, req.Msg.Urls); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	checkedUrls, err := s.service.FilterNewURLs(ctx, req.Msg.Urls)
	if err != nil {
		s.log.Info("Failed to check urls", "urls", req.Msg.Urls)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	newURLsCount, publishErrors := s.service.PublishURLsToNATS(ctx, checkedUrls)

	if len(publishErrors) > 0 {
		for _, err := range publishErrors {
			s.log.Error("Failed to publish URL to NATS", slog.String("error", err.Error()))
		}
	}

	existingURLsCount := len(req.Msg.Urls) - len(checkedUrls)

	return connect.NewResponse(&imgdownloaderv1.DownloadImagesResponse{
		ExistingUrls: int32(existingURLsCount),
		NewUrls:      newURLsCount,
	}), nil
}
