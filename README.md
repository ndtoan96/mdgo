[![Go](https://github.com/ndtoan96/mdgo/actions/workflows/go.yml/badge.svg)](https://github.com/ndtoan96/mdgo/actions/workflows/go.yml)
# mdgo
CLI tool for downloading manga from mangadex

# Install
Make sure go is installed on your computer. Open command prompt or terminal, run
```
go install github.com/ndtoan96/mdgo@latest
```

# Usage
`mdgo` includes two commands, which is `manga` and `chapter`. `manga` is used to download multiple chapters of a manga, while `chapter` is used to download one single chapter.
Run `mdgo manga -h` or `mdgo chapter -h` for their respective usage.

# Example
To download chapter in range [0, 10] of Kaguya-sama, save in folder `kaguya` with each chapter named `kaguya_chap_<chap-number>`, and zip all of them in "cbz" file, run this command:
```
mdgo manga -i https://mangadex.org/title/37f5cce0-8070-4ada-96e5-fa24b1bd4ff9/kaguya-sama-wa-kokurasetai-tensai-tachi-no-renai-zunousen -o kaguya -p kaguya_chap_ -C="0,10" -a "cbz"
```
