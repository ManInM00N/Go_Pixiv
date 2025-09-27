package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ManInM00N/go-tool/statics"
	"github.com/teris-io/shortid"
	"main/internal/taskQueue"

	. "main/configs"
	. "main/internal/cache/DAO"
	. "main/internal/pixivlib/DAO"
	"main/pkg/utils"
	"time"
)

func GetFollowMsg(followType, page, Type string) []FollowData {
	op := NewOption(WithR18(true), SufWithRankmode(Type), SufWithPage(page), WithFollowType(followType))
	data, err := GetFollow(op)
	if err != nil {
		utils.DebugLog.Println("Error getting Follow")
		return nil
	}
	list := make([]FollowData, 0, 100)
	all := data.Get("thumbnails").Get(followType).Array()
	for _, value := range all {
		var tmp FollowData
		err := json.Unmarshal([]byte(value.Raw), &tmp)
		if err != nil {
			utils.DebugLog.Printf("get follow page %s failed", page)
			return nil
		}
		tmp.R18 = utils.HasR18(&tmp.Tags)
		list = append(list, tmp)
	}
	return list
}

func Download_By_FollowPage(page, Type string, callEvent func(name string, data ...interface{})) {
	taskQueue.WaitingTasks++

	callEvent("Push", fmt.Sprint("follow page", page, Type))
	id, _ := shortid.Generate()
	progressInfo := &TaskInfo{
		Status: "Waiting",
		ID:     id,
		Name:   "Follow Illust Page " + page,
	}
	task, _ := taskQueue.TaskPool.NewTask(func() {
		if taskQueue.IsClosed {
			progressInfo.Status = "Pool closed"
			return
		}
		op := NewOption(SufWithRankmode(Type), SufWithPage(page))
		utils.InfoLog.Println("follow page", page, " "+Type+" pushed queue")
		progressInfo.Status = "Running"
		// println(page)
		c := make(chan string, 2000)
		tmp, err := GetFollow(op)
		all := tmp.Get("page").Get("ids").Array()
		progressInfo.Total = uint64(len(all))
		progressInfo.BeginTime = time.Now()
		taskQueue.WaitingTasks--
		if err != nil {
			utils.DebugLog.Println("Error getting Follow", err)
			callEvent("UpdateTerminal", fmt.Sprintln("follow page", page, Type, err))
			return
		}
		ProcessMax = int64(len(all))

		utils.InfoLog.Println("follow page", page, " "+Type+" Start download")
		satisfy := 0
		options := NewOption(WithMode(ByPid), WithR18(Setting.Agelimit), WithLikeLimit(0), WithDiffAuthor(false), SufWithRankmode(Type))
		var cnt int64
		for _, key := range all {
			k := key
			if taskQueue.IsClosed {
				return
			}

			if Db.Model(&Cache{}).Where("download_id = ?", k.String()).Count(&cnt); cnt == 1 {
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
				illust, err := work(statics.StringToInt64(temp.String()), options)
				defer func() {
					satisfy++
					progressInfo.Current++
				}()
				if err != nil {
					if !ContainMyerror(err) {
						c <- temp.Str
					}
					return nil, nil
				}
				if illust.IllustType == UgoiraType {
					callEvent("downloadugoira", illust.Pid, illust.Width, illust.Height, illust.Frames, illust.Source)
					time.Sleep(10 * time.Second)
					return nil, nil
				}
				Download(illust, options)
				time.Sleep(time.Second * 10)
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
		utils.InfoLog.Println("follow page", page, " "+Type, "Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all))
		callEvent("UpdateTerminal", fmt.Sprintln("follow page", page, " "+Type, "Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all)))
		satisfy = 0
		close(c)
		progressInfo.EndTime = time.Now()
		progressInfo.Status = "Done"
	}, progressInfo, 0)
	taskQueue.TaskPool.Add(task)
}
