package imageService

import (
	"context"
	"errors"
	"main/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	server *http.Server
	R      *gin.Engine
)

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

	R = gin.New()
	R.Use(Cors())
	Api := R.Group("/api")
	saucenao := Api.Group("/saucenao")
	saucenao.POST("/search", SauceNaoSearch)
	saucenao.GET("/quota", SauceNaoQuota)

	server = &http.Server{
		Addr:    ":7235",
		Handler: R,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			utils.ErrorLog.Printf("%s", err.Error())

		}
		//R.Run(":7234")
	}()
}

func ServerDown() {
	err := server.Shutdown(context.Background())
	if err != nil {
		utils.ErrorLog.Printf("%s", err.Error())
		return
	}
}
