package imageService

import (
	"context"
	"log/slog"
)

type ImageRepo interface {
	CheckURLs(ctx context.Context, urls []string) ([]string, error)
}

type Service struct {
	log  *slog.Logger
	repo ImageRepo
}

func New(log *slog.Logger, repo ImageRepo) *Service {
	return &Service{log: log, repo: repo}
}

func (s *Service) CheckUrls(ctx context.Context, urls []string) ([]string, error) {
	const op = "service.image.CheckUrls"
	s.log.With(slog.String("op", op))

	checkedUrls, err := s.repo.CheckURLs(ctx, urls)

	if err != nil {
		return nil, err
	}

	return checkedUrls, nil
}
