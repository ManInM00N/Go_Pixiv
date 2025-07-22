package configs

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"sync"
)

var Setting *Settings
var rwl sync.RWMutex

type Settings struct {
	Prefix           string `yml:"prefix" json:"prefix"`
	Proxy            string `yml:"proxy" json:"proxy"`
	Cookie           string `yml:"cookie" json:"cookie"`
	Agelimit         bool   `yml:"r_18" json:"r_18"`
	Downloadposition string `yml:"downloadposition" json:"downloadposition"`
	LikeLimit        int    `yml:"likelimit" json:"likelimit"`
	Retry429         int    `yml:"retry429" json:"retry429"`
	Downloadinterval int    `yml:"downloadinterval" json:"downloadinterval"`
	Retryinterval    int    `yml:"retryinterval" json:"retryinterval"`
	DifferAuthor     bool   `yml:"differauthor" json:"differauthor"`
	ExpiredTime      int    `yml:"expiretime" json:"expired_time"`
	UseProxy         bool   `yml:"useProxy" json:"useproxy"`
}

func (s *Settings) MsgDetail() string {
	res, _ := json.Marshal(s)
	return string(res)
}
func NowSetting() Settings {
	return *Setting
}
func (s *Settings) UpdateSettings(NewSetting Settings) {
	rwl.Lock()
	defer rwl.Unlock()
	tmp := *s
	tmp.UseProxy = s.UseProxy
	tmp.LikeLimit = max(NewSetting.LikeLimit, 0)
	_, err := os.Stat(NewSetting.Downloadposition)
	if err != nil {
		tmp.Downloadposition = "Download"
	}
	tmp.DifferAuthor = NewSetting.DifferAuthor
	tmp.Agelimit = NewSetting.Agelimit
	tmp.LikeLimit = NewSetting.LikeLimit
	tmp.Cookie = NewSetting.Cookie
	tmp.Retry429 = max(NewSetting.Retry429, 5000)
	tmp.Retryinterval = max(NewSetting.Retryinterval, 400)
	tmp.Downloadinterval = max(NewSetting.Downloadinterval, 400)
	tmp.ExpiredTime = max(7, min(366, NewSetting.ExpiredTime))
	*s = tmp
}

func UpdateSettings() {
	out, _ := yaml.Marshal(Setting)
	ioutil.WriteFile("configs/settings.yml", out, 0644)
}
