package pixivgo

import (
	"context"
	"net/url"
	"strconv"
)

// IllustDetailParams are the parameters for IllustDetail.
type IllustDetailParams struct {
	IllustID int
	NoAuth   bool
}

// IllustDetail returns detail information for an illustration.
func (c *Client) IllustDetail(ctx context.Context, params IllustDetailParams) (*IllustDetailResponse, error) {
	v := url.Values{"illust_id": {strconv.Itoa(params.IllustID)}}
	resp, err := c.doGet(ctx, "/v1/illust/detail", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[IllustDetailResponse](resp)
}

// IllustCommentsParams are the parameters for IllustComments.
type IllustCommentsParams struct {
	IllustID             int
	Offset               *int
	IncludeTotalComments *bool
	NoAuth               bool
}

// IllustComments returns comments on an illustration.
func (c *Client) IllustComments(ctx context.Context, params IllustCommentsParams) (*IllustCommentsResponse, error) {
	v := url.Values{"illust_id": {strconv.Itoa(params.IllustID)}}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	if params.IncludeTotalComments != nil {
		v.Set("include_total_comments", formatBool(*params.IncludeTotalComments))
	}
	resp, err := c.doGet(ctx, "/v1/illust/comments", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[IllustCommentsResponse](resp)
}

// IllustRelatedParams are the parameters for IllustRelated.
type IllustRelatedParams struct {
	IllustID      int
	Filter        Filter
	SeedIllustIDs []string
	Offset        *int
	Viewed        []string
	NoAuth        bool
}

// IllustRelated returns illustrations related to the specified illustration.
func (c *Client) IllustRelated(ctx context.Context, params IllustRelatedParams) (*IllustListResponse, error) {
	v := url.Values{
		"illust_id": {strconv.Itoa(params.IllustID)},
		"filter":    {defaultFilter(params.Filter)},
	}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	for _, id := range params.SeedIllustIDs {
		v.Add("seed_illust_ids[]", id)
	}
	for _, id := range params.Viewed {
		v.Add("viewed[]", id)
	}
	resp, err := c.doGet(ctx, "/v2/illust/related", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[IllustListResponse](resp)
}

// IllustRecommendedParams are the parameters for IllustRecommended.
type IllustRecommendedParams struct {
	ContentType                  IllustType
	IncludeRankingLabel          *bool
	Filter                       Filter
	MaxBookmarkIDForRecommend    *int
	MinBookmarkIDForRecentIllust *int
	Offset                       *int
	IncludeRankingIllusts        *bool
	BookmarkIllustIDs            []string
	IncludePrivacyPolicy         *string
	Viewed                       []string
	NoAuth                       bool
}

// IllustRecommended returns recommended illustrations.
// When NoAuth is true, the no-login endpoint is used.
func (c *Client) IllustRecommended(ctx context.Context, params IllustRecommendedParams) (*IllustListResponse, error) {
	path := "/v1/illust/recommended"
	if params.NoAuth {
		path = "/v1/illust/recommended-nologin"
	}

	ct := string(params.ContentType)
	if ct == "" {
		ct = string(IllustTypeIllust)
	}

	includeRanking := "true"
	if params.IncludeRankingLabel != nil {
		includeRanking = formatBool(*params.IncludeRankingLabel)
	}

	v := url.Values{
		"content_type":          {ct},
		"include_ranking_label": {includeRanking},
		"filter":                {defaultFilter(params.Filter)},
	}
	if params.MaxBookmarkIDForRecommend != nil {
		v.Set("max_bookmark_id_for_recommend", strconv.Itoa(*params.MaxBookmarkIDForRecommend))
	}
	if params.MinBookmarkIDForRecentIllust != nil {
		v.Set("min_bookmark_id_for_recent_illust", strconv.Itoa(*params.MinBookmarkIDForRecentIllust))
	}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	if params.IncludeRankingIllusts != nil {
		v.Set("include_ranking_illusts", formatBool(*params.IncludeRankingIllusts))
	}
	for _, id := range params.Viewed {
		v.Add("viewed[]", id)
	}
	if params.NoAuth && len(params.BookmarkIllustIDs) > 0 {
		v.Set("bookmark_illust_ids", joinStrings(params.BookmarkIllustIDs))
	}
	if params.IncludePrivacyPolicy != nil {
		v.Set("include_privacy_policy", *params.IncludePrivacyPolicy)
	}

	resp, err := c.doGet(ctx, path, v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[IllustListResponse](resp)
}

// IllustRankingParams are the parameters for IllustRanking.
type IllustRankingParams struct {
	Mode   RankingMode
	Filter Filter
	Date   *string // "YYYY-MM-DD"
	Offset *int
	NoAuth bool
}

// IllustRanking returns illustration rankings.
func (c *Client) IllustRanking(ctx context.Context, params IllustRankingParams) (*IllustListResponse, error) {
	mode := string(params.Mode)
	if mode == "" {
		mode = string(ModeDay)
	}
	v := url.Values{
		"mode":   {mode},
		"filter": {defaultFilter(params.Filter)},
	}
	if params.Date != nil {
		v.Set("date", *params.Date)
	}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v1/illust/ranking", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[IllustListResponse](resp)
}

// IllustFollowParams are the parameters for IllustFollow.
type IllustFollowParams struct {
	Restrict Restrict
	Offset   *int
	NoAuth   bool
}

// IllustFollow returns new illustrations from followed users. Requires authentication.
func (c *Client) IllustFollow(ctx context.Context, params IllustFollowParams) (*IllustListResponse, error) {
	restrict := string(params.Restrict)
	if restrict == "" {
		restrict = string(RestrictPublic)
	}
	v := url.Values{"restrict": {restrict}}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v2/illust/follow", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[IllustListResponse](resp)
}

// IllustNewParams are the parameters for IllustNew.
type IllustNewParams struct {
	ContentType  IllustType
	Filter       Filter
	MaxIllustID  *int
	NoAuth       bool
}

// IllustNew returns the latest illustrations.
func (c *Client) IllustNew(ctx context.Context, params IllustNewParams) (*IllustListResponse, error) {
	ct := string(params.ContentType)
	if ct == "" {
		ct = string(IllustTypeIllust)
	}
	v := url.Values{
		"content_type": {ct},
		"filter":       {defaultFilter(params.Filter)},
	}
	if params.MaxIllustID != nil {
		v.Set("max_illust_id", strconv.Itoa(*params.MaxIllustID))
	}
	resp, err := c.doGet(ctx, "/v1/illust/new", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[IllustListResponse](resp)
}

// UgoiraMetadataParams are the parameters for UgoiraMetadata.
type UgoiraMetadataParams struct {
	IllustID int
	NoAuth   bool
}

// UgoiraMetadata returns metadata for an animated illustration (ugoira).
func (c *Client) UgoiraMetadata(ctx context.Context, params UgoiraMetadataParams) (*UgoiraMetadataResponse, error) {
	v := url.Values{"illust_id": {strconv.Itoa(params.IllustID)}}
	resp, err := c.doGet(ctx, "/v1/ugoira/metadata", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[UgoiraMetadataResponse](resp)
}
