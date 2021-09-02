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
	"regexp"

	"github.com/ndtoan96/mgdex"
	"github.com/spf13/cobra"
)

// chapterCmd represents the chapter command
var chapterCmd = &cobra.Command{
	Use:   "chapter",
	Short: "Download a single chapter",
	Long: `This command is used to download a single chapter.

It takes chapter id or url as input, the downloaded pages are named
with format "page_xx" and image extension will be automatically deduced.
Chapter can be downloaded to a folder, or compressed in an archive. See
the flags for detail.`,
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		raw, _ := cmd.Flags().GetBool("raw")
		path, _ := cmd.Flags().GetString("output")
		zip, _ := cmd.Flags().GetBool("zip")

		url_pattern := regexp.MustCompile(`mangadex\.org/chapter/([\w-]+)`)
		var id string
		if m := url_pattern.FindStringSubmatch(input); len(m) == 2 {
			id = m[1]
		} else {
			id = input
		}
		chapter, err := mgdex.GetChapter(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		if zip {
			err = chapter.DownloadAsZip(!raw, path)
		} else {
			err = chapter.Download(!raw, path)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(chapterCmd)

	chapterCmd.Flags().StringP("input", "i", "", "chapter id or url")
	chapterCmd.Flags().BoolP("zip", "z", false, "archive the downloaded files")
	chapterCmd.Flags().StringP("output", "o", ".", "output path, unexisting parent folders will be created.")
	chapterCmd.Flags().BoolP("raw", "r", false, `by default compressed images are downloaded to save data, 
turn on this flag to download original quality images`)
	chapterCmd.MarkFlagRequired("input")
	chapterCmd.MarkFlagDirname("output")

	chapterCmd.Example = `mdgo chapter -i abc-dxy-zhtkfj-skfk -z -o manga/chapter.cbz`
}
