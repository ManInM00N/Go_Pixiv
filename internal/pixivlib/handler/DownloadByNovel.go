package handler

import (
	"encoding/json"
	"fmt"
	"github.com/teris-io/shortid"
	"main/configs"
	"main/internal/pixivlib/DAO"
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

func Download_By_SeriesId(id string, callEvent func(name string, data ...interface{})) {
	Setting := configs.NowSetting()
	uid, _ := shortid.Generate()
	progressInfo := &DAO.TaskInfo{
		Status: "Waiting",
		ID:     uid,
		Name:   fmt.Sprintf("Series %s ", id),
	}
	task, _ := taskQueue.TaskPool.NewTask(
		func() {
			defer func() {
				fmt.Println(progressInfo.Status)
			}()
			if taskQueue.IsClosed {
				progressInfo.Status = "Pool closed"
				progressInfo.Message = "Pool closed"
				return
			}
			info, err := GetSeriesInfo(id)
			if err != nil {
				progressInfo.Status = "Failed"
				progressInfo.Message = "Falied get Series Info"
				return
			}
			title := info.Get("title").String()
			progressInfo.Name += title
			data, err := GetNovelSeries(id)
			if err != nil {
				progressInfo.Status = "Failed"
				progressInfo.Message = "Falied get Series List"
				return
			}
			var novels []DAO.FollowData
			err = json.Unmarshal([]byte(data.Raw), &novels)
			if err != nil {
				progressInfo.Message = err.Error()
				progressInfo.Status = "Failed"
				fmt.Println(err)
				return
			}
			progressInfo.Status = "Running"
			progressInfo.BeginTime = time.Now()
			progressInfo.Total = uint64(len(novels))
			for _, it := range novels {
				DownloadNovel(it.ID)
				progressInfo.Current++
				time.Sleep(time.Duration(Setting.Downloadinterval) * time.Millisecond)
			}
			callEvent("UpdateTerminal", fmt.Sprintf("%s %s finished download %d page total\n", id, title, len(novels)))

		}, progressInfo, 0)
	taskQueue.TaskPool.Add(task)
}
