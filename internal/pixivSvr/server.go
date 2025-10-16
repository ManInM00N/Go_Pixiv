package pixivSvr

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var R *gin.Engine

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func ServerInit() {
	//gin.SetMode(gin.DebugMode)
	gin.SetMode(gin.ReleaseMode)
	R = gin.Default()
	R.Use(Cors())
	Api := R.Group("/api")
	Api.POST("/update", UpdateSetting)
	Api.GET("/getsetting", GetSetting)
	Api.GET("/followpage", Followlist)
	Api.GET("/rankpage", RankList)
	Api.GET("/ws", UpdateProgress)
	Api.GET("/preview", PreviewUrl)
	Api.GET("/novel_rproxy", NovelContent)
	Api.GET("/getugoira", GIFResource)
	Api.POST("/saveugoira", FetchGIF)
	Api.GET("/get_illust_page", GetIllustPage)
	Api.GET("/get_novel", NovelContent)
	Api.GET("/get_author", FetchAuthorInfo)
	Ws := R.Group("/ws")
	Ws.GET("/progress", UpdateProgress)
	Ws.GET("/rank", Transform)
}
