package DAO

import (
	"main/configs"
	"main/pkg/utils"
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
			DisableKeepAlives:     false,            // ⭐ 启用 keep-alive 以重用连接
			MaxIdleConns:          100,              // ⭐ 最大空闲连接数
			MaxIdleConnsPerHost:   10,               // ⭐ 每个主机的最大空闲连接数
			IdleConnTimeout:       90 * time.Second, // ⭐ 空闲连接超时
			ResponseHeaderTimeout: 10 * time.Second, // ⭐ 增加超时时间
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 3 * time.Second,
		},
		Timeout: 30 * time.Second, // ⭐ 添加总体超时
	}
}

// 创建 HTTP 请求
func CreatePixivRequest(url string, setting *configs.Settings) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.DebugLog.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	req.Header.Set("referer", "https://www.pixiv.net")
	req.Header.Set("cookie", "PHPSESSID="+setting.PixivConf.Cookie)

	return req, nil
}
