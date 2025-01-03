package init

import (
	"context"
	"io/ioutil"
	. "main/backend/src/DAO"
	"os"
	"time"

	"github.com/ManInM00N/go-tool/goruntine"
	"github.com/devchat-ai/gopool"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ymlfile *os.File
	Setting Settings
	is      = false
)

func UpdateSettings() {
	out, _ := yaml.Marshal(&Setting)
	ioutil.WriteFile("settings.yml", out, 0644)
}

func init() {
	if is {
		return
	}
	is = true
	Log_init()
	// logoInit()

	Ctx, Cancel = context.WithCancel(context.Background())
	ymlfile, _ = os.OpenFile("settings.yml", os.O_RDWR, 0644)
	defer ymlfile.Close()

	bytevalue, _ := ioutil.ReadAll(ymlfile)
	yaml.Unmarshal(bytevalue, &Setting)
	Setting.Prefix = "http://127.0.0.1:"
	Setting.LikeLimit = max(Setting.LikeLimit, 0)
	_, err := os.Stat(Setting.Downloadposition)
	if err != nil {
		Setting.Downloadposition = "Download"
	}
	Setting.Retry429 = max(Setting.Retry429, 5000)
	Setting.Retryinterval = max(Setting.Retryinterval, 1000)
	Setting.Downloadinterval = max(Setting.Downloadinterval, 200)
	DebugLog.Println("Check settings:"+Setting.Proxy, "PHPSESSID="+Setting.Cookie, "Download Position=", Setting.Downloadposition)
	UpdateSettings()
	RankPool = goruntine.NewGoPool(200, 1)
	TaskPool = goruntine.NewGoPool(200, 1)
	FollowPool = goruntine.NewGoPool(200, 1)
	TaskPool.Run()
	RankPool.Run()
	FollowPool.Run()
	FollowLoadPool = gopool.NewGoPool(2, gopool.WithTaskQueueSize(400))
	RankloadPool = gopool.NewGoPool(2, gopool.WithTaskQueueSize(5000))
	SinglePool = gopool.NewGoPool(1, gopool.WithTaskQueueSize(100))
	P = gopool.NewGoPool(4, gopool.WithTaskQueueSize(5000))
}

func CacheInit() {
	var err error
	Db, err = gorm.Open(sqlite.Open("cache.db"), &gorm.Config{})
	if err != nil {
		DebugLog.Fatalln(err)
	}
	err = Db.AutoMigrate(&Cache{})
	Clean()
}

func Clean() {
	tx := Db.Begin()
	expireDuration := time.Duration(Setting.ExpiredTime) * 24 * time.Hour
	if err := tx.Where("created_at <= ?", time.Now().Add(-expireDuration)).Delete(&Cache{}).Error; err != nil {
		tx.Rollback()
		DebugLog.Println(err)
	}
	if err := tx.Commit().Error; err != nil {
		DebugLog.Println(err)
	}
}
