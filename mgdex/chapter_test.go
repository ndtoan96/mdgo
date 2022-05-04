package mgdex

import (
	"testing"
)

func Test_GetChapter_OK(t *testing.T) {
	chapter, err := GetChapter("7ff854cf-dc17-4fdd-99d4-bc8f5d623b60")
	if err != nil {
		t.Fatal(err)
	}
	if chapter.GetId() != "7ff854cf-dc17-4fdd-99d4-bc8f5d623b60" || chapter.GetChapter() != "38" || chapter.GetVolume() != "4" {
		t.FailNow()
	}
}

func Test_GetChapter_Fail(t *testing.T) {
	_, err := GetChapter("rubbish")
	if err == nil {
		t.FailNow()
	}
}
