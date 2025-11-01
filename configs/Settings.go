package configs

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"sync"
)

var setting *Settings
var rwl sync.RWMutex

type Settings struct {
	Prefix   string `yml:"prefix" json:"prefix" default:"http://127.0.0.1:"`
	Proxy    string `yml:"proxy" json:"proxy"`
	UseProxy bool   `yml:"useProxy" json:"useproxy"`

	PixivConf PixivConf `yml:"pixivConf" json:"pixivConf"`

	SearchEngine SearchEngine `yml:"imageEngine" json:"imageEngine"`

	// droped
	//Cookie           string `yml:"cookie" json:"cookie"`
	//Agelimit         bool   `yml:"r_18" json:"r_18"`
	//Downloadposition string `yml:"downloadposition" json:"downloadposition"`
	//LikeLimit        int    `yml:"likelimit" json:"likelimit"`
	//Retry429         int    `yml:"retry429" json:"retry429"`
	//Downloadinterval int    `yml:"downloadinterval" json:"downloadinterval"`
	//Retryinterval    int    `yml:"retryinterval" json:"retryinterval"`
	//DifferAuthor     bool   `yml:"differauthor" json:"differauthor"`
	//ExpiredTime      int    `yml:"expiretime" json:"expired_time"`
}

type SearchEngine struct {
	SauceNaoConf SauceNaoConf `json:"sauceNaoConf"`
}

type PixivConf struct {
	Cookie           string `yml:"cookie" json:"cookie"`
	Agelimit         bool   `yml:"r_18" json:"r_18"`
	Downloadposition string `yml:"downloadposition" json:"downloadposition"`
	LikeLimit        int    `yml:"likelimit" json:"likelimit"`
	Retry429         int    `yml:"retry429" json:"retry429"`
	Downloadinterval int    `yml:"downloadinterval" json:"downloadinterval"`
	Retryinterval    int    `yml:"retryinterval" json:"retryinterval"`
	DifferAuthor     bool   `yml:"differauthor" json:"differauthor"`
	ExpiredTime      int    `yml:"expiretime" json:"expired_time"`
}

type SauceNaoConf struct {
	ApiKey  string `json:"api_key"`
	Numbers int    `json:"numbers"`
}

type GIFEncoderConf struct {
	Quality         int     `yml:"quality" json:"quality"`
	Dither          string  `yml:"dither" json:"dither"`
	ContrastBoost   float64 `yml:"contrast_boost" json:"contrast_boost"`
	SaturationBoost float64 `yml:"saturation_boost" json:"saturation_boost"`
}

func (s *Settings) MsgDetail() string {
	res, _ := json.Marshal(s)
	return string(res)
}
func NowSetting() Settings {
	return *setting
}
func (s *Settings) UpdateSettings(NewSetting Settings) {
	rwl.Lock()
	defer rwl.Unlock()
	tmp := *s
	tmp.UseProxy = NewSetting.UseProxy
	tmp.Proxy = NewSetting.Proxy

	tmpPixivConf := tmp.PixivConf

	tmpPixivConf.LikeLimit = max(NewSetting.PixivConf.LikeLimit, 0)
	err := os.MkdirAll(NewSetting.PixivConf.Downloadposition, os.ModePerm)
	if err != nil {
		tmpPixivConf.Downloadposition = "Download"
		fmt.Println(NewSetting.PixivConf.Downloadposition)
	} else {
		tmpPixivConf.Downloadposition = NewSetting.PixivConf.Downloadposition
	}
	tmpPixivConf.DifferAuthor = NewSetting.PixivConf.DifferAuthor
	tmpPixivConf.Agelimit = NewSetting.PixivConf.Agelimit
	tmpPixivConf.LikeLimit = NewSetting.PixivConf.LikeLimit
	tmpPixivConf.Cookie = NewSetting.PixivConf.Cookie
	tmpPixivConf.Retry429 = max(NewSetting.PixivConf.Retry429, 5000)
	tmpPixivConf.Retryinterval = max(NewSetting.PixivConf.Retryinterval, 800)
	tmpPixivConf.Downloadinterval = max(NewSetting.PixivConf.Downloadinterval, 500)
	tmpPixivConf.ExpiredTime = max(7, min(366, NewSetting.PixivConf.ExpiredTime))
	*s = tmp
}

func UpdateSettings(NewSetting Settings) {
	rwl.Lock()
	defer rwl.Unlock()
	tmp := *setting
	tmp.UseProxy = NewSetting.UseProxy
	tmp.Proxy = NewSetting.Proxy
	tmp.Prefix = NewSetting.Prefix

	tmpPixivConf := tmp.PixivConf
	tmpPixivConf.LikeLimit = max(NewSetting.PixivConf.LikeLimit, 0)
	err := os.MkdirAll(NewSetting.PixivConf.Downloadposition, os.ModePerm)
	if err != nil {
		tmpPixivConf.Downloadposition = "Download"
	} else {
		tmpPixivConf.Downloadposition = NewSetting.PixivConf.Downloadposition
	}
	tmpPixivConf.DifferAuthor = NewSetting.PixivConf.DifferAuthor
	tmpPixivConf.Agelimit = NewSetting.PixivConf.Agelimit
	tmpPixivConf.LikeLimit = NewSetting.PixivConf.LikeLimit
	tmpPixivConf.Cookie = NewSetting.PixivConf.Cookie
	tmpPixivConf.Retry429 = max(NewSetting.PixivConf.Retry429, 5000)
	tmpPixivConf.Retryinterval = max(NewSetting.PixivConf.Retryinterval, 800)
	tmpPixivConf.Downloadinterval = max(NewSetting.PixivConf.Downloadinterval, 500)
	tmpPixivConf.ExpiredTime = max(7, min(366, NewSetting.PixivConf.ExpiredTime))
	tmp.PixivConf = tmpPixivConf

	tmpImageEngine := tmp.SearchEngine

	tmpImageEngine.SauceNaoConf.ApiKey = NewSetting.SearchEngine.SauceNaoConf.ApiKey
	tmpImageEngine.SauceNaoConf.Numbers = NewSetting.SearchEngine.SauceNaoConf.Numbers

	tmp.SearchEngine = tmpImageEngine

	*setting = tmp
}

func SaveSettings() {
	out, _ := yaml.Marshal(setting)
	os.MkdirAll("configs", os.ModePerm)
	ioutil.WriteFile("configs/settings.yml", out, 0644)
}

func init() {
	setting = &Settings{}
}
