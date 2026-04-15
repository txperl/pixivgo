package pixivgo

import (
	"context"
	"net/url"
	"strconv"
)

// UserDetailParams are the parameters for UserDetail.
type UserDetailParams struct {
	UserID int
	Filter Filter
	NoAuth bool
}

// UserDetail returns detailed information for a user.
func (c *Client) UserDetail(ctx context.Context, params UserDetailParams) (*UserInfoDetailed, error) {
	v := url.Values{
		"user_id": {strconv.Itoa(params.UserID)},
		"filter":  {defaultFilter(params.Filter)},
	}
	resp, err := c.doGet(ctx, "/v1/user/detail", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[UserInfoDetailed](resp)
}

// UserIllustsParams are the parameters for UserIllusts.
type UserIllustsParams struct {
	UserID int
	Type   IllustType
	Filter Filter
	Offset *int
	NoAuth bool
}

// UserIllusts returns a user's illustrations.
func (c *Client) UserIllusts(ctx context.Context, params UserIllustsParams) (*UserIllustrations, error) {
	v := url.Values{
		"user_id": {strconv.Itoa(params.UserID)},
		"filter":  {defaultFilter(params.Filter)},
	}
	if params.Type != "" {
		v.Set("type", string(params.Type))
	}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v1/user/illusts", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[UserIllustrations](resp)
}

// UserBookmarksIllustParams are the parameters for UserBookmarksIllust.
type UserBookmarksIllustParams struct {
	UserID        int
	Restrict      Restrict
	Filter        Filter
	MaxBookmarkID *int
	Tag           *string
	NoAuth        bool
}

// UserBookmarksIllust returns a user's bookmarked illustrations.
func (c *Client) UserBookmarksIllust(ctx context.Context, params UserBookmarksIllustParams) (*UserBookmarksIllustrations, error) {
	restrict := string(params.Restrict)
	if restrict == "" {
		restrict = string(RestrictPublic)
	}
	v := url.Values{
		"user_id":  {strconv.Itoa(params.UserID)},
		"restrict": {restrict},
		"filter":   {defaultFilter(params.Filter)},
	}
	if params.MaxBookmarkID != nil {
		v.Set("max_bookmark_id", strconv.Itoa(*params.MaxBookmarkID))
	}
	if params.Tag != nil {
		v.Set("tag", *params.Tag)
	}
	resp, err := c.doGet(ctx, "/v1/user/bookmarks/illust", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[UserBookmarksIllustrations](resp)
}

// UserBookmarksNovelParams are the parameters for UserBookmarksNovel.
type UserBookmarksNovelParams struct {
	UserID        int
	Restrict      Restrict
	Filter        Filter
	MaxBookmarkID *int
	Tag           *string
	NoAuth        bool
}

// UserBookmarksNovel returns a user's bookmarked novels.
func (c *Client) UserBookmarksNovel(ctx context.Context, params UserBookmarksNovelParams) (*UserBookmarksNovel, error) {
	restrict := string(params.Restrict)
	if restrict == "" {
		restrict = string(RestrictPublic)
	}
	v := url.Values{
		"user_id":  {strconv.Itoa(params.UserID)},
		"restrict": {restrict},
		"filter":   {defaultFilter(params.Filter)},
	}
	if params.MaxBookmarkID != nil {
		v.Set("max_bookmark_id", strconv.Itoa(*params.MaxBookmarkID))
	}
	if params.Tag != nil {
		v.Set("tag", *params.Tag)
	}
	resp, err := c.doGet(ctx, "/v1/user/bookmarks/novel", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[UserBookmarksNovel](resp)
}

// UserFollowingParams are the parameters for UserFollowing.
type UserFollowingParams struct {
	UserID   int
	Restrict Restrict
	Offset   *int
	NoAuth   bool
}

// UserFollowing returns the users that a user is following.
func (c *Client) UserFollowing(ctx context.Context, params UserFollowingParams) (*UserFollowing, error) {
	restrict := string(params.Restrict)
	if restrict == "" {
		restrict = string(RestrictPublic)
	}
	v := url.Values{
		"user_id":  {strconv.Itoa(params.UserID)},
		"restrict": {restrict},
	}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v1/user/following", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[UserFollowing](resp)
}

// UserFollowerParams are the parameters for UserFollower.
type UserFollowerParams struct {
	UserID int
	Filter Filter
	Offset *int
	NoAuth bool
}

// UserFollower returns a user's followers.
func (c *Client) UserFollower(ctx context.Context, params UserFollowerParams) (*UserListResponse, error) {
	v := url.Values{
		"user_id": {strconv.Itoa(params.UserID)},
		"filter":  {defaultFilter(params.Filter)},
	}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v1/user/follower", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[UserListResponse](resp)
}

// UserMyPixivParams are the parameters for UserMyPixiv.
type UserMyPixivParams struct {
	UserID int
	Offset *int
	NoAuth bool
}

// UserMyPixiv returns a user's "My Pixiv" (close friends).
func (c *Client) UserMyPixiv(ctx context.Context, params UserMyPixivParams) (*UserListResponse, error) {
	v := url.Values{"user_id": {strconv.Itoa(params.UserID)}}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v1/user/mypixiv", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[UserListResponse](resp)
}

// UserListParams are the parameters for UserList.
type UserListParams struct {
	UserID int
	Filter Filter
	Offset *int
	NoAuth bool
}

// UserList returns a user's blacklist.
func (c *Client) UserList(ctx context.Context, params UserListParams) (*UserListResponse, error) {
	v := url.Values{
		"user_id": {strconv.Itoa(params.UserID)},
		"filter":  {defaultFilter(params.Filter)},
	}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v2/user/list", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[UserListResponse](resp)
}

// UserRelatedParams are the parameters for UserRelated.
type UserRelatedParams struct {
	SeedUserID int
	Filter     Filter
	Offset     *int
	NoAuth     bool
}

// UserRelated returns users related to the specified user.
func (c *Client) UserRelated(ctx context.Context, params UserRelatedParams) (*UserListResponse, error) {
	v := url.Values{
		"filter":       {defaultFilter(params.Filter)},
		"offset":       {"0"},
		"seed_user_id": {strconv.Itoa(params.SeedUserID)},
	}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v1/user/related", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[UserListResponse](resp)
}

// UserRecommendedParams are the parameters for UserRecommended.
type UserRecommendedParams struct {
	Filter Filter
	Offset *int
	NoAuth bool
}

// UserRecommended returns recommended users.
func (c *Client) UserRecommended(ctx context.Context, params UserRecommendedParams) (*UserListResponse, error) {
	v := url.Values{"filter": {defaultFilter(params.Filter)}}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v1/user/recommended", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[UserListResponse](resp)
}

// UserNovelsParams are the parameters for UserNovels.
type UserNovelsParams struct {
	UserID int
	Filter Filter
	Offset *int
	NoAuth bool
}

// UserNovels returns a user's novels.
func (c *Client) UserNovels(ctx context.Context, params UserNovelsParams) (*UserNovels, error) {
	v := url.Values{
		"user_id": {strconv.Itoa(params.UserID)},
		"filter":  {defaultFilter(params.Filter)},
	}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v1/user/novels", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[UserNovels](resp)
}

// UserBookmarkTagsIllustParams are the parameters for UserBookmarkTagsIllust.
type UserBookmarkTagsIllustParams struct {
	UserID   int
	Restrict Restrict
	Offset   *int
	NoAuth   bool
}

// UserBookmarkTagsIllust returns a user's bookmark tags for illustrations.
func (c *Client) UserBookmarkTagsIllust(ctx context.Context, params UserBookmarkTagsIllustParams) (*UserBookmarkTagsResponse, error) {
	restrict := string(params.Restrict)
	if restrict == "" {
		restrict = string(RestrictPublic)
	}
	v := url.Values{
		"user_id":  {strconv.Itoa(params.UserID)},
		"restrict": {restrict},
	}
	if params.Offset != nil {
		v.Set("offset", strconv.Itoa(*params.Offset))
	}
	resp, err := c.doGet(ctx, "/v1/user/bookmark-tags/illust", v, params.NoAuth)
	if err != nil {
		return nil, err
	}
	return parseResponse[UserBookmarkTagsResponse](resp)
}

// UserFollowAddParams are the parameters for UserFollowAdd.
type UserFollowAddParams struct {
	UserID   int
	Restrict Restrict
	NoAuth   bool
}

// UserFollowAdd follows a user. Requires authentication.
func (c *Client) UserFollowAdd(ctx context.Context, params UserFollowAddParams) error {
	restrict := string(params.Restrict)
	if restrict == "" {
		restrict = string(RestrictPublic)
	}
	data := url.Values{
		"user_id":  {strconv.Itoa(params.UserID)},
		"restrict": {restrict},
	}
	resp, err := c.doPost(ctx, "/v1/user/follow/add", data, params.NoAuth)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// UserFollowDeleteParams are the parameters for UserFollowDelete.
type UserFollowDeleteParams struct {
	UserID int
	NoAuth bool
}

// UserFollowDelete unfollows a user. Requires authentication.
func (c *Client) UserFollowDelete(ctx context.Context, params UserFollowDeleteParams) error {
	data := url.Values{"user_id": {strconv.Itoa(params.UserID)}}
	resp, err := c.doPost(ctx, "/v1/user/follow/delete", data, params.NoAuth)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// UserEditAIShowSettingsParams are the parameters for UserEditAIShowSettings.
type UserEditAIShowSettingsParams struct {
	ShowAI bool
	NoAuth bool
}

// UserEditAIShowSettings toggles the visibility of AI-generated works. Requires authentication.
func (c *Client) UserEditAIShowSettings(ctx context.Context, params UserEditAIShowSettingsParams) error {
	data := url.Values{"show_ai": {formatBool(params.ShowAI)}}
	resp, err := c.doPost(ctx, "/v1/user/ai-show-settings/edit", data, params.NoAuth)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
