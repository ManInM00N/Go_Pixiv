package imageService

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"main/configs"
	"main/internal/pixivlib/handler"
	"main/pkg/utils"
	"net/http"
	"path/filepath"
)

func SauceNaoSearch(c *gin.Context) {
	if SauceNaoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "SauceNAO 服务未初始化",
		})
		return
	}

	if SauceNaoService.APIKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "SauceNAO API Key 未配置，请先在设置中配置",
		})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		utils.DebugLog.Println("接收文件失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "未找到图片文件",
		})
		return
	}
	defer file.Close()

	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "图片大小不能超过 10MB",
		})
		return
	}

	ext := filepath.Ext(header.Filename)
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "只支持 JPG 和 PNG 格式的图片",
		})
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.DebugLog.Println("读取文件失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "读取文件失败",
		})
		return
	}

	upload, err := ImgManager.Upload(fileBytes, header.Filename)
	if err != nil {
		utils.DebugLog.Println("文件上传云失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "文件上传云失败",
		})
		return
	}
	set := configs.NowSetting().SearchEngine.SauceNaoConf
	utils.WarnLog.Printf("Url: %s ,Svr: %s ,Del: %s , api_key: %s", upload.URL, upload.ServiceName, upload.DeleteToken, set.ApiKey)
	result, err := SauceNaoService.SearchByURL(upload.URL)
	if err != nil {
		utils.DebugLog.Println("SauceNAO 搜索失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("搜索失败: %v", err),
		})
		return
	}

	// 8. 返回搜索结果（直接返回原始结果，前端已经有解析逻辑）
	c.JSON(http.StatusOK, result)
}

// SauceNaoQuota 获取 SauceNAO API 配额信息
func SauceNaoQuota(c *gin.Context) {
	if SauceNaoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "SauceNAO 服务未初始化",
		})
		return
	}

	if SauceNaoService.APIKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "SauceNAO API Key 未配置",
		})
		return
	}

	quota, err := SauceNaoService.GetAPIQuota()
	if err != nil {
		utils.DebugLog.Println("获取配额失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取配额失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, quota)
}

func UploadPicAndSearch(c *gin.Context) {
	p := c.Query("p")
	types := c.Query("type")
	data := handler.GetFollowMsg(types, p, "all")
	if data == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取关注数据失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"data": data,
		"num":  len(data),
	})
}
