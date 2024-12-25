package imageServer

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	imgdownloaderv1 "img_downloader/gen/img_downloader/v1"
	natsProducer "img_downloader/internal/nats/producer"
	"img_downloader/internal/services/image"
	"log/slog"
	"net/url"
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

	if err := validateUrls(req); err != nil {
		return nil, err
	}

	checkedUrls, err := s.service.CheckUrls(ctx, req.Msg.Urls)
	if err != nil {
		s.log.Info("Failed to check urls", "urls", req.Msg.Urls)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	for _, u := range checkedUrls {
		err = s.natsProducer.Publish([]byte(u))

		if err != nil {
			s.log.Error("nats publish failed",
				slog.String("url", u),
				slog.String("err", err.Error()))
		}
	}

	return connect.NewResponse(&imgdownloaderv1.DownloadImagesResponse{
		ExistingUrls: 0,
		NewUrls:      0,
	}), nil
}

func validateUrls(req *connect.Request[imgdownloaderv1.DownloadImagesRequest]) error {
	if len(req.Msg.GetUrls()) == 0 {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("urls are required"))
	}

	for _, u := range req.Msg.GetUrls() {
		_, err := url.ParseRequestURI(u)

		if err != nil {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}
	}

	return nil
}
