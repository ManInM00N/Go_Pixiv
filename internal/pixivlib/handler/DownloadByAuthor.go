package handler

import (
	"fmt"
	"github.com/ManInM00N/go-tool/statics"
	"github.com/teris-io/shortid"
	. "main/configs"
	. "main/internal/cache/DAO"
	. "main/internal/pixivlib/DAO"
	"main/internal/taskQueue"
	"main/pkg/utils"
	"time"
)

func Download_By_Author(text string, callEvent func(name string, data ...interface{})) {
	text = statics.CatchNumber(text)
	if text == "" {
		return
	}
	utils.InfoLog.Println(text + " pushed TaskQueue")
	taskQueue.WaitingTasks++
	callEvent("Push", text)
	id, _ := shortid.Generate()
	progressInfo := &TaskInfo{
		Status:  "Waiting",
		ID:      id,
		Name:    fmt.Sprintf("%s artworks ", text),
		Current: 0,
		Total:   0,
	}
	Setting := NowSetting()
	task, _ := taskQueue.TaskPool.NewTask(
		func() {
			if taskQueue.IsClosed {
				progressInfo.Status = "Pool closed"
				return
			}
			c := make(chan string, 2000)
			all, err := GetAuthorArtworks(statics.StringToInt64(text))
			authorMsg, _ := GetAuthorInfo(statics.StringToInt64(text))
			progressInfo.Name = fmt.Sprintf("%s(%s) artworks ", text, authorMsg["name"].String())
			utils.InfoLog.Println(authorMsg["name"], text, "'s artwork Start downloading")
			callEvent("UpdateTerminal", fmt.Sprintln("Author ", authorMsg["name"], text, " start downloading"))
			progressInfo.Status = "Running"
			progressInfo.BeginTime = time.Now()
			progressInfo.Total = uint64(len(all))

			if err != nil {
				utils.DebugLog.Println("Error getting author", err)
				callEvent("UpdateTerminal", fmt.Sprintf("Error getting author %s %s\n", text, err.Error()))
				//callEvent("Pop")
				return
			}

			satisfy := 0
			options := NewOption(WithMode(ByAuthor), WithR18(Setting.Agelimit), WithLikeLimit(Setting.LikeLimit))
			var cnt int64
			for key := range all {
				k := key
				if taskQueue.IsClosed {
					return
				}

				if Db.Model(&Cache{}).Where("download_id = ?", k).Count(&cnt); cnt == 1 {
					satisfy++
					progressInfo.Current++
					continue
				}
				task, _ := taskQueue.P.NewTaskWithCost(func() {
					temp := k
					illust, err := work(statics.StringToInt64(temp), options)
					defer func() {
						progressInfo.Current++
					}()
					if err != nil {
						if !ContainMyerror(err) {
							c <- temp
						}
						return
					}
					if !Download(illust, options) {
						c <- temp
						return
					}
					satisfy++
					return
				}, id, 1, 3)
				taskQueue.P.Add(task)
			}
			taskQueue.P.Wait()
			utils.DebugLog.Printf("%s ,failed %d , satisfied %d \n", progressInfo.Name, len(c), satisfy)
			if len(c) > 0 {
				fmt.Println("failed illust", len(c))
			}
			utils.InfoLog.Println(text+"'s artworks -> Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all))
			close(c)
			callEvent("UpdateTerminal", fmt.Sprintf("%s(%s)'s artworks -> Satisfied %d in %d", authorMsg["name"], text, satisfy, len(all)))
			satisfy = 0
			progressInfo.EndTime = time.Now()
			progressInfo.Status = "Done"
		}, progressInfo, 0)
	taskQueue.TaskPool.Add(task)
}
