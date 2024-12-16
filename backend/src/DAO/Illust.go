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

type Illust struct {
	Pid             int64     `db:"pid"`
	Title           string    `db:"title"`
	Caption         string    `db:"caption"`
	Tags            []string  `db:"tags"`
	ImageUrl        []string  `db:"image_url"`
	PreviewImageUrl string    `db:"preview_image"`
	AgeLimit        string    `db:"age_limit"`
	CreatedTime     string    `db:"created_time"`
	UserID          int64     `db:"userId"`
	UserName        string    `db:"user_name"`
	Pages           int       `db:"pages"`
	Likecount       int       `db:"likecount"`
	UploadedTime    time.Time `db:"uploaded_time"`
}

func (i *Illust) msg() string {
	return strconv.FormatInt(i.Pid, 10) +
		"\n  " + i.PreviewImageUrl
}

type FollowData struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	PreviewUrl string   `json:"url"`
	UserID     string   `json:"userId"`
	UserName   string   `json:"userName"`
	PageCount  int      `json:"pageCount"`
	AiType     int      `json:"aiType"`
	Tags       []string `json:"tags"`
	R18        bool     `json:"r18"`
}
type RankData struct {
	ID         int64    `json:"illust_id"`
	Title      string   `json:"title"`
	PreviewUrl string   `json:"url"`
	UserID     int64    `json:"user_id"`
	UserName   string   `json:"user_name"`
	PageCount  string   `json:"illust_page_count"`
	Tags       []string `json:"tags"`
	R18        bool     `json:"illust_content_type.sexual"`
	Type       string   `json:"illust_type"`
}
