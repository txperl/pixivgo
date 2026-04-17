package pixivgo

import (
	"fmt"
	"strconv"
)

// FlexInt is an int that can be unmarshaled from both JSON numbers and strings.
// Pixiv's API is inconsistent: the auth endpoint returns user IDs as strings,
// while other endpoints return them as integers.
type FlexInt int

func (fi *FlexInt) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("FlexInt: empty input")
	}
	// JSON number: starts with digit or '-'
	if data[0] != '"' {
		n, err := strconv.Atoi(string(data))
		if err != nil {
			return fmt.Errorf("FlexInt: cannot parse %s as int: %w", string(data), err)
		}
		*fi = FlexInt(n)
		return nil
	}
	// JSON string: strip quotes then parse
	if len(data) < 2 || data[len(data)-1] != '"' {
		return fmt.Errorf("FlexInt: invalid string %s", string(data))
	}
	n, err := strconv.Atoi(string(data[1 : len(data)-1]))
	if err != nil {
		return fmt.Errorf("FlexInt: cannot parse %s as int: %w", string(data), err)
	}
	*fi = FlexInt(n)
	return nil
}

func (fi FlexInt) MarshalJSON() ([]byte, error) {
	return strconv.AppendInt(nil, int64(fi), 10), nil
}

// Int returns the underlying int value.
func (fi FlexInt) Int() int {
	return int(fi)
}

// --- Type constants (mapped from Python Literal types) ---

// Filter controls platform-specific filtering.
type Filter string

const (
	FilterForIOS Filter = "for_ios"
	FilterNone   Filter = "none"
)

// IllustType is the type of illustration.
type IllustType string

const (
	IllustTypeIllust IllustType = "illust"
	IllustTypeManga  IllustType = "manga"
)

// Restrict controls public/private visibility.
type Restrict string

const (
	RestrictPublic  Restrict = "public"
	RestrictPrivate Restrict = "private"
)

// RankingMode is the ranking period/category.
type RankingMode string

const (
	ModeDay         RankingMode = "day"
	ModeWeek        RankingMode = "week"
	ModeMonth       RankingMode = "month"
	ModeDayMale     RankingMode = "day_male"
	ModeDayFemale   RankingMode = "day_female"
	ModeWeekOriginal RankingMode = "week_original"
	ModeWeekRookie  RankingMode = "week_rookie"
	ModeDayManga    RankingMode = "day_manga"
	ModeDayR18      RankingMode = "day_r18"
	ModeDayMaleR18  RankingMode = "day_male_r18"
	ModeDayFemaleR18 RankingMode = "day_female_r18"
	ModeWeekR18     RankingMode = "week_r18"
	ModeWeekR18G    RankingMode = "week_r18g"
)

// SearchTarget controls how search queries are matched.
type SearchTarget string

const (
	PartialMatchForTags SearchTarget = "partial_match_for_tags"
	ExactMatchForTags   SearchTarget = "exact_match_for_tags"
	TitleAndCaption     SearchTarget = "title_and_caption"
	Keyword             SearchTarget = "keyword"
)

// Sort controls the ordering of search results.
type Sort string

const (
	SortDateDesc    Sort = "date_desc"
	SortDateAsc     Sort = "date_asc"
	SortPopularDesc Sort = "popular_desc"
)

// Duration controls the time range for search results.
type Duration string

const (
	WithinLastDay   Duration = "within_last_day"
	WithinLastWeek  Duration = "within_last_week"
	WithinLastMonth Duration = "within_last_month"
)

// --- Response models (mapped from Python Pydantic models) ---

// ProfileImageUrls contains user profile image URLs.
type ProfileImageUrls struct {
	Medium string `json:"medium"`
}

// UserInfo contains basic user information.
type UserInfo struct {
	ID                   FlexInt          `json:"id"`
	Name                 string           `json:"name"`
	Account              string           `json:"account"`
	ProfileImageUrls     ProfileImageUrls `json:"profile_image_urls"`
	Comment              *string          `json:"comment,omitempty"`
	IsFollowed           *bool            `json:"is_followed"`
	IsAccessBlockingUser *bool            `json:"is_access_blocking_user,omitempty"`
	IsAcceptRequest      *bool            `json:"is_accept_request,omitempty"`
}

