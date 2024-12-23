package imageRepo

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	imageDomain "img_downloader/internal/domain/image"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CheckURLs(ctx context.Context, urls []string) ([]string, error) {
	const op = "repository.image.CheckURLs"
	var newURLs []string

	for _, url := range urls {
		var count int64
		err := r.db.Model(&imageDomain.Image{}).
			Where("id = ?", url).
			Count(&count).Error

		if err != nil {
			return newURLs, fmt.Errorf("%s: %v", op, err)
		}

		if count == 0 {
			newURLs = append(newURLs, url)
		}
	}

	return newURLs, nil
}
