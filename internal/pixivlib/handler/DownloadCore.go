package handler

import (
	"encoding/json"
	"fmt"
	_ "image/png"
	. "main/configs"
	. "main/internal/cache/DAO"
	. "main/internal/pixivlib/DAO"
	"main/internal/taskQueue"
	"main/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/ManInM00N/go-tool/statics"
)

var (
	Logo             []byte
	NowTaskMsg       = ""
	QueueTaskMsg     = ""
	ProcessMax       = int64(0)
	ProcessNow       = int64(0)
	RankLoadingNow   = false
	FollowLoadingNow = false
)

func Download_By_Pid(text string, callEvent func(name string, data ...interface{})) {
	text = statics.CatchNumber(text)
	println(text)
	if text == "" {
		return
	}
	taskQueue.SinglePool.AddTask(func() (interface{}, error) {
		op := NewOption(WithMode(ByPid), WithLikeLimit(0), WithR18(true), WithShowSingle(true), WithDiffAuthor(false))
		JustDownload(text, op, callEvent)
		return nil, nil
	})
}

func Download_By_NovelId(id string) {
	taskQueue.SinglePool.AddTask(func() (any, error) {
		DownloadNovel(id)
		time.Sleep(time.Millisecond * time.Duration(Setting.Downloadinterval))
		return nil, nil
	})
}

func Download_By_Author(text string, callEvent func(name string, data ...interface{})) {
	text = statics.CatchNumber(text)
	if text == "" {
		return
	}
	utils.InfoLog.Println(text + " pushed TaskQueue")
	taskQueue.WaitingTasks++
	callEvent("Push", text)
	taskQueue.TaskPool.Add(func() {
		if taskQueue.IsClosed {
			return
		}
		c := make(chan string, 2000)
		all, err := GetAuthor(statics.StringToInt64(text))
		taskQueue.WaitingTasks--
		if err != nil {
			utils.DebugLog.Println("Error getting author", err)
			callEvent("UpdateTerminal", fmt.Sprintln("Error getting author", text, err))
			callEvent("Pop")
			return
		}
		ProcessMax = int64(len(all))
		ProcessNow = 0
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
				ProcessNow++
				callEvent("UpdateProcess", 100*ProcessNow/max(ProcessMax, 1))
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
					ProcessNow++
					callEvent("UpdateProcess", 100*ProcessNow/max(ProcessMax, 1))
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
					return nil, nil
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
		println(len(c), " ", satisfy)
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
		NowTaskMsg = "No Task in queue"
		callEvent("UpdateTerminal", fmt.Sprintln(text+"'s artworks -> Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all)))
		ProcessNow = 0
		callEvent("UpdateProcess", 100*ProcessNow/max(ProcessMax, 1))
		callEvent("Pop")
		satisfy = 0
	})
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
	println(date)

	for i := int64(1); i < int64(3); i++ {
		temp := i
		callEvent("Push", fmt.Sprint(date+" page", i, " "+Type))
		utils.InfoLog.Println(text, Type, " page ", i, " pushed TaskQueue")
		taskQueue.TaskPool.Add(func() {
			if taskQueue.IsClosed {
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
			taskQueue.WaitingTasks--
			if err != nil {
				utils.DebugLog.Println("Error getting Rank", err)
				callEvent("UpdateTerminal", fmt.Sprintln(date+" page", i, " "+Type, err))
				callEvent("Pop")
				return
			}
			ProcessMax = int64(len(all))
			ProcessNow = 0
			callEvent("UpdateProcess", 100*ProcessNow/max(ProcessMax, 1))

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
					ProcessNow++
					callEvent("UpdateProcess", 100*ProcessNow/max(ProcessMax, 1))
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
						ProcessNow++
						callEvent("UpdateProcess", 100*ProcessNow/max(ProcessMax, 1))
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
			println(len(c), " ", satisfy)
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
			utils.InfoLog.Println(op.RankDate+" "+op.Rank+"'s artworks -> Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all))
			callEvent("UpdateTerminal", fmt.Sprintln(op.RankDate+" "+op.Rank+"'s artworks -> Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all)))
			satisfy = 0
			close(c)
			ProcessNow = 0
			callEvent("UpdateProcess", ProcessNow/max(ProcessMax, 1))
			callEvent("Pop")
		})

	}
}

func Download_By_FollowPage(page, Type string, callEvent func(name string, data ...interface{})) {
	taskQueue.WaitingTasks++

	callEvent("Push", fmt.Sprint("follow page", page, Type))
	taskQueue.TaskPool.Add(func() {
		if taskQueue.IsClosed {
			return
		}
		op := NewOption(SufWithRankmode(Type), SufWithPage(page))
		utils.InfoLog.Println("follow page", page, " "+Type+" pushed queue")

		// println(page)
		c := make(chan string, 2000)
		tmp, err := GetFollow(op)
		all := tmp.Get("page").Get("ids").Array()
		taskQueue.WaitingTasks--
		if err != nil {
			utils.DebugLog.Println("Error getting Follow", err)
			callEvent("UpdateTerminal", fmt.Sprintln("follow page", page, Type, err))
			callEvent("Pop")
			return
		}
		ProcessMax = int64(len(all))
		ProcessNow = 0
		callEvent("UpdateProcess", 100*ProcessNow/max(ProcessMax, 1))

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
				ProcessNow++
				callEvent("UpdateProcess", 100*ProcessNow/max(ProcessMax, 1))
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
					ProcessNow++
					callEvent("UpdateProcess", 100*ProcessNow/max(ProcessMax, 1))
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
				return nil, nil
			})
		}
		taskQueue.P.Wait()
		println(len(c), " ", satisfy)
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
		ProcessNow = 0
		callEvent("UpdateProcess", ProcessNow/max(ProcessMax, 1))
		callEvent("Pop")
	})
}

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

func GetFollowMsg(page, Type string) []FollowData {
	op := NewOption(WithR18(true), SufWithRankmode(Type), SufWithPage(page))
	data, err := GetFollow(op)
	if err != nil {
		utils.DebugLog.Println("Error getting Follow")
		return nil
	}
	list := make([]FollowData, 0, 100)
	all := data.Get("thumbnails").Get("illust").Array()
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
