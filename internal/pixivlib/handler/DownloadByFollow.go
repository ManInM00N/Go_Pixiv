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
		Status:  "Waiting",
		ID:      id,
		Name:    "Follow Illust Page " + page,
		Current: 0,
		Total:   0,
	}
	task, _ := taskQueue.TaskPool.NewTask(func() {
		if taskQueue.IsClosed {
			progressInfo.Status = "Pool closed"
			return
		}
		op := NewOption(SufWithRankmode(Type), SufWithPage(page))
		utils.InfoLog.Println("follow page", page, " "+Type+" start downloading")
		callEvent("UpdateTerminal", fmt.Sprintln("follow page", page, " "+Type+" start downloading"))
		progressInfo.Status = "Running"
		c := make(chan string, 2000)
		tmp, err := GetFollow(op)
		all := tmp.Get("page").Get("ids").Array()
		thumbs := tmp.Get("thumbnails").Get("illust").Array()
		progressInfo.Total = uint64(len(all))
		progressInfo.BeginTime = time.Now()
		if err != nil {
			utils.DebugLog.Println("Error getting Follow", err)
			callEvent("UpdateTerminal", fmt.Sprintln("follow page", page, Type, err))
			return
		}
		setting := NowSetting()
		utils.InfoLog.Println("follow page", page, " "+Type+" Start download")
		options := NewOption(WithMode(ByPid), WithR18(setting.PixivConf.Agelimit), WithLikeLimit(0), WithDiffAuthor(false), SufWithRankmode(Type))
		var cnt int64
		for _, v := range thumbs {
			k := v.Get("id").String()
			if taskQueue.IsClosed {
				return
			}

			if Db.Model(&Cache{}).Where("download_id = ?", k).Count(&cnt); cnt == 1 {
				progressInfo.Current++
				continue
			}
			weight := 1
			if v.Get("illustType").Int() == UgoiraType {
				weight = 3
			}
			task, _ := taskQueue.P.NewTaskWithCost(func() {
				if taskQueue.IsClosed {
					return
				}
				illust, err := work(statics.StringToInt64(k), options)
				defer func() {
					progressInfo.Current++
				}()
				if err != nil {
					if !ContainMyerror(err) {
						c <- k
					}
					return
				}
				Download(illust, options)
				time.Sleep(time.Second * 3)
				return
			}, k, 1, weight)
			taskQueue.P.Add(task)
			if len(c) > 0 {
				fmt.Println("Illust Failed ", len(c))
			}
		}
		taskQueue.P.Wait()
		utils.InfoLog.Println("follow page", page, " "+Type, "Successfully downloaded illusts: ", len(all))
		callEvent("UpdateTerminal", fmt.Sprintf("follow page %s type %s all %d\n", page, Type, len(all)))
		close(c)
		progressInfo.EndTime = time.Now()
		progressInfo.Status = "Done"
	}, progressInfo, 0)
	taskQueue.TaskPool.Add(task)
}
