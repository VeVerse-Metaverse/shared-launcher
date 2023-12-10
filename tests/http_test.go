package tests

import (
	"context"
	"dev.hackerman.me/artheon/l7-shared-launcher/download"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadFileUnknownSize(t *testing.T) {
	path := "temp/download-unknown-size.bin"
	url := "http://speedtest.ftp.otenet.gr/files/test100k.db"

	err := download.File(context.Background(), path, url, 0)
	if err != nil {
		t.Fatal(err)
	}

	stat, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if stat.Size() == 0 {
		t.Fatal("file is empty")
	}
}

func TestDownloadFileKnownSize(t *testing.T) {
	path := "temp/download-known-size.bin"
	url := "http://speedtest.ftp.otenet.gr/files/test100k.db"

	// delete file if exists
	stat, err := os.Stat(path)
	if err == nil {
		err := os.Remove(path)
		if err != nil {
			t.Fatal(err)
		}
	}

	// 100kb
	size := int64(100 * 1024)

	err = download.File(context.Background(), path, url, size)
	if err != nil {
		t.Fatal(err)
	}

	stat, err = os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if stat.Size() == 0 {
		t.Fatal("file is empty")
	}
}

func TestDownloadFileKnownSizeExists(t *testing.T) {
	path := "temp/download-known-size-exists.bin"
	url := "http://speedtest.ftp.otenet.gr/files/test100k.db"

	// delete file if exists
	stat, err := os.Stat(path)
	if err == nil {
		err := os.Remove(path)
		if err != nil {
			t.Fatal(err)
		}
	}

	// 100kb
	size := int64(100 * 1024)

	err = download.File(context.Background(), path, url, size)
	if err != nil {
		t.Fatal(err)
	}

	stat, err = os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if stat.Size() == 0 {
		t.Fatal("file is empty")
	}

	// try to download again having the file
	err = download.File(context.Background(), path, url, size)
	if err != nil && !download.ErrFileExists.Is(err) {
		t.Fatal(err)
	}

	stat, err = os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if stat.Size() == 0 {
		t.Fatal("file is empty")
	}
}

func TestDownloadFileInvalidPath(t *testing.T) {
	path := ""
	url := "http://speedtest.ftp.otenet.gr/files/test100k.db"

	err := download.File(context.Background(), path, url, 0)
	if err != nil && download.ErrInvalidPath.Is(err) {
		return // pass, expected error
	}

	t.Fatal(err)
}

func TestDownloadFileExecutable(t *testing.T) {
	srcPath, err := filepath.Abs("download/Metaverse.exe")
	if err != nil {
		t.Fatal(err)
	}

	stat, err := os.Stat(srcPath)
	if err != nil {
		t.Fatal(err)
	}

	if stat.Size() == 0 {
		t.Fatal("file is empty")
	}

	f, err := os.Open(srcPath)
	if err != nil {
		t.Fatal(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(f)

	fileContent := make([]byte, stat.Size())
	_, err = f.Read(fileContent)
	if err != nil {
		t.Fatal(err)
	}

	// Start a local HTTP server to serve the file
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(fileContent)
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	path := "temp/download-executable.exe"
	url := ts.URL

	err = download.File(context.Background(), path, url, 0)
	if err != nil {
		t.Fatal(err)
	}

	stat, err = os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if stat.Size() == 0 {
		t.Fatal("file is empty")
	}

	// todo: check mode
}
