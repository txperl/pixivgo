package pixivgo

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrAuthRequired is returned when an API method requires authentication
// but no access token has been set.
var ErrAuthRequired = errors.New("pixivgo: authentication required; call Auth() or SetAuth() first")

// PixivError represents an error from the Pixiv API or the HTTP transport layer.
type PixivError struct {
	Message    string
	StatusCode int
	Header     http.Header
	Body       string
	Err        error
}

func (e *PixivError) Error() string {
	if e.StatusCode != 0 {
		return fmt.Sprintf("pixivgo: %s (HTTP %d)", e.Message, e.StatusCode)
	}
	return fmt.Sprintf("pixivgo: %s", e.Message)
}

func (e *PixivError) Unwrap() error {
	return e.Err
}
