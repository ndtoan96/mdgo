package mgdex

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// MangaData represents data of manga gotten from manga feed api of mangadex. It does not
// include information of author, artist,... since it only serves the purpose of getting chapters data
type MangaData struct{ Data []ChapterData }

type mangaQuery struct {
	id           string
	lang         string
	limit        int
	offset       int
	order        string
	includeGroup bool
}

// MangaQuery returns a query params builder
func MangaQuery(id string) *mangaQuery {
	return &mangaQuery{
		id:           id,
		lang:         "en",
		limit:        100,
		offset:       0,
		order:        "asc",
		includeGroup: false,
	}
}

// Language specifies language of manga, default is "en"
func (q mangaQuery) Language(lang string) *mangaQuery {
	q.lang = lang
	return &q
}

// Limit specifies maximum number of chapter data gotten from one query, default is 100
func (q mangaQuery) Limit(limit int) *mangaQuery {
	q.limit = limit
	return &q
}

// Offset specifies offset value of chapter data gotten from the query, default is 0
func (q mangaQuery) Offset(offset int) *mangaQuery {
	q.offset = offset
	return &q
}

// Order specifies order of chapter data (by chapter number) gotten from the query, only
// accepts two values "asc" and "desc", default is "asc".
func (q mangaQuery) Order(order string) *mangaQuery {
	q.order = order
	return &q
}

// IncludeScanlationGroup enabled translation_group data gotten from the query.
func (q mangaQuery) IncludeScanlationGroup() *mangaQuery {
	q.includeGroup = true
	return &q
}

// Verify checks values of query params.
func (q mangaQuery) Verify() error {
	if q.lang == "" {
		return errors.New("language is empty")
	}
	if q.limit < 1 || q.limit > 500 {
		return errors.New("limit is not in range [1..500]")
	}
	if q.offset < 0 {
		return errors.New("offset is negative")
	}
	if q.order != "asc" && q.order != "desc" {
		return errors.New(`expect order to be "asc" or "desc", found "` + q.order + `"`)
	}
	return nil
}

// GetManga returns manga data gotten from manga feed api of mangadex
func (q mangaQuery) GetManga() (*MangaData, error) {
	err := q.Verify()
	if err != nil {
		return nil, err
	}

	base, err := url.Parse(fmt.Sprintf("https://api.mangadex.org/manga/%v/feed", q.id))
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("translatedLanguage[]", q.lang)
	params.Add("limit", fmt.Sprint(q.limit))
	params.Add("offset", fmt.Sprint(q.offset))
	params.Add("order[chapter]", q.order)
	if q.includeGroup {
		params.Add("includes[]", "scanlation_group")
	}
	base.RawQuery = params.Encode()

	resp, err := http.Get(base.String())
	if err != nil {
		return nil, fmt.Errorf("error getting %v, %v", base, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error getting %v, %v", base, resp.Status)
	}

	var manga MangaData
	err = json.NewDecoder(resp.Body).Decode(&manga)
	if err != nil {
		return nil, err
	}

	return &manga, nil
}

// Length returns number of chapters in manga data
func (manga MangaData) Length() int {
	return len(manga.Data)
}

// Append adds chapter data from m2 to m1
func (m1 *MangaData) Append(m2 *MangaData) {
	m1.Data = append(m1.Data, m2.Data...)
}
