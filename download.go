package pixivgo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const defaultReferer = "https://app-api.pixiv.net/"

// DownloadOptions controls how images are downloaded.
type DownloadOptions struct {
	Path    string // directory to save to (default: current directory)
	Name    string // filename override (default: basename from URL)
	Prefix  string // filename prefix
	Replace bool   // overwrite existing files (default: false)
	Referer string // Referer header (default: "https://app-api.pixiv.net/")
}

// Download downloads an image from the given URL to a file.
// Returns the path of the saved file.
func (c *Client) Download(ctx context.Context, imageURL string, opts *DownloadOptions) (string, error) {
	if opts == nil {
		opts = &DownloadOptions{}
	}

	// Determine file path
	name := opts.Name
	if name == "" {
		name = filepath.Base(imageURL)
	}
	if opts.Prefix != "" {
		name = opts.Prefix + name
	}

	dir := opts.Path
	if dir == "" {
		dir = "."
	}

	filePath := filepath.Join(dir, name)

	// Check if file exists
	if !opts.Replace {
		if _, err := os.Stat(filePath); err == nil {
			return filePath, nil // file already exists, skip
		}
	}

	// Download
	referer := opts.Referer
	if referer == "" {
		referer = defaultReferer
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return "", &PixivError{Message: fmt.Sprintf("create download request error: %v", err), Err: err}
	}
	req.Header.Set("Referer", referer)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", &PixivError{Message: fmt.Sprintf("download request error: %v", err), Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return "", &PixivError{
			Message:    fmt.Sprintf("download failed: HTTP %d", resp.StatusCode),
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", &PixivError{Message: fmt.Sprintf("create file error: %v", err), Err: err}
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", &PixivError{Message: fmt.Sprintf("write file error: %v", err), Err: err}
	}

	return filePath, nil
}

// DownloadToWriter downloads an image from the given URL and writes it to w.
func (c *Client) DownloadToWriter(ctx context.Context, imageURL string, w io.Writer) error {
	referer := defaultReferer

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return &PixivError{Message: fmt.Sprintf("create download request error: %v", err), Err: err}
	}
	req.Header.Set("Referer", referer)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &PixivError{Message: fmt.Sprintf("download request error: %v", err), Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return &PixivError{
			Message:    fmt.Sprintf("download failed: HTTP %d", resp.StatusCode),
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
		}
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		return &PixivError{Message: fmt.Sprintf("write to writer error: %v", err), Err: err}
	}

	return nil
}
