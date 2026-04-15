package pixivgo

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestDownload_Success(t *testing.T) {
	content := []byte("fake image data for testing")
	var receivedReferer string

	client, ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		receivedReferer = r.Header.Get("Referer")
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	})

	tmpDir := t.TempDir()
	imageURL := ts.URL + "/img/test_image.png"

	path, err := client.Download(context.Background(), imageURL, &DownloadOptions{
		Path: tmpDir,
	})
	if err != nil {
		t.Fatalf("Download error: %v", err)
	}

	// Verify file was created with correct content
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if !bytes.Equal(data, content) {
		t.Errorf("file content = %q, want %q", data, content)
	}

	// Verify filename derived from URL
	if filepath.Base(path) != "test_image.png" {
		t.Errorf("filename = %q, want \"test_image.png\"", filepath.Base(path))
	}

	// Verify Referer header (migrated from Python test_download)
	if receivedReferer != "https://app-api.pixiv.net/" {
		t.Errorf("Referer = %q, want \"https://app-api.pixiv.net/\"", receivedReferer)
	}
}

func TestDownloadToWriter_Success(t *testing.T) {
	content := []byte("binary content here")
	var receivedReferer string

	client, ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		receivedReferer = r.Header.Get("Referer")
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	})

	var buf bytes.Buffer
	err := client.DownloadToWriter(context.Background(), ts.URL+"/img/test.png", &buf)
	if err != nil {
		t.Fatalf("DownloadToWriter error: %v", err)
	}

	if !bytes.Equal(buf.Bytes(), content) {
		t.Errorf("buffer = %q, want %q", buf.Bytes(), content)
	}

	if receivedReferer != "https://app-api.pixiv.net/" {
		t.Errorf("Referer = %q, want \"https://app-api.pixiv.net/\"", receivedReferer)
	}
}

func TestDownload_SkipExisting(t *testing.T) {
	client, ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("new content"))
	})

	tmpDir := t.TempDir()
	existingContent := []byte("existing content")
	existingFile := filepath.Join(tmpDir, "test.png")
	if err := os.WriteFile(existingFile, existingContent, 0644); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	path, err := client.Download(context.Background(), ts.URL+"/img/test.png", &DownloadOptions{
		Path:    tmpDir,
		Replace: false,
	})
	if err != nil {
		t.Fatalf("Download error: %v", err)
	}

	// File should not be overwritten
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if !bytes.Equal(data, existingContent) {
		t.Errorf("file was overwritten; got %q, want %q", data, existingContent)
	}
}

func TestDownload_NilOpts(t *testing.T) {
	client, ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("data"))
	})

	// Change to temp dir so default path "." doesn't write in project root
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	t.Cleanup(func() { os.Chdir(origDir) })

	_, err := client.Download(context.Background(), ts.URL+"/img/file.png", nil)
	if err != nil {
		t.Fatalf("Download with nil opts should not error: %v", err)
	}
}

func TestDownload_HTTPError(t *testing.T) {
	client, ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	tmpDir := t.TempDir()
	_, err := client.Download(context.Background(), ts.URL+"/img/fail.png", &DownloadOptions{
		Path: tmpDir,
	})
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}

	var pe *PixivError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *PixivError, got %T", err)
	}
	if pe.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", pe.StatusCode)
	}
}

func TestDownloadToWriter_HTTPError(t *testing.T) {
	client, ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	var buf bytes.Buffer
	err := client.DownloadToWriter(context.Background(), ts.URL+"/img/fail.png", &buf)
	if err == nil {
		t.Fatal("expected error for HTTP 403")
	}

	var pe *PixivError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *PixivError, got %T", err)
	}
	if pe.StatusCode != 403 {
		t.Errorf("StatusCode = %d, want 403", pe.StatusCode)
	}
}
