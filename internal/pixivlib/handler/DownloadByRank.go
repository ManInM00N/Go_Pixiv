package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ManInM00N/go-tool/statics"
	"github.com/teris-io/shortid"
	. "main/configs"
	. "main/internal/cache/DAO"
	. "main/internal/pixivlib/DAO"
	"main/internal/taskQueue"
	"main/pkg/utils"
	"strconv"
	"strings"
	"time"
)

func GetRankMsg(date, Type, page, content string) []RankData {
	num := int64(0)
	if content == "ugoira" {
		num = 2
	} else if content == "manga" {
		num = 1
	}
	op := NewOption(WithR18(true), SufWithRankmode(Type), SufWithPage(page), SufWithType(num), SufWithDate(date))
	data, err := GetRank(op)
	if err != nil {
		utils.DebugLog.Printf("Error getting Rank %v\n", err)
		return nil
	}
	list := make([]RankData, 0, 100)
	all := data.Array()
	for _, value := range all {
		var tmp RankData
		err := json.Unmarshal([]byte(value.Raw), &tmp)
		if err != nil {
			utils.DebugLog.Printf("get follow page %s failed :%v\n", page, err)
			return nil
		}
		// tmp.R18 = utils.HasR18(&tmp.Tags)
		list = append(list, tmp)

	}
	// DebugLog.Println(list)
	return list
}

func Download_By_Rank(text, Type string, callEvent func(name string, data ...interface{})) {
	taskQueue.WaitingTasks++

	date := ""
	println(text, Type)
	temp := strings.Split(text, "-")
	for i := range temp {
		date += temp[i]
	}
	if len(date) != 8 {
		return
	}
	Setting := NowSetting()
	for i := int64(1); i < int64(3); i++ {
		temp := i
		callEvent("Push", fmt.Sprint(date+" page", i, " "+Type))
		utils.InfoLog.Println(text, Type, " page ", i, " pushed TaskQueue")
		id, _ := shortid.Generate()
		progressInfo := &TaskInfo{
			Status: "Waiting",
			ID:     id,
			Name:   "Rank " + Type + " " + text + " Page " + strconv.FormatInt(i, 10),
		}
		task, _ := taskQueue.TaskPool.NewTask(func() {
			if taskQueue.IsClosed {
				progressInfo.Status = "Pool closed"
				return
			}
			page := temp
			dd := date
			op := NewOption(SufWithType(0), SufWithRankmode(Type), SufWithDate(dd), WithMode(ByRank), WithR18(true), WithLikeLimit(Setting.LikeLimit), SufWithPage(strconv.FormatInt(page, 10)))
			utils.InfoLog.Println(op.RankDate+" page", page, " "+op.Rank+"Rank pushed queue")
			c := make(chan string, 2000)
			tmp, err := GetRank(op)
			// DebugLog.Println(tmp.Get("#.illust_id"))
			all := tmp.Get("#.illust_id").Array()

			progressInfo.BeginTime = time.Now()
			progressInfo.Total = uint64(len(all))
			progressInfo.Status = "Running"
			taskQueue.WaitingTasks--
			if err != nil {
				utils.DebugLog.Println("Error getting Rank", err)
				callEvent("UpdateTerminal", fmt.Sprintln(date+" page", i, " "+Type, err))
				progressInfo.Message = "Task failed"
				progressInfo.Status = "Done"
				return
			}
			ProcessMax = int64(len(all))

			utils.InfoLog.Println(op.RankDate + " " + op.Rank + "'s artworks Start download")
			satisfy := 0
			options := NewOption(WithMode(ByRank), WithR18(Setting.Agelimit), WithLikeLimit(Setting.LikeLimit), WithDiffAuthor(false), SufWithDate(op.RankDate), SufWithRankmode(Type))

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
				taskQueue.P.AddTask(func() (interface{}, error) {
					if a, b := JustDownload(ss, options, callEvent); b {
						satisfy += a
					}
					return nil, nil
				})
			}
			taskQueue.P.Wait()
			utils.InfoLog.Println(op.RankDate+" "+op.Rank+"'s artworks -> Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all))
			callEvent("UpdateTerminal", fmt.Sprintf("%s %s's artworks -> Satisfied %d in %d \n", op.RankDate, op.Rank, satisfy, len(all)))
			satisfy = 0
			close(c)
			progressInfo.EndTime = time.Now()
			progressInfo.Status = "Done"
		}, progressInfo, 0)
		taskQueue.TaskPool.Add(task)
	}
}
