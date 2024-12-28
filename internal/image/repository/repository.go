package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	imageDomain "img_downloader/internal/image/domain"
)

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) FilterNewURLs(ctx context.Context, urls []string) ([]string, error) {
	const op = "ImageRepository.FilterNewURLs"
	var newURLs []string

	for _, url := range urls {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context closed: %w", ctx.Err())
		default:
			var count int64
			err := r.db.WithContext(ctx).Model(&imageDomain.Image{}).
				Where("id = ?", url).
				Count(&count).Error

			if err != nil {
				return newURLs, fmt.Errorf("%s: %v", op, err)
			}

			if count == 0 {
				newURLs = append(newURLs, url)
			}
		}
	}

	return newURLs, nil
}