// CommentUser contains user info as it appears in comments.
type CommentUser struct {
	ID               FlexInt          `json:"id"`
	Name             string           `json:"name"`
	Account          string           `json:"account"`
	ProfileImageUrls ProfileImageUrls `json:"profile_image_urls"`
}

// Profile contains extended user profile details.
type Profile struct {
	Webpage                    *string `json:"webpage"`
	Gender                     string  `json:"gender"`
	Birth                      string  `json:"birth"`
	BirthDay                   string  `json:"birth_day"`
	BirthYear                  int     `json:"birth_year"`
	Region                     string  `json:"region"`
	AddressID                  int     `json:"address_id"`
	CountryCode                string  `json:"country_code"`
	Job                        string  `json:"job"`
	JobID                      int     `json:"job_id"`
	TotalFollowUsers           int     `json:"total_follow_users"`
	TotalMyPixivUsers          int     `json:"total_mypixiv_users"`
	TotalIllusts               int     `json:"total_illusts"`
	TotalManga                 int     `json:"total_manga"`
	TotalNovels                int     `json:"total_novels"`
	TotalIllustBookmarksPublic int     `json:"total_illust_bookmarks_public"`
	TotalIllustSeries          int     `json:"total_illust_series"`
	TotalNovelSeries           int     `json:"total_novel_series"`
	BackgroundImageURL         string  `json:"background_image_url"`
	TwitterAccount             string  `json:"twitter_account"`
	TwitterURL                 *string `json:"twitter_url"`
	PawooURL                   *string `json:"pawoo_url"`
	IsPremium                  bool    `json:"is_premium"`
	IsUsingCustomProfileImage  bool    `json:"is_using_custom_profile_image"`
}

// ProfilePublicity contains privacy settings.
type ProfilePublicity struct {
	Gender   string `json:"gender"`
	Region   string `json:"region"`
	BirthDay string `json:"birth_day"`
	BirthYear string `json:"birth_year"`
	Job      string `json:"job"`
	Pawoo    bool   `json:"pawoo"`
}

// Workspace contains user workspace information.
type Workspace struct {
	PC                string  `json:"pc"`
	Monitor           string  `json:"monitor"`
	Tool              string  `json:"tool"`
	Scanner           string  `json:"scanner"`
	Tablet            string  `json:"tablet"`
	Mouse             string  `json:"mouse"`
	Printer           string  `json:"printer"`
	Desktop           string  `json:"desktop"`
	Music             string  `json:"music"`
	Desk              string  `json:"desk"`
	Chair             string  `json:"chair"`
	Comment           string  `json:"comment"`
	WorkspaceImageURL *string `json:"workspace_image_url"`
}

// ImageUrls contains illustration image URLs at various sizes.
type ImageUrls struct {
	SquareMedium string `json:"square_medium"`
	Medium       string `json:"medium"`
	Large        string `json:"large"`
}

// IllustrationTag is a tag on an illustration.
type IllustrationTag struct {
	Name           string  `json:"name"`
	TranslatedName *string `json:"translated_name"`
}

// NovelTag is a tag on a novel.
type NovelTag struct {
	Name                string  `json:"name"`
	TranslatedName      *string `json:"translated_name"`
	AddedByUploadedUser bool    `json:"added_by_uploaded_user"`
}

// Series contains series information.
// Pixiv returns {} instead of null for empty series; in that case ID will be 0.
type Series struct {
	ID    FlexInt `json:"id"`
	Title string `json:"title"`
}

// MetaSinglePage contains the original image URL for single-page illustrations.
type MetaSinglePage struct {
	OriginalImageURL *string `json:"original_image_url,omitempty"`
}

