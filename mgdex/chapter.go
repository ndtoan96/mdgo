// Package mgdex provides interfaces to get information as well as download chapters and manga from mangadex.
package mgdex

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/ndtoan96/mdgo/mgdex/imgdl"
)

// A ChapterData represents data of chapter gotten from mangadex api. It does not
// include all possible information, only the ones commonly used.

type ChapterData struct {
	Id         string
	Attributes struct {
		Volume             string
		Chapter            string
		Title              string
		TranslatedLanguage string
		Pages              uint
	}
	Relationships []map[string]interface{}
}

type serverData struct {
	BaseUrl string
	Chapter struct {
		Hash      string
		Data      []string
		DataSaver []string
	}
}

// GetChapter send request to mangadex api and get back chapter data.
func GetChapter(id string) (*ChapterData, error) {
	// Request chapter via api
	url := fmt.Sprintf("https://api.mangadex.org/chapter/%v", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error getting %v, %v", url, resp.Status)
	}

	// Deserialize chapter json response to struct
	var chapter struct{ Data ChapterData }
	err = json.NewDecoder(resp.Body).Decode(&chapter)
	if err != nil {
		return nil, err
	}
	return &chapter.Data, nil
}

// GetId returns chapter id
func (chapter ChapterData) GetId() string {
	return chapter.Id
}

// GetVolume returns volume number of chapter, default is empty string.
func (chapter ChapterData) GetVolume() string {
	return chapter.Attributes.Volume
}

// GetChapter returns chapter number of chapter, default is empty string.
func (chapter ChapterData) GetChapter() string {
	return chapter.Attributes.Chapter
}

// GetTitle returns title of chapter, default is empty string.
func (chapter ChapterData) GetTitle() string {
	return chapter.Attributes.Title
}

// GetLanguage returns language of chapter
func (chapter ChapterData) GetLanguage() string {
	return chapter.Attributes.TranslatedLanguage
}

// GetPages returns number of pages in the chapter
func (chapter ChapterData) GetPages() uint {
	return chapter.Attributes.Pages
}

// GetScanlationGroup returns scanlation group of chapter. Note that this requires
// an additional parameter in chapter request and the function GetChapter does not
// implements it. So this function is only useful with ChapterData gotten from
// manga query where includeGroup is enabled.
func (chapter ChapterData) GetScanlationGroup() string {
	for _, rel := range chapter.Relationships {
		if rel["type"].(string) == "scanlation_group" {
			return rel["attributes"].(map[string]interface{})["name"].(string)
		}
	}
	return ""
}

// GetPageUrls returns urls for all pages in the chapter.
func (chapter ChapterData) GetPageUrls(dataSaver bool) ([]string, error) {
	// Get base url
	serverUrl := fmt.Sprintf("https://api.mangadex.org/at-home/server/%v", chapter.Id)
	resp, err := http.Get(serverUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error getting %v: %v", serverUrl, resp.Status)
	}
	var server serverData
	err = json.NewDecoder(resp.Body).Decode(&server)
	if err != nil {
		return nil, err
	}

	// Construct images urls
	var quality string
	var database []string
	if dataSaver {
		quality = "data-saver"
		database = server.Chapter.DataSaver
	} else {
		quality = "data"
		database = server.Chapter.Data
	}
	var urls []string
	for _, img := range database {
		urls = append(urls, fmt.Sprintf("%v/%v/%v/%v", server.BaseUrl, quality, server.Chapter.Hash, img))
	}
	return urls, nil
}

// Download downloads the chapter and save to folder specified by 'path'.
// If 'path' is empty, current folder will be used.
func (chapter ChapterData) Download(dataSaver bool, path string) error {
	// Get urls
	urls, err := chapter.GetPageUrls(dataSaver)
	if err != nil {
		return err
	}

	if len(urls) == 0 {
		return errors.New("Chapter " + chapter.GetChapter() + " is empty")
	}

	// Download images
	err = imgdl.DownloadImgages(urls, filepath.Join(path, "page_"))
	if err != nil {
		return err
	}
	return nil
}

// DownloadAsZip downloads the chapter and save to zip file specified by 'path'.
func (chapter ChapterData) DownloadAsZip(dataSaver bool, path string) error {
	// Get page urls
	urls, err := chapter.GetPageUrls(dataSaver)
	if err != nil {
		return err
	}

	if len(urls) == 0 {
		return errors.New("Chapter " + chapter.GetChapter() + " is empty")
	}

	// Download images to zip
	err = imgdl.DownloadImagesZip(urls, path, "page_")
	if err != nil {
		return err
	}

	return nil
}
