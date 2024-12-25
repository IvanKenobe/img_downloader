package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 2,
}

func DownloadImage(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create requst: %v", err)
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	imgData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return imgData, nil
}
