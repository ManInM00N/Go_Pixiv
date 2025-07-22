package init

import (
	"context"
	"github.com/ManInM00N/go-tool/goruntine"
	"github.com/devchat-ai/gopool"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	. "main/configs"
	. "main/internal/cache/DAO"
	"main/internal/taskQueue"
	"main/pkg/utils"
	"os"
	"time"
)

var (
	ymlfile *os.File
	is      = false
)

func init() {
	if is {
		return
	}
	is = true
	utils.Log_init()
	// logoInit()
	taskQueue.Ctx, taskQueue.Cancel = context.WithCancel(context.Background())
	ymlfile, _ = os.OpenFile("configs/settings.yml", os.O_RDWR, 0644)
	defer ymlfile.Close()

	bytevalue, _ := ioutil.ReadAll(ymlfile)
	Setting = &Settings{}
	yaml.Unmarshal(bytevalue, Setting)
	Setting.Prefix = "http://127.0.0.1:"
	Setting.LikeLimit = max(Setting.LikeLimit, 0)
	_, err := os.Stat(Setting.Downloadposition)
	if err != nil {
		Setting.Downloadposition = "Download"
	}
	Setting.Retry429 = max(Setting.Retry429, 5000)
	Setting.Retryinterval = max(Setting.Retryinterval, 1000)
	Setting.Downloadinterval = max(Setting.Downloadinterval, 200)
	utils.DebugLog.Println("Check settings:"+Setting.Proxy, "PHPSESSID="+Setting.Cookie, "Download Position=", Setting.Downloadposition)
	UpdateSettings()
	taskQueue.RankPool = goruntine.NewGoPool(200, 1)
	taskQueue.TaskPool = goruntine.NewGoPool(200, 1)
	taskQueue.FollowPool = goruntine.NewGoPool(200, 1)
	taskQueue.TaskPool.Run()
	taskQueue.RankPool.Run()
	taskQueue.FollowPool.Run()
	taskQueue.FollowLoadPool = gopool.NewGoPool(2, gopool.WithTaskQueueSize(400))
	taskQueue.RankloadPool = gopool.NewGoPool(2, gopool.WithTaskQueueSize(5000))
	taskQueue.SinglePool = gopool.NewGoPool(1, gopool.WithTaskQueueSize(100))
	taskQueue.P = gopool.NewGoPool(4, gopool.WithTaskQueueSize(5000))
}

func CacheInit() {
	var err error
	Db, err = gorm.Open(sqlite.Open("cache.db"), &gorm.Config{})
	if err != nil {
		utils.DebugLog.Fatalln(err)
	}
	err = Db.AutoMigrate(&Cache{})
	Clean()
}

func Clean() {
	tx := Db.Begin()
	expireDuration := time.Duration(Setting.ExpiredTime) * 24 * time.Hour
	if err := tx.Where("created_at <= ?", time.Now().Add(-expireDuration)).Delete(&Cache{}).Error; err != nil {
		tx.Rollback()
		utils.DebugLog.Println(err)
	}
	if err := tx.Commit().Error; err != nil {
		utils.DebugLog.Println(err)
	}
}

func Close() {
	taskQueue.IsClosed = true
	taskQueue.P.Wait()
	defer func() {
		taskQueue.P.Release()
		taskQueue.TaskPool.Close()
		taskQueue.SinglePool.Release()
	}()
}
