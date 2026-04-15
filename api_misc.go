package pixivgo

import (
	"context"
	"net/url"
	"strconv"
)

// TrendingTagsIllustParams are the parameters for TrendingTagsIllust.
type TrendingTagsIllustParams struct {
	Filter Filter
	NoAuth bool
}

// TrendingTagsIllust returns trending tags for illustrations.
func (c *Client) TrendingTagsIllust(ctx context.Context, params TrendingTagsIllustParams) (*TrendingTagsResponse, error) {
	v := url.Values{"filter": {defaultFilter(params.Filter)}}
	resp, err := c.doGet(ctx, "/v1/trending-tags/illust", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[TrendingTagsResponse](resp)
}

// ShowcaseArticleParams are the parameters for ShowcaseArticle.
type ShowcaseArticleParams struct {
	ShowcaseID int
}

// ShowcaseArticle returns a featured showcase article.
// This uses the Pixiv Web API (no authentication required).
func (c *Client) ShowcaseArticle(ctx context.Context, params ShowcaseArticleParams) (*ShowcaseArticleResponse, error) {
	articleURL := "https://www.pixiv.net/ajax/showcase/article?article_id=" + strconv.Itoa(params.ShowcaseID)

	v := url.Values{"article_id": {strconv.Itoa(params.ShowcaseID)}}

	// Override headers for Web API — use Chrome UA instead of iOS app UA
	origUA := c.userAgent
	c.userAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
	c.additionalHeaders.Set("Referer", "https://www.pixiv.net")
	defer func() {
		c.userAgent = origUA
		c.additionalHeaders.Del("Referer")
	}()

	_ = articleURL // the full URL is constructed via doGet with absolute path
	resp, err := c.doGet(ctx, "https://www.pixiv.net/ajax/showcase/article", v, true)
	if err != nil {
		return nil, err
	}
	return parseResponse[ShowcaseArticleResponse](resp)
}
