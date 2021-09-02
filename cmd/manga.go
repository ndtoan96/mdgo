/*
Copyright © 2021 Nguyen Duc Toan <ntoan96@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/ndtoan96/mgdex"
	"github.com/spf13/cobra"
)

// mangaCmd represents the manga command
var mangaCmd = &cobra.Command{
	Use:   "manga",
	Short: "Download chapters from a manga",
	Long: `This command is used to download chapters from a manga.
	
It takes manga id or url as input and provides several filters (by default,
all chapters will be downloaded), if more than one filters are used,
the condition will be AND together. Chapters can be downloaded to folder or
to archive. See the flags for more detail.`,
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		output, _ := cmd.Flags().GetString("output")
		prefix, _ := cmd.Flags().GetString("prefix")
		language, _ := cmd.Flags().GetString("language")
		chapters, _ := cmd.Flags().GetStringSlice("chapters")
		volumes, _ := cmd.Flags().GetStringSlice("volumes")
		chapterRange, _ := cmd.Flags().GetFloat64Slice("chapter-range")
		volumeRange, _ := cmd.Flags().GetFloat64Slice("volume-range")
		groups, _ := cmd.Flags().GetStringSlice("groups")
		archive, _ := cmd.Flags().GetString("archive")
		noRun, _ := cmd.Flags().GetBool("dry-run")
		raw, _ := cmd.Flags().GetBool("raw")

		if len(chapterRange) > 0 && len(chapterRange) != 2 {
			fmt.Println("chapter-range takes 2 values, found", len(chapterRange))
			return
		}

		if len(volumeRange) > 0 && len(volumeRange) != 2 {
			fmt.Println("volume-range takes 2 values, found", len(volumeRange))
			return
		}

		url_pattern := regexp.MustCompile(`mangadex\.org/title/([\w-]+)`)
		var id string
		if m := url_pattern.FindStringSubmatch(input); len(m) == 2 {
			id = m[1]
		} else {
			id = input
		}

		query := mgdex.MangaQuery(id).Language(language).Limit(500).Order("asc")
		if groups != nil {
			query = query.IncludeScanlationGroup()
		}
		manga, err := query.GetManga()
		if err != nil {
			fmt.Println(err)
		}

		filter := manga.Filter()
		if len(chapterRange) == 2 {
			filter = filter.ChapterRange([2]float64{chapterRange[0], chapterRange[1]})
		}
		if len(volumeRange) == 2 {
			filter = filter.VolumeRange([2]float64{volumeRange[0], volumeRange[1]})
		}
		if len(chapters) > 0 {
			filter = filter.Chapters(chapters)
		}
		if len(volumes) > 0 {
			filter = filter.Volumes(volumes)
		}
		if len(groups) > 0 {
			filter.PreferGroups(groups)
		}

		chapterList := filter.GetChapters()
		if noRun {
			fmt.Println("These chapters will be downloaded:")
			for i, chap := range chapterList {
				fmt.Println(i, "- chapter", chap.Chapter())
			}
		} else {
			if archive != "" {
				chapterList.DownloadAsZip(!raw, filepath.Join(output, prefix), archive)
			} else {
				chapterList.Download(!raw, filepath.Join(output, prefix))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(mangaCmd)

	mangaCmd.Flags().StringP("input", "i", "", "manga id or url")
	mangaCmd.Flags().StringP("output", "o", ".", "output folder")
	mangaCmd.Flags().StringP("prefix", "p", "chapter_", "prefix of downloaded folders or archives name")
	mangaCmd.Flags().StringSliceP("chapters", "c", nil, `list of chapters, can be used multiple times
or values can be seperated by comma`)
	mangaCmd.Flags().StringSliceP("volumes", "v", nil, `list of volumes, can be used multiple times
or values can be seperated by comma`)
	mangaCmd.Flags().Float64SliceP("chapter-range", "C", nil, `range of chapter, take two values seperated by comma`)
	mangaCmd.Flags().Float64SliceP("volume-range", "V", nil, `range of volume, take two values seperated by comma`)
	mangaCmd.Flags().StringP("language", "l", "en", "translated language, multiple choice is not supported")
	mangaCmd.Flags().StringSliceP("groups", "g", nil, `list of prefered scanlation groups. In case a chapter has several versions
made by different groups, the groups specified here will take precedence
according to the order they are listed. Otherwise first version will be taken`)
	mangaCmd.Flags().StringP("archive", "a", "", `extension of zip files. If not specified, chapters will not be zipped`)
	mangaCmd.Flags().BoolP("dry-run", "n", false, "only print list of chapters, not actually download them")
	mangaCmd.Flags().BoolP("raw", "r", false, `by default compressed images are downloaded to save data, 
turn on this flag to download original quality images`)
	mangaCmd.MarkFlagRequired("input")
	mangaCmd.MarkFlagDirname("output")

	mangaCmd.Example = `mdgo manga -i https://mangadex.org/title/37f5cce0-8070-4ada-96e5-fa24b1bd4ff9/kaguya-sama-wa-kokurasetai-tensai-tachi-no-renai-zunousen -o kaguya -p kaguya_chap_ -C="0,10" -a "cbz"`
}
