package DAO

import (
	"main/configs"
	"net/http"
	url2 "net/url"
	"time"
)

func GetClient() *http.Client {
	setting := configs.NowSetting()
	proxyURL, _ := url2.Parse(setting.Prefix + setting.Proxy)
	if !setting.UseProxy {
		proxyURL = nil
	}
	return &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyURL(proxyURL),
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: time.Second * 5,
			TLSHandshakeTimeout:   time.Second * 5,
			ExpectContinueTimeout: time.Second * 3,
		},
	}
}
