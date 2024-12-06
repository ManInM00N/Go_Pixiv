package init

import (
	"io"
	"main/backend/src/DAO"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateSetting(c *gin.Context) {
	var set DAO.Settings
	c.BindJSON(&set)
	DebugLog.Println(set.MsgDetail())
	Setting.UpdateSettings(set)
	UpdateSettings()
}

func GetSetting(c *gin.Context) {
	c.JSON(200, gin.H{
		"setting": Setting,
	})
}

func PreviewUrl(c *gin.Context) {
	imageURL, ok := c.GetQuery("url")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		DebugLog.Println("unknown error")
		return
	}
	if imageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'url' query parameter"})
		DebugLog.Println("Missing url")
		return
	}

	// 发起请求到目标图片地址
	client := &http.Client{}
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		DebugLog.Println("request failed")
		return
	}

	// 设置 Referer 头
	req.Header.Set("Referer", "https://www.pixiv.net")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch image"})
		DebugLog.Println("request failed")
		return
	}
	defer resp.Body.Close()

	// 将目标图片的内容和 Content-Type 返回给前端
	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Status(resp.StatusCode)
	DebugLog.Println(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
