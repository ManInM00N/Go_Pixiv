package init

import (
	"net/http"
	url2 "net/url"
	"time"
)

func GetClient() *http.Client {
	proxyURL, _ := url2.Parse(Setting.Proxy)
	return &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyURL(proxyURL),
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: time.Second * 5,
		},
	}
}
