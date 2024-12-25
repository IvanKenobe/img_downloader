package uploader

type Uploader interface {
	UploadToS3(originalURL string) (string, error)
	UploadToSFTP(originalURL string) (string, error)
}
