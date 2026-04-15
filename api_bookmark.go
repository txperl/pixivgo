package pixivgo

import (
	"context"
	"net/url"
	"strconv"
	"strings"
)

// IllustBookmarkDetailParams are the parameters for IllustBookmarkDetail.
type IllustBookmarkDetailParams struct {
	IllustID int
	NoAuth   bool
}

// IllustBookmarkDetail returns bookmark detail for an illustration.
func (c *Client) IllustBookmarkDetail(ctx context.Context, params IllustBookmarkDetailParams) (*BookmarkDetailResponse, error) {
	v := url.Values{"illust_id": {strconv.Itoa(params.IllustID)}}
	resp, err := c.doGet(ctx, "/v2/illust/bookmark/detail", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[BookmarkDetailResponse](resp)
}

// IllustBookmarkAddParams are the parameters for IllustBookmarkAdd.
type IllustBookmarkAddParams struct {
	IllustID int
	Restrict Restrict
	Tags     []string
	NoAuth   bool
}

// IllustBookmarkAdd adds an illustration to bookmarks. Requires authentication.
func (c *Client) IllustBookmarkAdd(ctx context.Context, params IllustBookmarkAddParams) error {
	restrict := string(params.Restrict)
	if restrict == "" {
		restrict = string(RestrictPublic)
	}
	data := url.Values{
		"illust_id": {strconv.Itoa(params.IllustID)},
		"restrict":  {restrict},
	}
	if len(params.Tags) > 0 {
		data.Set("tags[]", strings.Join(params.Tags, " "))
	}
	resp, err := c.doPost(ctx, "/v2/illust/bookmark/add", data, params.NoAuth)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// IllustBookmarkDeleteParams are the parameters for IllustBookmarkDelete.
type IllustBookmarkDeleteParams struct {
	IllustID int
	NoAuth   bool
}

// IllustBookmarkDelete removes an illustration from bookmarks. Requires authentication.
func (c *Client) IllustBookmarkDelete(ctx context.Context, params IllustBookmarkDeleteParams) error {
	data := url.Values{"illust_id": {strconv.Itoa(params.IllustID)}}
	resp, err := c.doPost(ctx, "/v1/illust/bookmark/delete", data, params.NoAuth)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
