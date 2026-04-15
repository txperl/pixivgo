package pixivgo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
)

// NovelDetailParams are the parameters for NovelDetail.
type NovelDetailParams struct {
	NovelID int
	NoAuth  bool
}

// NovelDetail returns detail information for a novel.
func (c *Client) NovelDetail(ctx context.Context, params NovelDetailParams) (*NovelInfo, error) {
	v := url.Values{"novel_id": {strconv.Itoa(params.NovelID)}}
	resp, err := c.doGet(ctx, "/v2/novel/detail", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	wrapper, err := parseResponse[NovelDetailResponse](resp)
	if err != nil {
		return nil, err
	}
	return &wrapper.Novel, nil
}

// NovelCommentsParams are the parameters for NovelComments.
type NovelCommentsParams struct {
	NovelID              int
	Offset               *int
	IncludeTotalComments *bool
	NoAuth               bool
}

// NovelComments returns comments on a novel.
func (c *Client) NovelComments(ctx context.Context, params NovelCommentsParams) (*NovelComments, error) {
	v := url.Values{"novel_id": {strconv.Itoa(params.NovelID)}}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	if params.IncludeTotalComments != nil {
		v.Set("include_total_comments", formatBool(*params.IncludeTotalComments))
	}
	resp, err := c.doGet(ctx, "/v1/novel/comments", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[NovelComments](resp)
}

// NovelRecommendedParams are the parameters for NovelRecommended.
type NovelRecommendedParams struct {
	IncludeRankingLabel       *bool
	Filter                    Filter
	Offset                    *int
	IncludeRankingNovels      *bool
	AlreadyRecommended        []string
	MaxBookmarkIDForRecommend *int
	IncludePrivacyPolicy      *string
	NoAuth                    bool
}

// NovelRecommended returns recommended novels.
func (c *Client) NovelRecommended(ctx context.Context, params NovelRecommendedParams) (*NovelListResponse, error) {
	includeRanking := "true"
	if params.IncludeRankingLabel != nil {
		includeRanking = formatBool(*params.IncludeRankingLabel)
	}
	v := url.Values{
		"include_ranking_label": {includeRanking},
		"filter":                {defaultFilter(params.Filter)},
	}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	if params.IncludeRankingNovels != nil {
		v.Set("include_ranking_novels", formatBool(*params.IncludeRankingNovels))
	}
	if params.MaxBookmarkIDForRecommend != nil {
		v.Set("max_bookmark_id_for_recommend", strconv.Itoa(*params.MaxBookmarkIDForRecommend))
	}
	if len(params.AlreadyRecommended) > 0 {
		v.Set("already_recommended", joinStrings(params.AlreadyRecommended))
	}
	if params.IncludePrivacyPolicy != nil {
		v.Set("include_privacy_policy", *params.IncludePrivacyPolicy)
	}
	resp, err := c.doGet(ctx, "/v1/novel/recommended", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[NovelListResponse](resp)
}

// NovelSeriesParams are the parameters for NovelSeries.
type NovelSeriesParams struct {
	SeriesID  int
	Filter    Filter
	LastOrder *string
	NoAuth    bool
}

// NovelSeries returns detail for a novel series.
func (c *Client) NovelSeries(ctx context.Context, params NovelSeriesParams) (*NovelSeriesResponse, error) {
	v := url.Values{
		"series_id": {strconv.Itoa(params.SeriesID)},
		"filter":    {defaultFilter(params.Filter)},
	}
	if params.LastOrder != nil {
		v.Set("last_order", *params.LastOrder)
	}
	resp, err := c.doGet(ctx, "/v2/novel/series", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[NovelSeriesResponse](resp)
}

// NovelNewParams are the parameters for NovelNew.
type NovelNewParams struct {
	Filter     Filter
	MaxNovelID *int
	NoAuth     bool
}

// NovelNew returns the latest novels.
func (c *Client) NovelNew(ctx context.Context, params NovelNewParams) (*NovelListResponse, error) {
	v := url.Values{"filter": {defaultFilter(params.Filter)}}
	if params.MaxNovelID != nil {
		v.Set("max_novel_id", strconv.Itoa(*params.MaxNovelID))
	}
	resp, err := c.doGet(ctx, "/v1/novel/new", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[NovelListResponse](resp)
}

// NovelFollowParams are the parameters for NovelFollow.
type NovelFollowParams struct {
	Restrict Restrict
	Offset   *int
	NoAuth   bool
}

// NovelFollow returns new novels from followed users. Requires authentication.
func (c *Client) NovelFollow(ctx context.Context, params NovelFollowParams) (*NovelListResponse, error) {
	restrict := string(params.Restrict)
	if restrict == "" {
		restrict = string(RestrictPublic)
	}
	v := url.Values{"restrict": {restrict}}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v1/novel/follow", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[NovelListResponse](resp)
}

var novelJSONRegex = regexp.MustCompile(`novel:\s({.+}),\s+isOwnWork`)

// WebviewNovelParams are the parameters for WebviewNovel.
type WebviewNovelParams struct {
	NovelID int
	Raw     bool // if true, returns raw HTML instead of parsed model
	NoAuth  bool
}

// WebviewNovelResult holds the result of WebviewNovel.
// Either Novel or RawHTML is set, depending on the Raw parameter.
type WebviewNovelResult struct {
	Novel   *WebviewNovel
	RawHTML string
}

// WebviewNovel fetches novel content from the webview endpoint.
// If params.Raw is true, RawHTML is populated. Otherwise, Novel is populated
// by extracting and parsing JSON from the HTML.
func (c *Client) WebviewNovel(ctx context.Context, params WebviewNovelParams) (*WebviewNovelResult, error) {
	v := url.Values{
		"id":             {strconv.Itoa(params.NovelID)},
		"viewer_version": {"20221031_ai"},
	}
	resp, err := c.doGet(ctx, "/webview/v2/novel", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &PixivError{
			Message:    fmt.Sprintf("read webview response error: %v", err),
			StatusCode: resp.StatusCode,
			Err:        err,
		}
	}

	htmlStr := string(body)

	if params.Raw {
		return &WebviewNovelResult{RawHTML: htmlStr}, nil
	}

	// Extract JSON from HTML: novel: {...}, isOwnWork
	match := novelJSONRegex.FindStringSubmatch(htmlStr)
	if len(match) < 2 {
		return nil, &PixivError{
			Message:    "extract novel content error: regex did not match",
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
			Body:       htmlStr,
		}
	}

	var novel WebviewNovel
	if err := json.Unmarshal([]byte(match[1]), &novel); err != nil {
		return nil, &PixivError{
			Message:    fmt.Sprintf("parse novel json error: %v", err),
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
			Body:       match[1],
			Err:        err,
		}
	}

	return &WebviewNovelResult{Novel: &novel}, nil
}
