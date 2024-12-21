package imageServer

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	imgdownloaderv1 "img_downloader/gen/img_downloader/v1"
	"log/slog"
	"net/url"
)

type imageServer struct {
	log *slog.Logger
}

func NewImageServer(log *slog.Logger) *imageServer {
	return &imageServer{
		log: log,
	}
}

func (s *imageServer) DownloadImages(
	ctx context.Context,
	req *connect.Request[imgdownloaderv1.DownloadImagesRequest],
) (*connect.Response[imgdownloaderv1.DownloadImagesResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := validateUrls(req); err != nil {
		return nil, err
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
