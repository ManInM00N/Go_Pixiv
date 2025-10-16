package handler

import (
	"main/configs"
	"main/internal/taskQueue"
	"time"
)

func Download_By_NovelId(id string) {
	Setting := configs.NowSetting()
	taskQueue.SinglePool.AddTask(func() (any, error) {
		DownloadNovel(id)
		time.Sleep(time.Millisecond * time.Duration(Setting.Downloadinterval))
		return nil, nil
	})
}

func Download_By_SeriesId(id string) {

}
