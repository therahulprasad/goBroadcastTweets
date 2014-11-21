package main

import (
	"sync"
)

type TwitterError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TwitterResponseError struct {
	errors []TwitterError `json:"errors"`
}
type User struct {
	ContributorsEnabled            bool   `json:"contributors_enabled"`
	CreatedAt                      string `json:"created_at"`
	DefaultProfile                 bool   `json:"default_profile"`
	DefaultProfileImage            bool   `json:"default_profile_image"`
	Description                    string `json:"description"`
	FavouritesCount                int    `json:"favourites_count"`
	FollowRequestSent              bool   `json:"follow_request_sent"`
	FollowersCount                 int    `json:"followers_count"`
	Following                      bool   `json:"following"`
	FriendsCount                   int    `json:"friends_count"`
	GeoEnabled                     bool   `json:"geo_enabled"`
	Id                             int64  `json:"id"`
	IdStr                          string `json:"id_str"`
	IsTranslator                   bool   `json:"is_translator"`
	Lang                           string `json:"lang"`
	ListedCount                    int64  `json:"listed_count"`
	Location                       string `json:"location"`
	Name                           string `json:"name"`
	Notifications                  bool   `json:"notifications"`
	ProfileBackgroundColor         string `json:"profile_background_color"`
	ProfileBackgroundImageURL      string `json:"profile_background_image_url"`
	ProfileBackgroundImageUrlHttps string `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool   `json:"profile_background_tile"`
	ProfileImageURL                string `json:"profile_image_url"`
	ProfileImageUrlHttps           string `json:"profile_image_url_https"`
	ProfileLinkColor               string `json:"profile_link_color"`
	ProfileSidebarBorderColor      string `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool   `json:"profile_use_background_image"`
	Protected                      bool   `json:"protected"`
	ScreenName                     string `json:"screen_name"`
	ShowAllInlineMedia             bool   `json:"show_all_inline_media"`
	Status                         *Tweet `json:"status"` // Only included if the user is a friend
	StatusesCount                  int64  `json:"statuses_count"`
	TimeZone                       string `json:"time_zone"`
	URL                            string `json:"url"`
	UtcOffset                      int    `json:"utc_offset"`
	Verified                       bool   `json:"verified"`
}

type XUserList struct {
	mu       sync.Mutex
	userlist map[int64]User
}

func (obj *XUserList) Init() {
	obj.userlist = make(map[int64]User)
}
func (obj *XUserList) Set(userid int64, value User) {
	obj.mu.Lock()
	obj.userlist[userid] = value
	obj.mu.Unlock()
}
func (obj *XUserList) Get(userid int64) (user User) {
	obj.mu.Lock()
	user = obj.userlist[userid]
	obj.mu.Unlock()
	return
}

type XReplyStatuses struct {
	repliedList map[int64]bool
	mu          sync.Mutex
}

func (c *XReplyStatuses) Init() {
	c.repliedList = make(map[int64]bool)
}

func (c *XReplyStatuses) Initiate(userid int64) {
	c.mu.Lock()
	c.repliedList[userid] = false
	c.mu.Unlock()
}
func (c *XReplyStatuses) IsSet(userid int64) (ok bool) {
	c.mu.Lock()
	_, ok = c.repliedList[userid]
	c.mu.Unlock()
	return
}
func (c *XReplyStatuses) ListUseridUnsent() (list map[int64]int64) {
	list = make(map[int64]int64)
	c.mu.Lock()
	j := int64(0)
	for userid := range c.repliedList {
		if c.repliedList[userid] == false {
			list[j] = userid
			j++
		}
	}
	c.mu.Unlock()
	return
}
func (c *XReplyStatuses) Sent(userid int64) {
	c.mu.Lock()
	c.repliedList[userid] = true
	c.mu.Unlock()
}

type List struct {
	CreatedAt       string `json:"created_at"`
	Slug            string `json:"slug"`
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	Description     string `json:"description"`
	Mode            string `json:"mode"`
	Following       bool   `json:"following"`
	User            User   `json:"user"`
	MemberCount     int    `json:"member_count"`
	IdStr           string `json:"id_str"`
	SubscriberCount int    `json:"subscriber_count"`
	Id              int64  `json:"id"`
	Uri             string `json:"uri"`
}

type Place struct {
	Attributes  map[string]string `json:"attributes"`
	BoundingBox struct {
		Coordinates [][][]float64 `json:"coordinates"`
		Type        string        `json:"type"`
	} `json:"bounding_box"`
	ContainedWithin []struct {
		Attributes  map[string]string `json:"attributes"`
		BoundingBox struct {
			Coordinates [][][]float64 `json:"coordinates"`
			Type        string        `json:"type"`
		} `json:"bounding_box"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
		FullName    string `json:"full_name"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		PlaceType   string `json:"place_type"`
		URL         string `json:"url"`
	} `json:"contained_within"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	FullName    string `json:"full_name"`
	Geometry    struct {
		Coordinates [][][]float64 `json:"coordinates"`
		Type        string        `json:"type"`
	} `json:"geometry"`
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	PlaceType string   `json:"place_type"`
	Polylines []string `json:"polylines"`
	URL       string   `json:"url"`
}
type Entities struct {
	Hashtags []struct {
		Indices []int
		Text    string
	}
	Urls []struct {
		Indices      []int
		Url          string
		Display_url  string
		Expanded_url string
	}
	User_mentions []struct {
		Name        string
		Indices     []int
		Screen_name string
		Id          int64
		Id_str      string
	}
	Media []struct {
		Id              int64
		Id_str          string
		Media_url       string
		Media_url_https string
		Url             string
		Display_url     string
		Expanded_url    string
		Sizes           MediaSizes
		Type            string
		Indices         []int
	}
}

type MediaSizes struct {
	Medium MediaSize
	Thumb  MediaSize
	Small  MediaSize
	Large  MediaSize
}

type MediaSize struct {
	W      int
	H      int
	Resize string
}
type Tweet struct {
	Contributors         []int64     `json:"contributors"`
	Coordinates          interface{} `json:"coordinates"`
	CreatedAt            string      `json:"created_at"`
	Entities             Entities    `json:"entities"`
	FavoriteCount        int         `json:"favorite_count"`
	Favorited            bool        `json:"favorited"`
	Geo                  interface{} `json:"geo"`
	Id                   int64       `json:"id"`
	IdStr                string      `json:"id_str"`
	InReplyToScreenName  string      `json:"in_reply_to_screen_name"`
	InReplyToStatusID    int64       `json:"in_reply_to_status_id"`
	InReplyToStatusIdStr string      `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64       `json:"in_reply_to_user_id"`
	InReplyToUserIdStr   string      `json:"in_reply_to_user_id_str"`
	Place                Place       `json:"place"`
	PossiblySensitive    bool        `json:"possibly_sensitive"`
	RetweetCount         int         `json:"retweet_count"`
	Retweeted            bool        `json:"retweeted"`
	RetweetedStatus      *Tweet      `json:"retweeted_status"`
	Source               string      `json:"source"`
	Text                 string      `json:"text"`
	Truncated            bool        `json:"truncated"`
	User                 User        `json:"user"`
}

type searchResponse struct {
	Statuses []Tweet
}
