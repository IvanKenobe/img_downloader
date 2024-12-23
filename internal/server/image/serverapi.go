package imageServer

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	imgdownloaderv1 "img_downloader/gen/img_downloader/v1"
	"img_downloader/internal/services/image"
	"log/slog"
	"net/url"
)

type Server struct {
	log     *slog.Logger
	service *imageService.Service
}

func New(log *slog.Logger, service *imageService.Service) *Server {
	return &Server{
		log:     log,
		service: service,
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

	newUrls, err := s.service.CheckUrls(ctx, req.Msg.Urls)

	if err != nil {
		s.log.Info("image download failed", "urls", req.Msg.Urls)
		return nil, connect.NewError(connect.CodeInternal, errors.New("internal server error"))
	}

	for _, u := range newUrls {
		s.log.Info("download url", "url", u)
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
