package pixivgo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultHosts     = "https://app-api.pixiv.net"
	defaultUserAgent = "PixivIOSApp/7.13.3 (iOS 14.6; iPhone13,2)"
	defaultAppOS     = "ios"
	defaultAppOSVer  = "14.6"

	defaultClientID     = "MOBrBDS8blbauoSck0ZfDbtuzpyT"
	defaultClientSecret = "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj"
	defaultHashSecret   = "28c1fdd170a5204386cb1313c7077b34f83e4aaf4aa829ce78c231e05b0bae2c"
)

// Client is the Pixiv App-API client.
// It is safe for concurrent use after construction.
type Client struct {
	httpClient        *http.Client
	hosts             string
	clientID          string
	clientSecret      string
	hashSecret        string
	userAgent         string
	additionalHeaders http.Header

	mu           sync.RWMutex
	accessToken  string
	refreshToken string
	userID       string
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets a custom http.Client for all requests.
// Use this to inject proxied clients or bypass-configured transports.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithBaseURL overrides the default API base URL (https://app-api.pixiv.net).
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.hosts = strings.TrimRight(baseURL, "/") }
}

// WithAdditionalHeaders sets additional HTTP headers for all requests.
func WithAdditionalHeaders(h http.Header) Option {
	return func(c *Client) { c.additionalHeaders = h.Clone() }
}

// WithAcceptLanguage sets the Accept-Language header for tag translations.
// Common values: "en-us", "zh-cn", "ja".
func WithAcceptLanguage(lang string) Option {
	return func(c *Client) { c.additionalHeaders.Set("Accept-Language", lang) }
}

// NewClient creates a new Pixiv API client with the given options.
func NewClient(opts ...Option) *Client {
	c := &Client{
		httpClient:        http.DefaultClient,
		hosts:             defaultHosts,
		clientID:          defaultClientID,
		clientSecret:      defaultClientSecret,
		hashSecret:        defaultHashSecret,
		userAgent:         defaultUserAgent,
		additionalHeaders: make(http.Header),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// SetAPIProxy changes the API base URL (e.g., "http://app-api.pixivlite.com").
func (c *Client) SetAPIProxy(proxyHosts string) {
	c.hosts = strings.TrimRight(proxyHosts, "/")
}

// SetAdditionalHeaders overrides the additional headers for all requests.
func (c *Client) SetAdditionalHeaders(h http.Header) {
	c.additionalHeaders = h.Clone()
}

// SetAcceptLanguage sets the Accept-Language header for tag translations.
func (c *Client) SetAcceptLanguage(lang string) {
	c.additionalHeaders.Set("Accept-Language", lang)
}

// getAccessToken returns the current access token (thread-safe).
func (c *Client) getAccessToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.accessToken
}

// requireAuth checks that the client is authenticated.
func (c *Client) requireAuth() error {
	if c.getAccessToken() == "" {
		return ErrAuthRequired
	}
	return nil
}

// doRequest executes an HTTP request against the Pixiv API.
// If noAuth is false, it adds the Bearer token.
// path can be an absolute URL or a relative path appended to c.hosts.
func (c *Client) doRequest(ctx context.Context, method, path string, params url.Values, data url.Values, noAuth bool) (*http.Response, error) {
	// Build URL
	reqURL := path
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		reqURL = c.hosts + path
	}

	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	var bodyReader io.Reader
	if data != nil {
		bodyReader = strings.NewReader(data.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return nil, &PixivError{Message: fmt.Sprintf("create request error: %v", err), Err: err}
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// Set iOS app headers
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("App-OS", defaultAppOS)
	req.Header.Set("App-OS-Version", defaultAppOSVer)

	// When using a proxy/bypass host, set the original Host header
	if c.hosts != defaultHosts && !strings.HasPrefix(path, "http") {
		req.Host = "app-api.pixiv.net"
	}

	// Apply additional headers
	for key, values := range c.additionalHeaders {
		for _, v := range values {
			req.Header.Set(key, v)
		}
	}

	// Apply auth
	if !noAuth {
		if err := c.requireAuth(); err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.getAccessToken())
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &PixivError{Message: fmt.Sprintf("request %s %s error: %v", method, path, err), Err: err}
	}

	return resp, nil
}

// doGet is a convenience wrapper for GET requests.
func (c *Client) doGet(ctx context.Context, path string, params url.Values, noAuth bool) (*http.Response, error) {
	return c.doRequest(ctx, http.MethodGet, path, params, nil, noAuth)
}

// doPost is a convenience wrapper for POST requests.
func (c *Client) doPost(ctx context.Context, path string, data url.Values, noAuth bool) (*http.Response, error) {
	return c.doRequest(ctx, http.MethodPost, path, nil, data, noAuth)
}

// parseResponse reads the response body, checks HTTP status, and unmarshals JSON into T.
func parseResponse[T any](resp *http.Response) (*T, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &PixivError{
			Message:    fmt.Sprintf("read response body error: %v", err),
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
			Err:        err,
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, &PixivError{
			Message:    fmt.Sprintf("API returned HTTP %d", resp.StatusCode),
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
			Body:       string(body),
		}
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, &PixivError{
			Message:    fmt.Sprintf("json unmarshal error: %v", err),
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
			Body:       string(body),
			Err:        err,
		}
	}

	return &result, nil
}
