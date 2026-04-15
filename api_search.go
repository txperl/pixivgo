package pixivgo

import (
	"context"
	"net/url"
	"strconv"
)

// SearchIllustParams are the parameters for SearchIllust.
type SearchIllustParams struct {
	Word         string       // required
	SearchTarget SearchTarget // default: PartialMatchForTags
	Sort         Sort         // default: SortDateDesc
	Duration     *Duration
	StartDate    *string // "YYYY-MM-DD"
	EndDate      *string // "YYYY-MM-DD"
	Filter       Filter
	SearchAIType *int // 0 or 1
	Offset       *int
	NoAuth       bool
}

// SearchIllust searches for illustrations by keyword.
func (c *Client) SearchIllust(ctx context.Context, params SearchIllustParams) (*SearchIllustrations, error) {
	searchTarget := string(params.SearchTarget)
	if searchTarget == "" {
		searchTarget = string(PartialMatchForTags)
	}
	sort := string(params.Sort)
	if sort == "" {
		sort = string(SortDateDesc)
	}
	v := url.Values{
		"word":          {params.Word},
		"search_target": {searchTarget},
		"sort":          {sort},
		"filter":        {defaultFilter(params.Filter)},
	}
	if params.StartDate != nil {
		v.Set("start_date", *params.StartDate)
	}
	if params.EndDate != nil {
		v.Set("end_date", *params.EndDate)
	}
	if params.Duration != nil {
		v.Set("duration", string(*params.Duration))
	}
	if params.SearchAIType != nil {
		v.Set("search_ai_type", strconv.Itoa(*params.SearchAIType))
	}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v1/search/illust", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[SearchIllustrations](resp)
}

// SearchNovelParams are the parameters for SearchNovel.
type SearchNovelParams struct {
	Word                          string
	SearchTarget                  SearchTarget
	Sort                          Sort
	MergePlainKeywordResults      *string // "true" or "false"
	IncludeTranslatedTagResults   *string // "true" or "false"
	StartDate                     *string
	EndDate                       *string
	Filter                        *string
	SearchAIType                  *int
	Offset                        *int
	NoAuth                        bool
}

// SearchNovel searches for novels by keyword.
func (c *Client) SearchNovel(ctx context.Context, params SearchNovelParams) (*SearchNovelResponse, error) {
	searchTarget := string(params.SearchTarget)
	if searchTarget == "" {
		searchTarget = string(PartialMatchForTags)
	}
	sort := string(params.Sort)
	if sort == "" {
		sort = string(SortDateDesc)
	}

	merge := "true"
	if params.MergePlainKeywordResults != nil {
		merge = *params.MergePlainKeywordResults
	}
	include := "true"
	if params.IncludeTranslatedTagResults != nil {
		include = *params.IncludeTranslatedTagResults
	}

	v := url.Values{
		"word":                             {params.Word},
		"search_target":                    {searchTarget},
		"sort":                             {sort},
		"merge_plain_keyword_results":      {merge},
		"include_translated_tag_results":   {include},
	}
	if params.Filter != nil {
		v.Set("filter", *params.Filter)
	}
	if params.StartDate != nil {
		v.Set("start_date", *params.StartDate)
	}
	if params.EndDate != nil {
		v.Set("end_date", *params.EndDate)
	}
	if params.SearchAIType != nil {
		v.Set("search_ai_type", strconv.Itoa(*params.SearchAIType))
	}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v1/search/novel", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[SearchNovelResponse](resp)
}

// SearchUserParams are the parameters for SearchUser.
type SearchUserParams struct {
	Word     string
	Sort     Sort
	Duration *Duration
	Filter   Filter
	Offset   *int
	NoAuth   bool
}

// SearchUser searches for users by keyword.
func (c *Client) SearchUser(ctx context.Context, params SearchUserParams) (*UserListResponse, error) {
	sort := string(params.Sort)
	if sort == "" {
		sort = string(SortDateDesc)
	}
	v := url.Values{
		"word":   {params.Word},
		"sort":   {sort},
		"filter": {defaultFilter(params.Filter)},
	}
	if params.Duration != nil {
		v.Set("duration", string(*params.Duration))
	}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v1/search/user", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[UserListResponse](resp)
}
