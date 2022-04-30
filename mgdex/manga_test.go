package mgdex

import (
	"math"
	"strconv"
	"testing"
)

func Test_Manga(t *testing.T) {
	var manga *MangaData
	var err error

	manga, err = MangaQuery("d7037b2a-874a-4360-8a7b-07f2899152fd").IncludeScanlationGroup().Language("vi").Limit(5).Order("asc").GetManga()
	if err != nil {
		t.Fatal(err)
	}
	chapterList := manga.Filter().GetChapters()
	if len(chapterList) != 5 {
		t.FailNow()
	}

	var old_chapter float64
	old_chapter = -1.0
	for _, chapter := range chapterList {
		if chapter.GetLanguage() != "vi" || chapter.GetScanlationGroup() == "" {
			t.FailNow()
		}
		new_chapter, _ := strconv.ParseFloat(chapter.GetChapter(), 64)
		if new_chapter < old_chapter {
			t.FailNow()
		}
	}

	manga, err = MangaQuery("d7037b2a-874a-4360-8a7b-07f2899152fd").Limit(5).Order("desc").GetManga()
	old_chapter = math.Inf(1)
	for _, chapter := range chapterList {
		new_chapter, _ := strconv.ParseFloat(chapter.GetChapter(), 64)
		if new_chapter > old_chapter {
			t.FailNow()
		}
	}
}

func Test_Manga_Fail(t *testing.T) {
	var err error
	// id
	_, err = MangaQuery("rubbish_id").GetManga()
	if err == nil {
		t.FailNow()
	}

	// negative limit
	_, err = MangaQuery("d7037b2a-874a-4360-8a7b-07f2899152fd").Limit(-1).GetManga()
	if err == nil {
		t.FailNow()
	}

	// limit too big
	_, err = MangaQuery("d7037b2a-874a-4360-8a7b-07f2899152fd").Limit(501).GetManga()
	if err == nil {
		t.FailNow()
	}

	// order
	_, err = MangaQuery("d7037b2a-874a-4360-8a7b-07f2899152fd").Order("abcd").GetManga()
	if err == nil {
		t.FailNow()
	}
}

func Test_Manga_Length_and_Append(t *testing.T) {
	m1 := MangaData{}
	if m1.Length() != 0 {
		t.FailNow()
	}
	c1 := ChapterData{}
	c2 := ChapterData{}
	c3 := ChapterData{}
	m2 := MangaData{Data: []ChapterData{c1, c2, c3}}
	if m2.Length() != 3 {
		t.FailNow()
	}
	m1.Append(&m2)
	if m1.Length() != 3 {
		t.FailNow()
	}
}
