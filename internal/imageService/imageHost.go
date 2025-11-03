package imageService

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

var ImgManager *ImageHostManager

// ImageHostService 图床服务接口
type ImageHostService interface {
	Upload(imageData []byte, filename string) (string, error)
	Delete(deleteToken string) error
	GetName() string
	SupportsDelete() bool
	//RegisterService(app *application.App)
}

// UploadResult 上传结果，包含删除令牌
type UploadResult struct {
	URL         string `json:"url"`
	DeleteToken string `json:"deleteToken"` // 用于删除的令牌
	DeleteURL   string `json:"deleteUrl"`   // 删除链接（如果有）
}

// ==========  ImgBB 图床服务 ==========
// 免费：每月 5000 次上传，图片永久保存
// 获取 API Key: https://api.imgbb.com/

type ImgBBService struct {
	APIKey     string
	HTTPClient *http.Client
}

func NewImgBBService(apiKey string) *ImgBBService {
	return &ImgBBService{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *ImgBBService) GetName() string {
	return "ImgBB"
}

func (s *ImgBBService) SupportsDelete() bool {
	return true
}

func (s *ImgBBService) Upload(imageData []byte, filename string) (string, error) {
	// 转换为 Base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// 准备表单数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("key", s.APIKey)
	writer.WriteField("image", base64Image)
	writer.WriteField("name", filename)

	if err := writer.Close(); err != nil {
		return "", err
	}

	// 发送请求
	req, err := http.NewRequest("POST", "https://api.imgbb.com/1/upload", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("上传失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("上传失败 %d: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var result struct {
		Data struct {
			URL       string `json:"url"`
			DeleteURL string `json:"delete_url"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if !result.Success {
		return "", fmt.Errorf("上传失败")
	}

	// 返回格式: URL|DeleteURL
	return result.Data.URL + "|" + result.Data.DeleteURL, nil
}

func (s *ImgBBService) Delete(deleteURL string) error {
	// ImgBB 通过访问 delete_url 来删除
	resp, err := s.HTTPClient.Get(deleteURL)
	if err != nil {
		return fmt.Errorf("删除失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("删除失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

// ==========  Imgur 图床服务 ==========
// 免费：每小时 50 次上传，图片永久保存
// 获取 Client ID: https://api.imgur.com/oauth2/addclient

type ImgurService struct {
	ClientID   string
	HTTPClient *http.Client
}

func NewImgurService(clientID string) *ImgurService {
	return &ImgurService{
		ClientID: clientID,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *ImgurService) GetName() string {
	return "Imgur"
}

func (s *ImgurService) SupportsDelete() bool {
	return true
}

func (s *ImgurService) Upload(imageData []byte, filename string) (string, error) {
	// 转换为 Base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// 准备表单数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("image", base64Image)
	writer.WriteField("type", "base64")

	if err := writer.Close(); err != nil {
		return "", err
	}

	// 发送请求
	req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Client-ID "+s.ClientID)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("上传失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("上传失败 %d: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var result struct {
		Data struct {
			Link       string `json:"link"`
			DeleteHash string `json:"deletehash"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if !result.Success {
		return "", fmt.Errorf("上传失败")
	}

	// 返回格式: URL|DeleteHash
	return result.Data.Link + "|" + result.Data.DeleteHash, nil
}

func (s *ImgurService) Delete(deleteHash string) error {
	req, err := http.NewRequest("DELETE", "https://api.imgur.com/3/image/"+deleteHash, nil)
	if err != nil {
		return fmt.Errorf("删除失败: %v", err)
	}
	req.Header.Set("Authorization", "Client-ID "+s.ClientID)

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("删除失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("删除失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

// ==========  SM.MS 图床服务 ==========
// 免费：每分钟 20 次，5MB 限制，图片保存 30 天
// 获取 Token: https://sm.ms/home/apitoken

type SMMSService struct {
	Token      string
	HTTPClient *http.Client
}

func NewSMMSService(token string) *SMMSService {
	return &SMMSService{
		Token: token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *SMMSService) GetName() string {
	return "SM.MS"
}

func (s *SMMSService) SupportsDelete() bool {
	return true
}

func (s *SMMSService) Upload(imageData []byte, filename string) (string, error) {
	// 准备表单数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件
	part, err := writer.CreateFormFile("smfile", filename)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(imageData); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	// 发送请求
	req, err := http.NewRequest("POST", "https://sm.ms/api/v2/upload", body)
	if err != nil {
		return "", err
	}

	if s.Token != "" {
		req.Header.Set("Authorization", s.Token)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("上传失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var result struct {
		Success bool   `json:"success"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    struct {
			URL    string `json:"url"`
			Hash   string `json:"hash"`
			Delete string `json:"delete"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if !result.Success {
		// 图片已存在的情况
		if result.Code == "image_repeated" {
			return result.Data.URL + "|" + result.Data.Hash, nil
		}
		return "", fmt.Errorf("上传失败: %s", result.Message)
	}

	// 返回格式: URL|Hash
	return result.Data.URL + "|" + result.Data.Hash, nil
}

func (s *SMMSService) Delete(hash string) error {
	req, err := http.NewRequest("GET", "https://sm.ms/api/v2/delete/"+hash, nil)
	if err != nil {
		return fmt.Errorf("删除失败: %v", err)
	}

	if s.Token != "" {
		req.Header.Set("Authorization", s.Token)
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("删除失败: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if !result.Success {
		return fmt.Errorf("删除失败: %s", result.Message)
	}

	return nil
}

// ==========  Catbox 图床服务 ==========
// 免费：无限制，200MB 限制，图片永久保存
// 无需 API Key

type CatboxService struct {
	HTTPClient *http.Client
}

func NewCatboxService() *CatboxService {
	return &CatboxService{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *CatboxService) GetName() string {
	return "Catbox"
}

func (s *CatboxService) SupportsDelete() bool {
	return false // Catbox 需要 userhash 才能删除，匿名上传不支持
}

func (s *CatboxService) Upload(imageData []byte, filename string) (string, error) {
	// 准备表单数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("reqtype", "fileupload")

	// 添加文件
	part, err := writer.CreateFormFile("fileToUpload", filename)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(imageData); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	// 发送请求
	req, err := http.NewRequest("POST", "https://catbox.moe/user/api.php", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("上传失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("上传失败 %d: %s", resp.StatusCode, string(respBody))
	}

	// 响应直接是 URL
	urlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	url := strings.TrimSpace(string(urlBytes))
	if !strings.HasPrefix(url, "http") {
		return "", fmt.Errorf("上传失败: %s", url)
	}

	// Catbox 匿名上传不支持删除，返回 URL|（空）
	return url + "|", nil
}

func (s *CatboxService) Delete(deleteToken string) error {
	return fmt.Errorf("Catbox 匿名上传不支持删除")
}

// ==========  临时文件服务 tmpfiles.org ==========
// 免费：无需注册，文件保存 1 小时到永久（根据需要）
// 无需 API Key

type TmpFilesService struct {
	HTTPClient *http.Client
}

func NewTmpFilesService() *TmpFilesService {
	return &TmpFilesService{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *TmpFilesService) GetName() string {
	return "TmpFiles"
}

func (s *TmpFilesService) SupportsDelete() bool {
	return false // TmpFiles 文件会自动过期，不支持手动删除
}

func (s *TmpFilesService) Upload(imageData []byte, filename string) (string, error) {
	// 准备表单数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(imageData); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	// 发送请求
	req, err := http.NewRequest("POST", "https://tmpfiles.org/api/v1/upload", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("上传失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var result struct {
		Status string `json:"status"`
		Data   struct {
			URL string `json:"url"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Status != "success" {
		return "", fmt.Errorf("上传失败")
	}

	// tmpfiles.org 返回的是下载页面 URL，需要转换为直链
	// 例如: https://tmpfiles.org/12345 -> https://tmpfiles.org/dl/12345
	url := strings.Replace(result.Data.URL, "tmpfiles.org/", "tmpfiles.org/dl/", 1)

	// TmpFiles 不支持删除，返回 URL|（空）
	return url + "|", nil
}

func (s *TmpFilesService) Delete(deleteToken string) error {
	return fmt.Errorf("TmpFiles 不支持手动删除，文件会自动过期")
}

// ==========  图床管理器 ==========

type ImageHostManager struct {
	services map[string]ImageHostService
	priority []string // 优先级列表
}

func NewImageHostManager() *ImageHostManager {
	return &ImageHostManager{
		services: make(map[string]ImageHostService),
		priority: []string{},
	}
}

// AddService 添加图床服务
func (m *ImageHostManager) AddService(service ImageHostService) {
	name := service.GetName()
	m.services[name] = service
	m.priority = append(m.priority, name)
}

// UploadResponse 上传响应
type UploadResponse struct {
	URL         string `json:"url"`
	DeleteToken string `json:"deleteToken"`
	ServiceName string `json:"serviceName"`
	CanDelete   bool   `json:"canDelete"`
}

// Upload 上传图片（使用第一个可用的服务）
func (m *ImageHostManager) Upload(imageData []byte, filename string) (*UploadResponse, error) {
	for _, name := range m.priority {
		service := m.services[name]
		result, err := service.Upload(imageData, filename)
		if err == nil {
			// 解析结果: URL|DeleteToken
			parts := strings.Split(result, "|")
			url := parts[0]
			deleteToken := ""
			if len(parts) > 1 {
				deleteToken = parts[1]
			}

			return &UploadResponse{
				URL:         url,
				DeleteToken: deleteToken,
				ServiceName: name,
				CanDelete:   service.SupportsDelete() && deleteToken != "",
			}, nil
		}
		fmt.Printf("%s 上传失败: %v，尝试下一个...\n", name, err)
	}
	return nil, fmt.Errorf("所有图床服务均上传失败")
}

// UploadWithRetry 上传图片（带重试）
func (m *ImageHostManager) UploadWithRetry(imageData []byte, filename string, maxRetries int) (*UploadResponse, error) {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		response, err := m.Upload(imageData, filename)
		if err == nil {
			return response, nil
		}

		if attempt < maxRetries {
			time.Sleep(2 * time.Second)
		}
	}
	return nil, fmt.Errorf("重试 %d 次后仍然失败", maxRetries)
}

// UploadToSpecific 上传到指定服务
func (m *ImageHostManager) UploadToSpecific(serviceName string, imageData []byte, filename string) (*UploadResponse, error) {
	service, exists := m.services[serviceName]
	if !exists {
		return nil, fmt.Errorf("服务 %s 不存在", serviceName)
	}

	result, err := service.Upload(imageData, filename)
	if err != nil {
		return nil, err
	}

	// 解析结果
	parts := strings.Split(result, "|")
	url := parts[0]
	deleteToken := ""
	if len(parts) > 1 {
		deleteToken = parts[1]
	}

	return &UploadResponse{
		URL:         url,
		DeleteToken: deleteToken,
		ServiceName: serviceName,
		CanDelete:   service.SupportsDelete() && deleteToken != "",
	}, nil
}

// Delete 删除图片
func (m *ImageHostManager) Delete(serviceName, deleteToken string) error {
	service, exists := m.services[serviceName]
	if !exists {
		return fmt.Errorf("服务 %s 不存在", serviceName)
	}

	if !service.SupportsDelete() {
		return fmt.Errorf("服务 %s 不支持删除", serviceName)
	}

	if deleteToken == "" {
		return fmt.Errorf("缺少删除令牌")
	}

	return service.Delete(deleteToken)
}

func init() {
	ImgManager = NewImageHostManager()
	ImgManager.AddService(NewTmpFilesService())

}
