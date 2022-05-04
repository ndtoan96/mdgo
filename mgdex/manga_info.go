package mgdex

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MangaInfo struct {
	Id         string
	Attributes struct {
		Title                        map[string]string
		AltTitles                    []map[string]string
		Description                  map[string]string
		Status                       string
		Year                         int
		availableTranslatedLanguages []string
	}
}

// GetMangaInfo returns manga information (GET manga/{id})
func GetMangaInfo(id string) (*MangaInfo, error) {
	url := fmt.Sprintf("https://api.mangadex.org/manga/%v", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error getting %v, %v", url, resp.Status)
	}

	// Deserialize manga json response to struct
	var manga struct{ Data MangaInfo }
	err = json.NewDecoder(resp.Body).Decode(&manga)
	if err != nil {
		return nil, err
	}
	return &manga.Data, nil
}

// GetAllTitles returns all titles in all languages
func (manga *MangaInfo) GetAllTitles() []map[string]string {
	allTitles := make([]map[string]string, 0)
	allTitles = append(allTitles, manga.Attributes.Title)
	allTitles = append(allTitles, manga.Attributes.AltTitles...)
	return allTitles
}

// GetTitles returns all titles in a specific language
func (manga *MangaInfo) GetTitles(lang string) []string {
	titles := make([]string, 0)
	allTitles := manga.GetAllTitles()
	for _, m := range allTitles {
		if val, exists := m[lang]; exists {
			titles = append(titles, val)
		}
	}
	return titles
}

// GetMainTitle returns the title in the main manga page
func (manga *MangaInfo) GetMainTitle() string {
	for _, val := range manga.Attributes.Title {
		return val
	}
	return ""
}
