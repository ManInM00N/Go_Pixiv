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
	Pages           int64     `db:"pages"`
	Likecount       int64     `db:"likecount"`
	UploadedTime    time.Time `db:"uploaded_time"`
}

func (i *Illust) msg() string {
	return strconv.FormatInt(i.Pid, 10) +
		"\n  " + i.PreviewImageUrl

}
