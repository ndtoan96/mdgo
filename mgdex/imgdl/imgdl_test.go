package imgdl

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func Test_GetImageData_ImageTypes(t *testing.T) {
	tls1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "image/jpg")
	}))
	defer tls1.Close()
	tls2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "image/png")
	}))
	defer tls2.Close()
	tls3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "image/gif")
	}))
	defer tls3.Close()
	tls4 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer tls4.Close()

	var ext string
	var err error
	ext, _, err = GetImageData(tls1.URL)
	if err != nil || ext != "jpg" {
		t.FailNow()
	}
	ext, _, err = GetImageData(tls2.URL)
	if err != nil || ext != "png" {
		t.FailNow()
	}
	ext, _, err = GetImageData(tls3.URL)
	if err != nil || ext != "gif" {
		t.FailNow()
	}
	ext, _, err = GetImageData(tls4.URL)
	if err != nil || ext != "jpg" {
		t.FailNow()
	}
}

func Test_GetImageData_HttpError(t *testing.T) {
	tls := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer tls.Close()

	_, _, err := GetImageData(tls.URL)
	if err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf(`Expected 404 Not found, got %v`, err)
	}

	_, _, err = GetImageData("http://some.bs.url/abcd")
	if err == nil {
		t.Fatal(`Expected error, got none`)
	}
}

func Test_GetImageData_NotImage(t *testing.T) {
	tls := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Some text response")
	}))
	defer tls.Close()

	_, _, err := GetImageData(tls.URL)
	if !strings.Contains(err.Error(), "ot an image") {
		t.Fatalf(`Expected not an image error, got %v`, err)
	}
}

func Test_DownLoadImage_ok(t *testing.T) {
	imgname := "imgdl/test-data/img1.png"
	tls := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "image/png")
		f, _ := os.Open(imgname)
		io.Copy(w, f)
	}))
	defer tls.Close()

	tmpfile, _ := ioutil.TempFile("", "a")
	err := DownloadImgage(tls.URL, tmpfile.Name())
	if err != nil {
		t.FailNow()
	}
	dat1, _ := os.ReadFile(imgname)
	dat2, _ := os.ReadFile(tmpfile.Name() + ".png")
	if bytes.Compare(dat1, dat2) != 0 {
		t.FailNow()
	}
}

func Test_DownLoadImage_failed(t *testing.T) {
	tls := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer tls.Close()

	tmpfile, _ := ioutil.TempFile("", "a")
	var err error
	err = DownloadImgage(tls.URL, tmpfile.Name())
	if err == nil {
		t.FailNow()
	}
	err = DownloadImgage("http://bs_url.fake_it.man/fff", tmpfile.Name())
	if err == nil {
		t.FailNow()
	}
}

func Test_DownLoadImages_ok(t *testing.T) {
	imgname1 := "imgdl/test-data/img1.png"
	tls1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "image/png")
		f, _ := os.Open(imgname1)
		io.Copy(w, f)
	}))
	defer tls1.Close()
	imgname2 := "imgdl/test-data/img2.jpg"
	tls2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, _ := os.Open(imgname2)
		io.Copy(w, f)
	}))
	defer tls2.Close()

	tmpfile, _ := ioutil.TempFile("", "a")
	err := DownloadImgages([]string{tls1.URL, tls2.URL}, tmpfile.Name())
	if err != nil {
		t.FailNow()
	}
	rd1, _ := os.ReadFile(imgname1)
	rd2, _ := os.ReadFile(imgname2)
	wr1, _ := os.ReadFile(tmpfile.Name() + "00.png")
	wr2, _ := os.ReadFile(tmpfile.Name() + "01.jpg")
	if bytes.Compare(rd1, wr1) != 0 || bytes.Compare(rd2, wr2) != 0 {
		t.FailNow()
	}
}
