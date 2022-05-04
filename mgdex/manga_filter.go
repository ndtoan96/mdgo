package mgdex

import (
	"strconv"
	"strings"
)

type ChapterList []*ChapterData

// Filter criterias of chapters in a manga. These criterias act in an AND manner.
type mangaFilter struct {
	manga        *MangaData
	volumes      map[string]struct{}
	chapters     map[string]struct{}
	volumeRange  *[2]float64
	chapterRange *[2]float64
	preferGroups map[string]int
}

// Filter returns a pointer to mangaFilter with default values. It works
// in a builder manner.
func (manga MangaData) Filter() *mangaFilter {
	return &mangaFilter{
		manga:        &manga,
		volumes:      nil,
		chapters:     nil,
		volumeRange:  nil,
		chapterRange: nil,
		preferGroups: nil,
	}
}

// Volumes specifies list of volumes.
func (filter mangaFilter) Volumes(vols []string) *mangaFilter {
	filter.volumes = make(map[string]struct{})
	for _, vol := range vols {
		filter.volumes[vol] = struct{}{}
	}
	return &filter
}

// Chapters specifies list of chapters.
func (filter mangaFilter) Chapters(chaps []string) *mangaFilter {
	filter.chapters = make(map[string]struct{})
	for _, chap := range chaps {
		filter.chapters[chap] = struct{}{}
	}
	return &filter
}

// VolumeRange specifies an inclusive range of volume
func (filter mangaFilter) VolumeRange(min float64, max float64) *mangaFilter {
	filter.volumeRange = &[2]float64{min, max}
	return &filter
}

// ChapterRange specifies an inclusive range of chapter
func (filter mangaFilter) ChapterRange(min float64, max float64) *mangaFilter {
	filter.chapterRange = &[2]float64{min, max}
	return &filter
}

// PreferGroups specifies the priority of each group in the order they are presented.
// Only takes effect if mangaQuery.IncludeScanlationGroup is enabled.
//
// Note that this does not filter only the chapter translated by these groups. In case
// there are several version of a chapter then the groups specifed here will take precedence
// when filtered.
func (filter mangaFilter) PreferGroups(groups []string) *mangaFilter {
	filter.preferGroups = make(map[string]int)
	for i, group := range groups {
		filter.preferGroups[strings.ToLower(group)] = len(groups) - i
	}
	return &filter
}

// GetChapters returns list of chapter sastified the criterias.
func (filter mangaFilter) GetChapters() (chapters ChapterList) {
	for i, chapter := range filter.manga.Data {
		// Parsed chapter is stored in a slice, if chapter with same name already exists end of slice,
		// compare it with the current parsed chapter by these criteria:
		// - skip whichever is empty
		// - skip one with lower priority of scanlation group
		if len(chapters) > 0 && chapters[len(chapters)-1].GetChapter() == chapter.GetChapter() {
			old_chapter := chapters[len(chapters)-1]
			old_empty := old_chapter.GetPages() == 0
			new_empty := chapter.GetPages() == 0
			if (old_empty && new_empty) || (!old_empty && !new_empty) {
				if filter.preferGroups != nil && chapter.GetPages() > 0 {
					old_group := strings.ToLower(old_chapter.GetScanlationGroup())
					new_group := strings.ToLower(chapter.GetScanlationGroup())
					if filter.preferGroups[old_group] < filter.preferGroups[new_group] {
						chapters[len(chapters)-1] = &filter.manga.Data[i]
					}
				}
			} else if old_empty && !new_empty {
				chapters[len(chapters)-1] = &filter.manga.Data[i]
			}
			continue
		}

		isGood := true // this flag indicates the current parsed chapter sastifies all criterias
		if filter.volumes != nil {
			_, exist := filter.volumes[chapter.GetVolume()]
			isGood = isGood && exist
		}
		if filter.chapters != nil {
			_, exist := filter.chapters[chapter.GetChapter()]
			isGood = isGood && exist
		}
		if filter.volumeRange != nil {
			val, err := strconv.ParseFloat(chapter.GetVolume(), 64)
			isGood = isGood && err == nil && val >= filter.volumeRange[0] && val <= filter.volumeRange[1]
		}
		if filter.chapterRange != nil {
			val, err := strconv.ParseFloat(chapter.GetChapter(), 64)
			isGood = isGood && err == nil && val >= filter.chapterRange[0] && val <= filter.chapterRange[1]
		}
		if isGood {
			// the chapter sastified all criterias, save it in the map
			chapters = append(chapters, &filter.manga.Data[i])
		}
	}

	return
}