// MetaPage contains image URLs for multi-page illustrations.
type MetaPage struct {
	ImageUrls ImageUrls `json:"image_urls"`
}

// IllustrationInfo contains complete illustration details.
type IllustrationInfo struct {
	ID                    FlexInt           `json:"id"`
	Title                 string            `json:"title"`
	Type                  string            `json:"type"`
	ImageUrls             ImageUrls         `json:"image_urls"`
	Caption               string            `json:"caption"`
	Restrict              int               `json:"restrict"`
	User                  UserInfo          `json:"user"`
	Tags                  []IllustrationTag `json:"tags"`
	Tools                 []string          `json:"tools"`
	CreateDate            string            `json:"create_date"`
	PageCount             int               `json:"page_count"`
	Width                 int               `json:"width"`
	Height                int               `json:"height"`
	SanityLevel           int               `json:"sanity_level"`
	XRestrict             int               `json:"x_restrict"`
	Series                *Series           `json:"series"`
	MetaSinglePage        MetaSinglePage    `json:"meta_single_page"`
	MetaPages             []MetaPage        `json:"meta_pages"`
	TotalView             int               `json:"total_view"`
	TotalBookmarks        int               `json:"total_bookmarks"`
	IsBookmarked          bool              `json:"is_bookmarked"`
	Visible               bool              `json:"visible"`
	IsMuted               bool              `json:"is_muted"`
	IllustAIType          int               `json:"illust_ai_type"`
	IllustBookStyle       int               `json:"illust_book_style"`
	TotalComments         *int              `json:"total_comments,omitempty"`
	RestrictionAttributes []string          `json:"restriction_attributes"`
}

// Comment represents a comment on an illustration or novel.
type Comment struct {
	ID            FlexInt      `json:"id"`
	Comment       string       `json:"comment"`
	Date          string       `json:"date"`
	User          *CommentUser `json:"user"`
	ParentComment *Comment     `json:"parent_comment,omitempty"`
}

// NovelInfo contains complete novel details.
type NovelInfo struct {
	ID                   FlexInt    `json:"id"`
	Title                string     `json:"title"`
	Caption              string     `json:"caption"`
	Restrict             int        `json:"restrict"`
	XRestrict            int        `json:"x_restrict"`
	IsOriginal           bool       `json:"is_original"`
	ImageUrls            ImageUrls  `json:"image_urls"`
	CreateDate           string     `json:"create_date"`
	Tags                 []NovelTag `json:"tags"`
	PageCount            int        `json:"page_count"`
	TextLength           int        `json:"text_length"`
	User                 UserInfo   `json:"user"`
	Series               *Series    `json:"series"`
	IsBookmarked         bool       `json:"is_bookmarked"`
	TotalBookmarks       int        `json:"total_bookmarks"`
	TotalView            int        `json:"total_view"`
	Visible              bool       `json:"visible"`
	TotalComments        int        `json:"total_comments"`
	IsMuted              bool       `json:"is_muted"`
	IsMyPixivOnly        bool       `json:"is_mypixiv_only"`
	IsXRestricted        bool       `json:"is_x_restricted"`
	NovelAIType          int        `json:"novel_ai_type"`
	CommentAccessControl *int       `json:"comment_access_control,omitempty"`
}

// NovelNavigationInfo contains navigation info within a novel series.
type NovelNavigationInfo struct {
	ID               FlexInt `json:"id"`
	Viewable         bool    `json:"viewable"`
	ContentOrder     string  `json:"content_order"`
	Title            string  `json:"title"`
	CoverURL         string  `json:"cover_url"`
	ViewableMessage  *string `json:"viewable_message"`
}

// NovelRating contains like/bookmark/view counts for a novel.
type NovelRating struct {
	Like     int `json:"like"`
	Bookmark int `json:"bookmark"`
	View     int `json:"view"`
}

