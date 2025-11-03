package imageService

import (
	"encoding/json"
	"fmt"
	"io"
	"main/configs"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var SauceNaoService *SauceNAOService

type SauceNAOService struct {
	APIKey     string
	HTTPClient *http.Client
}

// StringOrArray 处理可能是字符串或字符串数组的字段
type StringOrArray []string

func (s *StringOrArray) UnmarshalJSON(data []byte) error {
	// 尝试解析为字符串数组
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*s = arr
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str != "" {
			*s = []string{str}
		} else {
			*s = []string{}
		}
		return nil
	}

	// 如果都失败，返回空数组
	*s = []string{}
	return nil
}

// String 返回字符串表示（用逗号连接）
func (s StringOrArray) String() string {
	if len(s) == 0 {
		return ""
	}
	return strings.Join(s, ", ")
}

// First 返回第一个元素
func (s StringOrArray) First() string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

// NewSauceNAOService 创建新的 SauceNAO 服务实例
func NewSauceNAOService(apiKey string) *SauceNAOService {
	setting := configs.NowSetting()
	proxyURL, _ := url.Parse(setting.Prefix + setting.Proxy)
	if !setting.UseProxy {
		proxyURL = nil
	}
	return &SauceNAOService{
		APIKey: apiKey,
		HTTPClient: &http.Client{
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
			Timeout: 30 * time.Second,
		},
	}
}

func (s *SauceNAOService) GetHTTPClient() *http.Client {
	setting := configs.NowSetting()
	proxyURL, _ := url.Parse(setting.Prefix + setting.Proxy)
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
		Timeout: 30 * time.Second,
	}
}

// SauceNAOResult SauceNAO 搜索结果
type SauceNAOResult struct {
	Header  ResultHeader `json:"header"`
	Results []Result     `json:"results"`
}

// ResultHeader 结果头信息
type ResultHeader struct {
	UserID            string  `json:"user_id"`
	AccountType       string  `json:"account_type"`
	ShortLimit        string  `json:"short_limit"`
	LongLimit         string  `json:"long_limit"`
	LongRemaining     int     `json:"long_remaining"`
	ShortRemaining    int     `json:"short_remaining"`
	Status            int     `json:"status"`
	ResultsRequested  int     `json:"results_requested"`
	SearchDepth       string  `json:"search_depth"`
	MinimumSimilarity float64 `json:"minimum_similarity"`
	QueryImageDisplay string  `json:"query_image_display"`
	QueryImage        string  `json:"query_image"`
	ResultsReturned   int     `json:"results_returned"`
}

// Result 单个搜索结果
type Result struct {
	Header ResultItemHeader `json:"header"`
	Data   ResultData       `json:"data"`
}

// ResultItemHeader 结果项头信息
type ResultItemHeader struct {
	Similarity string `json:"similarity"`
	Thumbnail  string `json:"thumbnail"`
	IndexID    int    `json:"index_id"`
	IndexName  string `json:"index_name"`
}

// ResultData 结果数据
type ResultData struct {
	ExtURLs    []string        `json:"ext_urls"`
	Title      json.RawMessage `json:"title"`       // 可能是字符串或数组
	Author     json.RawMessage `json:"author"`      // 可能是字符串或数组
	AuthorName json.RawMessage `json:"author_name"` // 可能是字符串或数组
	MemberName json.RawMessage `json:"member_name"` // 可能是字符串或数组
	Characters json.RawMessage `json:"characters"`  // 可能是字符串或数组
	Material   json.RawMessage `json:"material"`    // 可能是字符串或数组
	Source     json.RawMessage `json:"source"`      // 可能是字符串或数组
	PixivID    int             `json:"pixiv_id"`
	MemberID   int             `json:"member_id"`
	TwitterID  string          `json:"twitter_id"`
	DanbooruID int             `json:"danbooru_id"`
	GelbooruID int             `json:"gelbooru_id"`
	Creator    json.RawMessage `json:"creator"`  // 可能是字符串或数组
	EngName    json.RawMessage `json:"eng_name"` // 可能是字符串或数组
	JpName     json.RawMessage `json:"jp_name"`  // 可能是字符串或数组
}

