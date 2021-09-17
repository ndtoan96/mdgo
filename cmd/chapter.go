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
	"regexp"

	"github.com/ndtoan96/mgdex"
	"github.com/spf13/cobra"
)

// chapterCmd represents the chapter command
var chapterCmd = &cobra.Command{
	Use:     "chapter <input> [output]",
	Aliases: []string{"c"},
	Short:   "Download a single chapter",
	Long: `This command is used to download a single chapter.

It takes chapter id or url as input, the downloaded pages are named
with format "page_xx" and image extension will be automatically deduced.
Chapter can be downloaded to a folder, or compressed in an archive. See
the flags for detail.`,
	Example: `mdgo chapter abc-dxy-zhtkfj-skfk -a "cbz" manga/chapter`,
	Args:    cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		input := cmd.Flags().Arg(0)
		path := cmd.Flags().Arg(1)
		raw, _ := cmd.Flags().GetBool("raw")
		archive, _ := cmd.Flags().GetString("archive")

		url_pattern := regexp.MustCompile(`mangadex\.org/chapter/([\w-]+)`)
		var id string
		if m := url_pattern.FindStringSubmatch(input); len(m) == 2 {
			id = m[1]
		} else {
			id = input
		}
		chapter, err := mgdex.GetChapter(id)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if archive != "" {
			if path == "" || path == "." {
				path = "chapter"
			}
			path += "." + archive
			err = chapter.DownloadAsZip(!raw, path)
		} else {
			err = chapter.Download(!raw, path)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(chapterCmd)

	chapterCmd.Flags().StringP("archive", "a", "", "Archive the downloaded files")
	chapterCmd.Flags().BoolP("raw", "r", false, `By default compressed images are downloaded to save data, 
turn on this flag to download original quality images`)
	chapterCmd.SetHelpTemplate(chapterCmd.HelpTemplate() + fmt.Sprintf(`Args:
  %-10s Chapter ID or url
  %-10s Folder (or file name in case --archive is set) to save downloaded chapter, if not set
  %-10s current folder will be used

`, "input", "output", ""))
}
