package DAO

import (
	"fmt"
	"os"
	"sync"
)

type Settings struct {
	Prefix           string `yml:"prefix" json:"prefix"`
	Proxy            string `yml:"proxy" json:"proxy"`
	Cookie           string `yml:"cookie" json:"cookie"`
	Agelimit         bool   `yml:"r_18" json:"r_18"`
	Downloadposition string `yml:"downloadposition" json:"downloadposition"`
	LikeLimit        int    `yml:"likelimit" json:"likelimit"`
	Retry429         int    `yml:"retry429" json:"retry429"`
	Downloadinterval int    `yml:"downloadinterval" json:"downloadinterval"`
	Retryinterval    int    `yml:"retryinterval"json:"retryinterval"`
	DifferAuthor     bool   `yml:"differauthor" json:"differauthor"`
	ExpiredTime      int    `yml:"expiretime" json:"expired_time"`
}

var rwl sync.RWMutex

func (s *Settings) MsgDetail() string {
	str := fmt.Sprintf("\n\tProxy :%v\n\tCookie :%v\n\tAgelimit :%v\n\tDownloadposition :%v\n", s.Proxy, s.Cookie, s.Agelimit, s.Downloadposition)
	return str
}

func (s *Settings) UpdateSettings(NewSetting Settings) {
	rwl.Lock()
	*s = NewSetting
	s.LikeLimit = max(s.LikeLimit, 0)
	_, err := os.Stat(s.Downloadposition)
	if err != nil {
		s.Downloadposition = "Download"
	}
	s.Retry429 = max(s.Retry429, 5000)
	s.Retryinterval = max(s.Retryinterval, 400)
	s.Downloadinterval = max(s.Downloadinterval, 400)
	s.ExpiredTime = max(1, min(366, s.ExpiredTime))
	rwl.Unlock()
}