// WebviewNovel contains novel content extracted from webview HTML.
// Fields use camelCase JSON tags matching the embedded JSON format.
type WebviewNovel struct {
	ID                   string               `json:"id"`
	Title                string               `json:"title"`
	SeriesID             *string              `json:"seriesId"`
	SeriesTitle          *string              `json:"seriesTitle"`
	SeriesIsWatched      *bool                `json:"seriesIsWatched"`
	UserID               string               `json:"userId"`
	CoverURL             string               `json:"coverUrl"`
	Tags                 []string             `json:"tags"`
	Caption              string               `json:"caption"`
	CDate                string               `json:"cdate"`
	Rating               NovelRating          `json:"rating"`
	Text                 string               `json:"text"`
	Marker               *string              `json:"marker"`
	Illusts              []string             `json:"illusts"`
	Images               []string             `json:"images"`
	SeriesNavigation     *NovelNavigationInfo `json:"seriesNavigation"`
	GlossaryItems        []string             `json:"glossaryItems"`
	ReplaceableItemIDs   []string             `json:"replaceableItemIds"`
	AIType               int                  `json:"aiType"`
	IsOriginal           bool                 `json:"isOriginal"`
}

// UserPreview contains a user with sample works.
type UserPreview struct {
	User    UserInfo           `json:"user"`
	Illusts []IllustrationInfo `json:"illusts"`
	Novels  []NovelInfo        `json:"novels"`
	IsMuted bool               `json:"is_muted"`
}

// --- Response wrapper types for API endpoints ---

// AuthResponse contains the token data from an auth request.
type AuthResponse struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	TokenType    string   `json:"token_type"`
	Scope        string   `json:"scope"`
	RefreshToken string   `json:"refresh_token"`
	User         UserInfo `json:"user"`
}

// authTokenResponse is the top-level wrapper for the auth endpoint.
type authTokenResponse struct {
	Response AuthResponse `json:"response"`
}

// UserInfoDetailed is the response from user_detail.
type UserInfoDetailed struct {
	User             UserInfo         `json:"user"`
	Profile          Profile          `json:"profile"`
	ProfilePublicity ProfilePublicity `json:"profile_publicity"`
	Workspace        Workspace        `json:"workspace"`
}

// UserIllustrations is the response from user_illusts.
type UserIllustrations struct {
	User    UserInfo           `json:"user"`
	Illusts []IllustrationInfo `json:"illusts"`
	NextURL *string            `json:"next_url"`
}

// UserFollowing is the response from user_following.
type UserFollowing struct {
	UserPreviews []UserPreview `json:"user_previews"`
	NextURL      *string       `json:"next_url"`
}

// UserBookmarksIllustrations is the response from user_bookmarks_illust.
type UserBookmarksIllustrations struct {
	Illusts []IllustrationInfo `json:"illusts"`
	NextURL *string            `json:"next_url"`
}

// UserBookmarksNovel is the response from user_bookmarks_novel.
type UserBookmarksNovel struct {
	Novels  []NovelInfo `json:"novels"`
	NextURL *string     `json:"next_url"`
}

// UserNovels is the response from user_novels.
type UserNovels struct {
	User    UserInfo    `json:"user"`
	Novels  []NovelInfo `json:"novels"`
	NextURL *string     `json:"next_url"`
}

// SearchIllustrations is the response from search_illust.
type SearchIllustrations struct {
	Illusts       []IllustrationInfo `json:"illusts"`
	NextURL       *string            `json:"next_url"`
	SearchSpanLimit int              `json:"search_span_limit"`
	ShowAI        bool               `json:"show_ai"`
}

// SearchNovelResponse is the response from search_novel.
type SearchNovelResponse struct {
	Novels          []NovelInfo `json:"novels"`
	NextURL         *string     `json:"next_url"`
	SearchSpanLimit int         `json:"search_span_limit"`
	ShowAI          bool        `json:"show_ai"`
}

// NovelComments is the response from novel_comments.
type NovelComments struct {
	TotalComments        int       `json:"total_comments"`
	Comments             []Comment `json:"comments"`
	NextURL              *string   `json:"next_url"`
	CommentAccessControl int       `json:"comment_access_control"`
}

