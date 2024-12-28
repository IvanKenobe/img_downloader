package uploader

type S3Uploader interface {
	UploadToS3(originalURL string) (string, error)
}

type SFTPUploader interface {
	UploadToSFTP(originalURL string) (string, error)
}
