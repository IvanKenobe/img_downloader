package imageService

import (
	"context"
	natsProducer "img_downloader/internal/nats/producer"
	"log/slog"
	"net/url"
	"sync"
)

type ImageRepo interface {
	FilterNewURLs(ctx context.Context, urls []string) ([]string, error)
}

type Service struct {
	log          *slog.Logger
	repo         ImageRepo
	natsProducer *natsProducer.Producer
}

func New(log *slog.Logger, repo ImageRepo, natsProducer *natsProducer.Producer) *Service {
	return &Service{log: log, repo: repo, natsProducer: natsProducer}
}

func (s *Service) FilterNewURLs(ctx context.Context, urls []string) ([]string, error) {
	const op = "service.image.CheckUrls"
	s.log.With(slog.String("op", op))

	checkedUrls, err := s.repo.FilterNewURLs(ctx, urls)

	if err != nil {
		return nil, err
	}

	return checkedUrls, nil
}

func (s *Service) ValidateURLs(ctx context.Context, urls []string) error {
	const op = "service.image.ValidateURLs"
	s.log.With(slog.String("op", op))

	errChan := make(chan error, len(urls))

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, u := range urls {
		go func(url string) {
			defer wg.Done()
			if err := s.validateSingleUrl(url); err != nil {
				errChan <- err
			}
		}(u)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) PublishURLsToNATS(ctx context.Context, urls []string) (int32, []error) {
	var wg sync.WaitGroup
	var newURLsCount int32
	var errors []error

	var errChan = make(chan error, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			err := s.natsProducer.Publish([]byte(u))
			if err != nil {
				errChan <- err
			} else {
				errChan <- nil
			}
		}(url)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			errors = append(errors, err)
		} else {
			newURLsCount++
		}
	}

	return newURLsCount, errors
}

func (s *Service) validateSingleUrl(u string) error {
	_, err := url.Parse(u)
	return err
}