// IllustDetailResponse is the response from illust_detail.
type IllustDetailResponse struct {
	Illust IllustrationInfo `json:"illust"`
}

// IllustCommentsResponse is the response from illust_comments.
type IllustCommentsResponse struct {
	TotalComments        int       `json:"total_comments"`
	Comments             []Comment `json:"comments"`
	NextURL              *string   `json:"next_url"`
	CommentAccessControl int       `json:"comment_access_control"`
}

// IllustListResponse is a generic response containing a list of illustrations.
// Used by illust_related, illust_recommended, illust_ranking, illust_follow, illust_new.
type IllustListResponse struct {
	Illusts []IllustrationInfo `json:"illusts"`
	NextURL *string            `json:"next_url"`
}

// UserListResponse is a generic response containing user previews.
// Used by user_follower, user_mypixiv, user_list, user_related, user_recommended.
type UserListResponse struct {
	UserPreviews []UserPreview `json:"user_previews"`
	NextURL      *string       `json:"next_url"`
}

// NovelDetailResponse wraps the novel_detail endpoint response.
type NovelDetailResponse struct {
	Novel NovelInfo `json:"novel"`
}

// NovelListResponse is a generic response containing a list of novels.
// Used by novel_recommended, novel_new, novel_follow.
type NovelListResponse struct {
	Novels  []NovelInfo `json:"novels"`
	NextURL *string     `json:"next_url"`
}

// NovelSeriesResponse is the response from novel_series.
type NovelSeriesResponse struct {
	NovelSeriesDetail any `json:"novel_series_detail"`
	Novels            []NovelInfo `json:"novels"`
	NextURL           *string     `json:"next_url"`
}

// UgoiraMetadataResponse is the response from ugoira_metadata.
type UgoiraMetadataResponse struct {
	UgoiraMetadata UgoiraMetadata `json:"ugoira_metadata"`
}

// UgoiraMetadata contains animated illustration metadata.
type UgoiraMetadata struct {
	ZipUrls UgoiraZipUrls  `json:"zip_urls"`
	Frames  []UgoiraFrame  `json:"frames"`
}

// UgoiraZipUrls contains URLs for ugoira zip files.
type UgoiraZipUrls struct {
	Medium string `json:"medium"`
}

// UgoiraFrame contains timing information for a single ugoira frame.
type UgoiraFrame struct {
	File  string `json:"file"`
	Delay int    `json:"delay"`
}

// TrendingTagsResponse is the response from trending_tags_illust.
type TrendingTagsResponse struct {
	TrendTags []TrendTag `json:"trend_tags"`
}

// TrendTag is a trending tag with a sample illustration.
type TrendTag struct {
	Tag              IllustrationTag  `json:"tag"`
	TranslatedName   *string          `json:"translated_name"`
	Illust           IllustrationInfo `json:"illust"`
}

// BookmarkDetailResponse is the response from illust_bookmark_detail.
type BookmarkDetailResponse struct {
	BookmarkDetail BookmarkDetail `json:"bookmark_detail"`
}

// BookmarkDetail contains bookmark information.
type BookmarkDetail struct {
	IsBookmarked bool             `json:"is_bookmarked"`
	Tags         []BookmarkTag    `json:"tags"`
	Restrict     string           `json:"restrict"`
}

// BookmarkTag contains a bookmark tag with registration status.
type BookmarkTag struct {
	Name         string `json:"name"`
	IsRegistered bool   `json:"is_registered"`
}

// UserBookmarkTagsResponse is the response from user_bookmark_tags_illust.
type UserBookmarkTagsResponse struct {
	BookmarkTags []UserBookmarkTag `json:"bookmark_tags"`
	NextURL      *string           `json:"next_url"`
}

// UserBookmarkTag contains a bookmark tag with count.
type UserBookmarkTag struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// ShowcaseArticleResponse is the response from showcase_article.
type ShowcaseArticleResponse struct {
	Body any `json:"body"`
}
