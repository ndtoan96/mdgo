// Pakage imgdl provides functions to download one or multiple image asynchronousely.
package imgdl

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// This constant specifies timeout (in second) downloading one image.
	TIMEOUT_SEC = 30
)

// GetImageData downloads a single image, returns its extension and a ReadCloser which holds image data.
func GetImageData(url string) (string, *io.ReadCloser, error) {
	// Get url
	resp, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	if resp.StatusCode != 200 {
		return "", nil, errors.New(fmt.Sprintf("error getting %v: %v", url, resp.Status))
	}

	// Guess file extension
	var ext string
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" || contentType == "image/jpeg" {
		ext = "jpg"
	} else if strings.HasPrefix(contentType, "image") {
		ext = strings.TrimPrefix(contentType, "image/")
	} else {
		return "", nil, errors.New(fmt.Sprintf("error getting %v: not an image", url))
	}

	return ext, &resp.Body, nil
}

// DownloadImgage downloads a single image and save to path. Unexist parent folder
// will be created and image extension is guess from MIME type.
func DownloadImgage(url string, path string) error {
	ext, data, err := GetImageData(url)
	if err != nil {
		return err
	}
	defer (*data).Close()

	// Create parent folder
	parent := filepath.Dir(path)
	err = os.MkdirAll(parent, fs.ModeDir)
	if err != nil {
		return err
	}

	// Create file
	file, err := os.Create(path + "." + ext)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write to file
	_, err = io.Copy(file, *data)
	if err != nil {
		return err
	}

	return nil
}

// DownloadImages asychronously donwloads multiple images at once and write to a folder.
func DownloadImgages(urls []string, prefix string) error {
	// Spawn go routines to download images
	c := make(chan error)
	numUrls := len(urls)
	for i, url := range urls {
		go func(url string, i int) {
			err := DownloadImgage(url, fmt.Sprintf("%v%02d", prefix, i))
			c <- err
		}(url, i)
	}

	// Listen to result
	for cnt := 0; cnt < numUrls; cnt++ {
		select {
		case err := <-c:
			if err != nil {
				return err
			}
		case <-time.After(TIMEOUT_SEC * time.Second):
			return errors.New("timeout")
		}
	}
	return nil
}

// DownloadImagesZip asychronously donwloads multiple images at once and write to a zip file.
// This function will spawn a go routine for each url so caller must ensure not too many runs at
// a time since it can lead to network error.
func DownloadImagesZip(urls []string, path, prefix string) error {
	// Create parent folders
	parent := filepath.Dir(path)
	err := os.MkdirAll(parent, 0777)
	if err != nil {
		return err
	}

	// Spawn go routines
	type result struct {
		idx  int
		ext  string
		data *io.ReadCloser
		err  error
	}
	c := make(chan result)
	for i, url := range urls {
		go func(url string, i int) {
			ext, data, err := GetImageData(url)
			c <- result{idx: i, ext: ext, data: data, err: err}
		}(url, i)
	}

	// Create archive file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := zip.NewWriter(file)
	defer writer.Close()

	// Collect result
	numUrls := len(urls)
	for cnt := 0; cnt < numUrls; cnt++ {
		select {
		case res := <-c:
			if res.err != nil {
				return res.err
			}
			// Write data to file in archive
			page, err := writer.Create(fmt.Sprintf("%v%02d.%v", prefix, res.idx, res.ext))
			if err != nil {
				return err
			}
			_, err = io.Copy(page, *res.data)
			if err != nil {
				return err
			}
		case <-time.After(TIMEOUT_SEC * time.Second):
			return errors.New("timeout")
		}
	}

	return nil
}
