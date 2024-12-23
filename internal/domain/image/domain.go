package imageDomain

type Image struct {
	ID   string   `gorm:"primaryKey"`
	URLs []string `gorm:"type:text[]"`
}
