package pixivSvr

import (
	"bufio"
	"fmt"
	"io"
	"main/configs"
	. "main/internal/pixivlib/DAO"
	"main/internal/pixivlib/handler"
	"main/pkg/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func UpdateSetting(c *gin.Context) {
	var set configs.Settings
	utils.InfoLog.Println(c.Request.Body)
	err := c.BindJSON(&set)
	if err != nil {
		fmt.Println(err)
		utils.DebugLog.Println(err)
		return
	}
	utils.DebugLog.Println(set.MsgDetail(), set)
	configs.Setting.UpdateSettings(set)
	configs.UpdateSettings()
	c.JSON(http.StatusOK, gin.H{
		"setting": configs.Setting,
	})
}

func GetSetting(c *gin.Context) {
	c.JSON(200, gin.H{
		"setting": configs.Setting,
	})
}

func PreviewUrl(c *gin.Context) {
	imageURL, ok := c.GetQuery("url")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		utils.DebugLog.Println("unknown error")
		return
	}
	if imageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'url' query parameter"})
		utils.DebugLog.Println("Missing url")
		return
	}
	// 发起请求到目标图片地址
	client := GetClient()
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		utils.DebugLog.Println("request failed")
		return
	}

	// 设置 Referer 头
	req.Header.Set("Referer", "https://www.pixiv.net")
	var resp *http.Response
	ok = false
	for i := 0; i < 5; i++ {
		resp, err = client.Do(req)
		if err == nil {
			ok = true
			break
		}
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch image"})
		utils.DebugLog.Println("request failed")
		return
	}

	// 将目标图片的内容和 Content-Type 返回给前端
	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Status(resp.StatusCode)
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		utils.DebugLog.Println(err.Error())
		return
	}
}

func RankList(c *gin.Context) {
	p := c.Query("p")
	mode := c.Query("mode")
	content := c.Query("content")
	rawdate := c.Query("date")
	val := strings.Split(rawdate, "-")
	date := ""
	for _, v := range val {
		date = date + v
	}
	utils.DebugLog.Println(val, date)
	data := handler.GetRankMsg(date, mode, p, content)
	if data == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取排行榜数据失败",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"data": data,
		"num":  len(data),
	})
}

func Followlist(c *gin.Context) {
	p := c.Query("p")
	data := handler.GetFollowMsg(p, "all")
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

func NovelContent(c *gin.Context) {
}
func GIFResource(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'url' query parameter"})
		utils.DebugLog.Println("Missing url")
		return
	}
	// 发起请求到目标图片地址
	client := GetClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		utils.DebugLog.Println("request failed")
		return
	}

	// 设置 Referer 头
	req.Header.Set("Referer", "https://www.pixiv.net")
	var resp *http.Response
	ok := false
	for i := 0; i < 5; i++ {
		resp, err = client.Do(req)
		if err == nil {
			ok = true
			break
		}
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch zip"})
		utils.DebugLog.Println("request failed")
		return
	}
	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Status(resp.StatusCode)
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		utils.DebugLog.Println(err.Error())
		return
	}
}
func FetchGIF(c *gin.Context) {
	Pid := c.Query("id")
	Path := configs.Setting.Downloadposition
	Path = filepath.Join(Path, "GIF")
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		utils.DebugLog.Println("Error retrieving the file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error retrieving file"})
		return
	}
	defer file.Close()

	// 创建保存文件的目录
	err = os.MkdirAll(Path, os.ModePerm)
	if err != nil {
		utils.DebugLog.Println("Error creating uploads directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating directory"})
		return
	}
	Path = filepath.Join(Path, Pid+".gif")
	f, err := os.OpenFile(Path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	bufWriter := bufio.NewWriter(f)
	defer func() {
		bufWriter.Flush()
		f.Close()
	}()
	if _, err = io.Copy(bufWriter, file); err != nil {
		utils.DebugLog.Println(Pid, "Error saving file: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving file"})

		return
	}

	// 返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}
