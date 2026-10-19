package freshrss

type authRequest struct {
	APIKey string `form:"api_key"`
}

type authCheckRequest struct {
	API    string `query:"api"`
	APIKey string `form:"api_key"`
}

type authCheckResponse struct {
	APIVersion          int    `json:"api_version"`
	Auth                int    `json:"auth"`
	LastRefreshedOnTime string `json:"last_refreshed_on_time"`
}

type apiRequest struct {
	API string `query:"api"`
}

type groupsRequest struct {
	API    string `query:"api"`
	Groups string `query:"groups"`
}

type feedsRequest struct {
	API   string `query:"api"`
	Feeds string `query:"feeds"`
}

type faviconsRequest struct {
	API      string `query:"api"`
	Favicons string `query:"favicons"`
}

type itemsRequest struct {
	API   string `query:"api"`
	Items string `query:"items"`

	SinceID string
	MaxID   string
	WithIDs string
}

type itemIDsRequest struct {
	API           string `query:"api"`
	UnreadItemIDs string `query:"unread_item_ids"`
	SavedItemIDs  string `query:"saved_item_ids"`
}

type markItemRequest struct {
	As string
	ID string
}

type markFeedRequest struct {
	ID     string
	Before string
}

type markGroupRequest struct {
	ID     string
	Before string
}

type unreadRecentlyReadRequest struct {
	API                string `query:"api"`
	UnreadRecentlyRead string `form:"unread_recently_read"`
}

type markItemForm struct {
	APIKey string `form:"api_key"`
	Mark   string `form:"mark"`
	As     string `form:"as"`
	ID     string `form:"id"`
}

type markFeedForm struct {
	APIKey string `form:"api_key"`
	Mark   string `form:"mark"`
	As     string `form:"as"`
	ID     string `form:"id"`
	Before string `form:"before"`
}

type markGroupForm markFeedForm

type unreadRecentlyReadForm struct {
	APIKey             string `form:"api_key"`
	UnreadRecentlyRead string `form:"unread_recently_read"`
}

type groupsResponse struct {
	Groups      []group      `json:"groups"`
	FeedsGroups []feedsGroup `json:"feeds_groups"`
}

type feedsResponse struct {
	Feeds       []feed       `json:"feeds"`
	FeedsGroups []feedsGroup `json:"feeds_groups"`
}

type faviconsResponse struct {
	Favicons []favicon `json:"favicons"`
}

type itemsResponse struct {
	Items      []item `json:"items"`
	TotalItems int    `json:"total_items"`
}

type itemIDsResponse struct {
	UnreadItemIDs string `json:"unread_item_ids"`
	SavedItemIDs  string `json:"saved_item_ids"`
}

type group struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type feed struct {
	ID                int    `json:"id"`
	FaviconID         int    `json:"favicon_id"`
	Title             string `json:"title"`
	URL               string `json:"url"`
	SiteURL           string `json:"site_url"`
	IsSpark           int    `json:"is_spark"`
	LastUpdatedOnTime int64  `json:"last_updated_on_time"`
}

type feedsGroup struct {
	GroupID int    `json:"group_id"`
	FeedIDs string `json:"feed_ids"`
}

type favicon struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

type item struct {
	ID            int    `json:"id"`
	FeedID        int    `json:"feed_id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	HTML          string `json:"html"`
	URL           string `json:"url"`
	IsSaved       int    `json:"is_saved"`
	IsRead        int    `json:"is_read"`
	CreatedOnTime int64  `json:"created_on_time"`
}
