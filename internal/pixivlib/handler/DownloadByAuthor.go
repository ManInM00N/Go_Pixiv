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
	authorMsg, _ := GetAuthorInfo(statics.StringToInt64(text))
	progressInfo := &TaskInfo{
		Status: "Waiting",
		ID:     id,
		Name:   fmt.Sprintf("%s(%s) artworks ", text, authorMsg["name"].String()),
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
			taskQueue.WaitingTasks--
			progressInfo.Status = "Running"
			progressInfo.BeginTime = time.Now()
			progressInfo.Total = uint64(len(all))
			if err != nil {
				utils.DebugLog.Println("Error getting author", err)
				callEvent("UpdateTerminal", fmt.Sprintf("Error getting author %s %s\n", text, err.Error()))
				//callEvent("Pop")
				return
			}

			ProcessMax = int64(len(all))
			utils.InfoLog.Println(text + "'s artworks Start download")
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
				taskQueue.P.AddTask(func() (interface{}, error) {
					// time.Sleep(1 * time.Second)
					if taskQueue.IsClosed {
						return nil, nil
					}

					temp := k
					illust, err := work(statics.StringToInt64(temp), options)
					defer func() {
						progressInfo.Current++
					}()
					if err != nil {
						if !ContainMyerror(err) {
							c <- temp
						}
						return nil, nil
					}

					if illust.IllustType == UgoiraType {
						callEvent("downloadugoira", illust.Pid, illust.Width, illust.Height, illust.Frames, illust.Source)
						time.Sleep(10 * time.Second)
					}
					if !Download(illust, options) {
						c <- temp
						return nil, nil
					}
					satisfy++
					return nil, nil
				})
			}
			taskQueue.P.Wait()
			utils.DebugLog.Printf("%s ,failed %d , satisfied %d \n", progressInfo.Name, len(c), satisfy)
			for len(c) > 0 {
				if taskQueue.IsClosed {
					return
				}
				ss := <-c
				// log.Println(ss, " Download failed Now retrying")
				taskQueue.P.AddTask(func() (interface{}, error) {
					if a, b := JustDownload(ss, options, callEvent); b {
						satisfy += a
					}
					return nil, nil
				})
			}
			taskQueue.P.Wait()
			utils.InfoLog.Println(text+"'s artworks -> Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all))
			close(c)
			callEvent("UpdateTerminal", fmt.Sprintf("%s(%s)'s artworks -> Satisfied %d in %d", authorMsg["name"], text, satisfy, len(all)))
			satisfy = 0
			progressInfo.EndTime = time.Now()
			progressInfo.Status = "Done"
		}, progressInfo, 0)
	taskQueue.TaskPool.Add(task)
}
