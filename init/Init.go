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
	"main/internal/pixivlib/DAO"
	"main/internal/taskQueue"
	"main/pkg/utils"
	"net/http"
	_ "net/http/pprof"
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
	Setting := &Settings{}
	yaml.Unmarshal(bytevalue, Setting)
	Setting.Prefix = "http://127.0.0.1:"
	Setting.PixivConf.LikeLimit = max(Setting.PixivConf.LikeLimit, 0)
	_, err := os.Stat(Setting.PixivConf.Downloadposition)
	if err != nil {
		Setting.PixivConf.Downloadposition = "Download"
	}
	Setting.PixivConf.Retry429 = max(Setting.PixivConf.Retry429, 5000)
	Setting.PixivConf.Retryinterval = max(Setting.PixivConf.Retryinterval, 1500)
	Setting.PixivConf.Downloadinterval = max(Setting.PixivConf.Downloadinterval, 700)
	UpdateSettings(*Setting)
	utils.DebugLog.Println("Check settings:"+Setting.Proxy, "PHPSESSID="+Setting.PixivConf.Cookie, "Download Position=", Setting.PixivConf.Downloadposition)
	SaveSettings()

	taskQueue.TaskPool = goruntine.NewTaskPool(1, 1,
		goruntine.WithLowestPriorityFirst(),
		goruntine.WithTaskEqualityByInfoFunc(func(a, b any) bool {
			return a.(DAO.TaskInfo).ID == b.(DAO.TaskInfo).ID
		}),
	)
	taskQueue.TaskPool.Run()

	taskQueue.SinglePool = gopool.NewGoPool(1, gopool.WithTaskQueueSize(100))
	taskQueue.P = goruntine.NewTaskPool(4, 10,
		goruntine.WithFIFO())
	taskQueue.P.Run()

	go func() {
		http.ListenAndServe("localhost:6060", nil) // 使用 pprof 监听端口
	}()
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
	expireDuration := time.Duration(NowSetting().PixivConf.ExpiredTime) * 24 * time.Hour
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
	taskQueue.P.Stop()
	taskQueue.P.Wait()
	defer func() {
		taskQueue.TaskPool.Stop()
		taskQueue.SinglePool.Release()
	}()
}
