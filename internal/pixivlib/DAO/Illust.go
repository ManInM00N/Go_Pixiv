package DAO

import (
	"strconv"
	"time"
)

type ImageData struct {
	URLs struct {
		ThumbMini string `json:"thumb_mini"`
		Small     string `json:"small"`
		Regular   string `json:"regular"`
		Original  string `json:"original"`
	} `json:"urls"`
	Width  int `json:"width"`
	Height int `json:"height"`
}
type Frame struct {
	File  string `json:"file"`
	Delay int    `json:"delay"`
}

const (
	IllustType = iota
	MangaType
	UgoiraType
)

type Novel struct {
	Id          string
	Content     string
	UserId      string
	UserName    string
	Tags        []string
	CoverUrl    string
	SeriesId    int
	SeriesTitle string
	Title       string
	R18         bool
}
type Illust struct {
	Pid             int64
	Title           string
	Caption         string
	Tags            []string
	ImageUrl        []string
	Source          string
	Frames          []Frame
	FileType        string
	PreviewImageUrl string
	AgeLimit        string
	CreatedTime     string
	UserID          int64
	UserName        string
	Pages           int
	Likecount       int
	UploadedTime    time.Time
	IllustType      int
	Callback        func(name string, data ...interface{})
	Width           int64
	Height          int64
}

func (i *Illust) msg() string {
	return strconv.FormatInt(i.Pid, 10) +
		"\n  " + i.PreviewImageUrl
}

type FollowData struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	UserID          string   `json:"userId"`
	UserName        string   `json:"userName"`
	PageCount       int      `json:"pageCount"`
	AiType          int      `json:"aiType"`
	Tags            []string `json:"tags"`
	R18             bool     `json:"r18"`
	ProfileImageUrl string   `json:"profileImageUrl"`
	PreviewUrl      string   `json:"url"`

	// illust

	// novel
	Genre         string `json:"genre"`
	Description   string `json:"description"`
	SeriesId      string `json:"seriesId"`
	SeriesTitle   string `json:"seriesTitle"`
	TextCount     int    `json:"textCount"`
	WordCount     int    `json:"wordCount"`
	BookMarkCount int    `json:"bookmarkCount"`
}

type RankData struct {
	ID                int64    `json:"illust_id"`
	Title             string   `json:"title"`
	PreviewUrl        string   `json:"url"`
	UserID            int64    `json:"user_id"`
	Date              string   `json:"date"`
	UserName          string   `json:"user_name"`
	PageCount         string   `json:"illust_page_count"`
	Tags              []string `json:"tags"`
	IllustContentType struct {
		Sexual     int  `json:"sexual"`
		Lo         bool `json:"lo"`
		Grotesque  bool `json:"grotesque"`
		Violent    bool `json:"violent"`
		Homosexual bool `json:"homosexual"`
		Drug       bool `json:"drug"`
		Thoughts   bool `json:"thoughts"`
		Antisocial bool `json:"antisocial"`
		Religion   bool `json:"religion"`
		Original   bool `json:"original"`
		Furry      bool `json:"furry"`
		BL         bool `json:"bl"`
		Yuri       bool `json:"yuri"`
	} `json:"illust_content_type"`
	Type string `json:"illust_type"`
}
