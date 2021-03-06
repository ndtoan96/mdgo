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
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ndtoan96/mgdex"
	"github.com/spf13/cobra"
)

// Return code:
//	1 - Request error
//	2 - Download error
//	3 - Invalid argument

// mangaCmd represents the manga command
var mangaCmd = &cobra.Command{
	Use:     "manga <input> [output]",
	Aliases: []string{"m"},
	Short:   "Download multiple chapters from a manga",
	Long: `This command is used to download chapters from a manga.
	
It takes manga id or url as input and provides several filters (by default,
all chapters will be downloaded), if more than one filters are used,
the condition will be AND together. Chapters can be downloaded to folder or
to archive. See the flags for more detail.`,
	Example: `mdgo manga https://mangadex.org/title/37f5cce0-8070-4ada-96e5-fa24b1bd4ff9/kaguya-sama-wa-kokurasetai-tensai-tachi-no-renai-zunousen kaguya -p kaguya_chap_ -C="0,10" -a "cbz"`,
	Args:    cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		input := cmd.Flags().Arg(0)
		output := cmd.Flags().Arg(1)
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
		getAll, _ := cmd.Flags().GetBool("all")
		lastN, _ := cmd.Flags().GetUint("last")

		// Check arguments validity
		if len(chapterRange) > 0 && len(chapterRange) != 2 {
			fmt.Println("chapter-range takes 2 values, found", len(chapterRange))
			os.Exit(3)
		}
		if len(volumeRange) > 0 && len(volumeRange) != 2 {
			fmt.Println("volume-range takes 2 values, found", len(volumeRange))
			os.Exit(3)
		}

		// Extract manga id from url
		url_pattern := regexp.MustCompile(`mangadex\.org/title/([\w-]+)`)
		var mangaId string
		if m := url_pattern.FindStringSubmatch(input); len(m) == 2 {
			mangaId = m[1]
		} else {
			mangaId = input
		}

		// Get manga title if prefix contains :m
		if strings.Contains(prefix, ":m") {
			mangaInfo, err := mgdex.GetMangaInfo(mangaId)
			if err == nil {
				mangaTitle := mangaInfo.GetMainTitle()
				prefix = strings.ReplaceAll(prefix, ":m", mangaTitle)
			} else {
				fmt.Println("Warning: Can't get manga title, skip.")
			}
		}

		query := mgdex.MangaQuery(mangaId).Language(language).Limit(500).Order("asc")
		if groups != nil {
			query = query.IncludeScanlationGroup()
		}
		manga, err := query.GetManga()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Try to query continously with offset 500 until there's no chapters left
		if getAll {
			cnt := 1
			for {
				query_tmp := mgdex.MangaQuery(mangaId).Language(language).Limit(500).Order("asc").Offset(500 * cnt)
				if groups != nil {
					query = query.IncludeScanlationGroup()
				}
				manga_tmp, err := query_tmp.GetManga()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if manga_tmp.Length() == 0 {
					break
				}
				manga.Append(manga_tmp)
				cnt++
			}
		}

		var chapterList mgdex.ChapterList
		filter := manga.Filter()
		if len(groups) > 0 {
			filter.PreferGroups(groups)
		}

		// check if --last/-L argument is used
		if lastN > 0 {
			// --last/-L argument cannot be used together with volume and chapter filters
			if len(chapterRange) > 0 || len(volumeRange) > 0 || len(chapters) > 0 || len(volumes) > 0 {
				fmt.Println("'--last/-L' argument can not be used together with chapter or volume filters")
				os.Exit(3)
			}
			// Fetch all chapter
			chapterList = filter.GetChapters()
			// Get the last N chapter
			if uint(len(chapterList)) > lastN {
				chapterList = chapterList[uint(len(chapterList))-lastN:]
			}
		} else {
			// Add argument to the filter
			if len(chapterRange) == 2 {
				filter = filter.ChapterRange(chapterRange[0], chapterRange[1])
			}
			if len(volumeRange) == 2 {
				filter = filter.VolumeRange(volumeRange[0], volumeRange[1])
			}
			if len(chapters) > 0 {
				filter = filter.Chapters(chapters)
			}
			if len(volumes) > 0 {
				filter = filter.Volumes(volumes)
			}
			chapterList = filter.GetChapters()
		}

		if noRun {
			// dry run, print the chapters name
			fmt.Println("These chapters will be downloaded:")
			for i, chap := range chapterList {
				fmt.Println(i+1, "- chapter", chap.GetChapter())
			}
		} else {
			var success bool
			if archive != "" {
				success = chapterList.DownloadAsZip(!raw, filepath.Join(output, prefix), archive)
			} else {
				success = chapterList.Download(!raw, filepath.Join(output, prefix))
			}
			if !success {
				fmt.Println("Failure in downloading one or more chapters!")
				os.Exit(2)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(mangaCmd)

	mangaCmd.Flags().StringP("prefix", "p", "chapter_", "Prefix of downloaded folders or archives name")
	mangaCmd.Flags().StringSliceP("chapters", "c", nil, `List of chapters, can be used multiple times
or values can be seperated by comma`)
	mangaCmd.Flags().StringSliceP("volumes", "v", nil, `List of volumes, can be used multiple times
or values can be seperated by comma`)
	mangaCmd.Flags().Float64SliceP("chapter-range", "C", nil, `Range of chapter, take two values seperated by comma`)
	mangaCmd.Flags().Float64SliceP("volume-range", "V", nil, `Range of volume, take two values seperated by comma`)
	mangaCmd.Flags().StringP("language", "l", "en", "Translated language, multiple choice is not supported")
	mangaCmd.Flags().StringSliceP("groups", "g", nil, `List of prefered scanlation groups. In case a chapter has several versions
made by different groups, the groups specified here will take precedence
according to the order they are listed. Otherwise first version will be taken`)
	mangaCmd.Flags().StringP("archive", "a", "", `Extension of zip files. If not specified, chapters will not be zipped`)
	mangaCmd.Flags().BoolP("dry-run", "n", false, "Only print list of chapters, not actually download them")
	mangaCmd.Flags().BoolP("raw", "r", false, `By default compressed images are downloaded to save data, 
	turn on this flag to download original quality images`)
	mangaCmd.Flags().Bool("all", false, `Try to download all chapters, use this tag when manga has too many chapters
	and normal mode cannot download them all (can use --dry-run to check first)`)
	mangaCmd.Flags().UintP("last", "L", 0, "Last n chapters to download, cannot be used together with volume and chapter filters")
	mangaCmd.SetHelpTemplate(mangaCmd.HelpTemplate() + fmt.Sprintf(`Args:
  %-10s Manga ID or url
  %-10s Folder to save downloaded chapters, if not set, current folder will be used

`, "input", "output"))

}
