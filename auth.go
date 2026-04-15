package pixivgo

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const authEndpoint = "https://oauth.secure.pixiv.net/auth/token"

// Auth authenticates with Pixiv using a refresh token.
// It stores the resulting access token and refresh token on the client.
// Password authentication is not supported (deprecated by Pixiv).
func (c *Client) Auth(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	if refreshToken == "" {
		refreshToken = c.getRefreshToken()
	}
	if refreshToken == "" {
		return nil, &PixivError{Message: "no refresh_token provided"}
	}

	// Build auth URL — use bypass hosts if configured
	authURL := authEndpoint
	var hostOverride string
	if c.hosts != defaultHosts {
		authURL = c.hosts + "/auth/token"
		hostOverride = "oauth.secure.pixiv.net"
	}

	// Request signing
	localTime := time.Now().UTC().Format("2006-01-02T15:04:05+00:00")
	hash := md5.Sum([]byte(localTime + c.hashSecret))
	clientHash := fmt.Sprintf("%x", hash)

	// Build form data
	data := url.Values{
		"get_secure_url": {"1"},
		"client_id":      {c.clientID},
		"client_secret":  {c.clientSecret},
		"grant_type":     {"refresh_token"},
		"refresh_token":  {refreshToken},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, &PixivError{Message: fmt.Sprintf("create auth request error: %v", err), Err: err}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("App-OS", defaultAppOS)
	req.Header.Set("App-OS-Version", defaultAppOSVer)
	req.Header.Set("X-Client-Time", localTime)
	req.Header.Set("X-Client-Hash", clientHash)

	if hostOverride != "" {
		req.Host = hostOverride
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &PixivError{Message: fmt.Sprintf("auth request error: %v", err), Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &PixivError{Message: fmt.Sprintf("read auth response error: %v", err), Err: err}
	}

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusMovedPermanently &&
		resp.StatusCode != http.StatusFound {
		return nil, &PixivError{
			Message:    fmt.Sprintf("auth failed, check refresh_token; HTTP %d: %s", resp.StatusCode, string(body)),
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
			Body:       string(body),
		}
	}

	var tokenResp authTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, &PixivError{
			Message:    fmt.Sprintf("parse auth response error: %v", err),
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
			Body:       string(body),
			Err:        err,
		}
	}

	authResp := &tokenResp.Response

	// Store tokens
	c.mu.Lock()
	c.accessToken = authResp.AccessToken
	c.refreshToken = authResp.RefreshToken
	c.userID = strconv.Itoa(authResp.User.ID.Int())
	c.mu.Unlock()

	return authResp, nil
}

// SetAuth manually sets the access token and optionally the refresh token.
// Use this when you already have valid tokens.
func (c *Client) SetAuth(accessToken string, refreshToken string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.accessToken = accessToken
	if refreshToken != "" {
		c.refreshToken = refreshToken
	}
}

// getRefreshToken returns the current refresh token (thread-safe).
func (c *Client) getRefreshToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.refreshToken
}

// UserID returns the authenticated user's ID.
func (c *Client) UserID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.userID
}
