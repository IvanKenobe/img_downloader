package service

import (
	"context"
	"fmt"
	natsProducer "img_downloader/internal/nats/producer"
	"log/slog"
	"net/url"
	"sync"
)

type ImageRepo interface {
	FilterNewURLs(ctx context.Context, urls []string) ([]string, error)
}

type ImageService struct {
	log          *slog.Logger
	repo         ImageRepo
	natsProducer *natsProducer.Producer
}

func NewImageService(log *slog.Logger, repo ImageRepo, natsProducer *natsProducer.Producer) *ImageService {
	return &ImageService{log: log, repo: repo, natsProducer: natsProducer}
}

func (s *ImageService) FilterNewURLs(ctx context.Context, urls []string) ([]string, error) {
	const op = "ImageService.FilterNewURLs"
	s.log.With(slog.String("op", op))

	checkedUrls, err := s.repo.FilterNewURLs(ctx, urls)

	if err != nil {
		return nil, err
	}

	return checkedUrls, nil
}

func (s *ImageService) ValidateURLs(ctx context.Context, urls []string) error {
	const op = "ImageService.ValidateURLs"
	s.log.With(slog.String("op", op))

	if len(urls) == 0 {
		return fmt.Errorf("urls are required")
	}

	errChan := make(chan error, len(urls))

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, u := range urls {
		go func(url string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
				if err := s.validateSingleUrl(url); err != nil {
					errChan <- err
				}
			}
		}(u)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ImageService) PublishURLsToNATS(ctx context.Context, urls []string) (int32, []error) {
	var wg sync.WaitGroup
	var newURLsCount int32
	var errors []error
	var errChan = make(chan error, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
				err := s.natsProducer.Publish([]byte(u))
				s.log.Info("Published", slog.String("url", u))
				if err != nil {
					errChan <- err
				} else {
					errChan <- nil
				}
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

func (s *ImageService) validateSingleUrl(u string) error {
	_, err := url.ParseRequestURI(u)
	return err
}