func (s *SauceNAOService) SearchByURL(imageURL string) (*SauceNAOResult, error) {
	params := url.Values{}
	params.Add("api_key", s.APIKey)
	params.Add("output_type", "2") // JSON 格式
	params.Add("numres", "10")     // 返回结果数量
	params.Add("url", imageURL)

	apiURL := "https://saucenao.com/search.php?" + params.Encode()

	resp, err := s.GetHTTPClient().Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 返回错误 %d: %s", resp.StatusCode, string(body))
	}

	var result SauceNAOResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}

// SearchResult 简化的搜索结果结构（用于前端）
type SearchResult struct {
	Similarity float64  `json:"similarity"`
	Thumbnail  string   `json:"thumbnail"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	URLs       []string `json:"urls"`
	Source     string   `json:"source"`
	IndexName  string   `json:"indexName"`
	PixivID    int      `json:"pixivId"`
	TwitterID  string   `json:"twitterId"`
}

// SauceNAO 的某些字段可能返回字符串或字符串数组
func parseStringOrArray(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}

	// 尝试解析为字符串
	var str string
	if err := json.Unmarshal(raw, &str); err == nil {
		return str
	}

	// 尝试解析为字符串数组
	var arr []string
	if err := json.Unmarshal(raw, &arr); err == nil {
		if len(arr) > 0 {
			// 返回第一个元素，或者用逗号连接所有元素
			// return arr[0]  // 只返回第一个
			return strings.Join(arr, ", ") // 连接所有元素
		}
		return ""
	}

	// 如果都失败，返回空字符串
	return ""
}

// ParseResults 解析并简化搜索结果
func ParseResults(result *SauceNAOResult) []SearchResult {
	if result == nil || len(result.Results) == 0 {
		return []SearchResult{}
	}

	simplified := make([]SearchResult, 0, len(result.Results))
	for _, r := range result.Results {
		similarity := 0.0
		fmt.Sscanf(r.Header.Similarity, "%f", &similarity)

		// 解析可能是字符串或数组的字段
		author := parseStringOrArray(r.Data.Author)
		if author == "" {
			author = parseStringOrArray(r.Data.AuthorName)
		}
		if author == "" {
			author = parseStringOrArray(r.Data.MemberName)
		}
		if author == "" {
			author = parseStringOrArray(r.Data.Creator)
		}

		title := parseStringOrArray(r.Data.Title)
		if title == "" {
			title = parseStringOrArray(r.Data.EngName)
		}
		if title == "" {
			title = parseStringOrArray(r.Data.JpName)
		}

		source := parseStringOrArray(r.Data.Source)

		simplified = append(simplified, SearchResult{
			Similarity: similarity,
			Thumbnail:  r.Header.Thumbnail,
			Title:      title,
			Author:     author,
			URLs:       r.Data.ExtURLs,
			Source:     source,
			IndexName:  r.Header.IndexName,
			PixivID:    r.Data.PixivID,
			TwitterID:  r.Data.TwitterID,
		})
	}

	return simplified
}

// GetAPIQuota 获取 API 配额信息
func (s *SauceNAOService) GetAPIQuota() (map[string]interface{}, error) {
	params := url.Values{}
	params.Add("api_key", s.APIKey)
	params.Add("output_type", "2")
	params.Add("test_mode", "1")

	apiURL := "https://saucenao.com/search.php?" + params.Encode()

	resp, err := s.GetHTTPClient().Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	var result SauceNAOResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	quota := map[string]interface{}{
		"short_remaining": result.Header.ShortRemaining,
		"long_remaining":  result.Header.LongRemaining,
		"short_limit":     result.Header.ShortLimit,
		"long_limit":      result.Header.LongLimit,
		"account_type":    result.Header.AccountType,
	}

	return quota, nil
}

func (s *SauceNAOService) SetApiKey(key string) error {
	s.APIKey = key
	return nil
}

func init() {
	setting := configs.NowSetting()
	SauceNaoService = NewSauceNAOService(setting.SearchEngine.SauceNaoConf.ApiKey)
}
