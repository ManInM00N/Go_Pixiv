package init

import (
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/ManInM00N/go-tool/statics"
	_ "image/png"
	"strings"

	. "main/backend/src/DAO"
	"strconv"
)

var (
	Appwindow    fyne.Window
	Logo         []byte
	NowTaskMsg   = ""
	QueueTaskMsg = ""
	ProcessMax   = int64(0)
	ProcessNow   = int64(0)
)

func Download_By_Pid(text string) {
	text = statics.CatchNumber(text)
	if text == "" {
		return
	}
	SinglePool.AddTask(func() (interface{}, error) {
		op := NewOption(WithMode(ByPid), WithLikeLimit(0), WithR18(true), WithShowSingle(true))
		JustDownload(text, op)
		return nil, nil
	})
}
func Download_By_Author(text string, callEvent func(name string, data ...interface{})) {
	text = statics.CatchNumber(text)
	if text == "" {
		return
	}
	InfoLog.Println(text + " pushed TaskQueue")
	WaitingTasks++
	callEvent("Push", text)
	TaskPool.Add(func() {
		if IsClosed {
			return
		}
		c := make(chan string, 2000)
		all, err := GetAuthor(statics.StringToInt64(text))
		WaitingTasks--
		if err != nil {
			DebugLog.Println("Error getting author", err)
			callEvent("UpdateTerminal", fmt.Sprintln("Error getting author", text, err))
			callEvent("Pop")
			return
		}
		ProcessMax = int64(len(all))
		ProcessNow = 0
		InfoLog.Println(text + "'s artworks Start download")
		satisfy := 0
		options := NewOption(WithMode(ByAuthor), WithR18(Setting.Agelimit), WithLikeLimit(Setting.LikeLimit))

		for key, _ := range all {
			k := key
			if IsClosed {
				return
			}
			P.AddTask(func() (interface{}, error) {
				//time.Sleep(1 * time.Second)
				if IsClosed {
					return nil, nil
				}

				temp := k
				illust, err := work(statics.StringToInt64(temp), options)
				if err != nil {
					//continue
					if !ContainMyerror(err) {
						c <- temp
					}
					ProcessNow++
					return nil, nil
				}
				if !Download(illust, options) {
					c <- temp
					ProcessNow++
					return nil, nil
				}
				satisfy++
				ProcessNow++
				callEvent("UpdateProcess", 100*ProcessNow/max(ProcessMax, 1))
				//process.Refresh()

				return nil, nil
			})
		}
		P.Wait()
		println(len(c), " ", satisfy)
		for len(c) > 0 {
			if IsClosed {
				return
			}
			ss := <-c
			//log.Println(ss, " Download failed Now retrying")
			P.AddTask(func() (interface{}, error) {
				if a, b := JustDownload(ss, options); b {
					satisfy += a
				}
				return nil, nil
			})
		}
		P.Wait()
		InfoLog.Println(text+"'s artworks -> Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all))

		satisfy = 0
		close(c)
		NowTaskMsg = "No Task in queue"
		ProcessNow = 0
		callEvent("UpdateTerminal", fmt.Sprintln(text+"'s artworks -> Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all)))
		callEvent("UpdateProcess", 100*ProcessNow/max(ProcessMax, 1))
		callEvent("Pop")
	})

}
func Download_By_Rank(text, Type string, callEvent func(name string, data ...interface{})) {
	WaitingTasks++

	date := ""
	println(text, Type)
	temp := strings.Split(text, "-")
	for i := range temp {
		date += temp[i]
	}
	if len(date) != 6 {
		return
	}
	//hasrank := time.Date(		,tt.Location())
	println(date)

	for i := int64(1); i < int64(3); i++ {
		temp := i
		callEvent("Push", fmt.Sprint(date+" page", i, " "+Type))
		TaskPool.Add(func() {

			if IsClosed {
				return
			}
			page := temp
			dd := date
			op := NewOption(WithType(0), WithRankmode(Type), WithDate(dd), WithR18(true), WithLikeLimit(Setting.LikeLimit), WithPage(strconv.FormatInt(page, 10)))
			InfoLog.Println(op.RankDate+" page", page, " "+op.Rank+"Rank pushed queue")

			//println(page)
			c := make(chan string, 2000)
			all, err := GetRank(op)
			WaitingTasks--
			if err != nil {
				DebugLog.Println("Error getting Rank", err)
				callEvent("UpdateTerminal", fmt.Sprintln(date+" page", i, " "+Type, err))
				callEvent("Pop")
				return
			}
			ProcessMax = int64(len(all))
			ProcessNow = 0
			callEvent("UpdateProcess", 100*ProcessNow/max(ProcessMax, 1))

			InfoLog.Println(op.RankDate + " " + op.Rank + "'s artworks Start download")
			satisfy := 0
			options := NewOption(WithMode(ByAuthor), WithR18(Setting.Agelimit), WithLikeLimit(Setting.LikeLimit), WithDiffAuthor(false), WithDate(op.RankDate), WithRankmode(Type))
			for _, key := range all {
				k := key
				if IsClosed {
					return
				}
				P.AddTask(func() (interface{}, error) {
					//time.Sleep(1 * time.Second)
					if IsClosed {
						return nil, nil
					}
					temp := k
					illust, err := work(statics.StringToInt64(temp.String()), options)
					if err != nil {
						//continue
						if !ContainMyerror(err) {
							c <- temp.Str
						}
						satisfy++
						ProcessNow++
						callEvent("UpdateProcess", 100*ProcessNow/max(ProcessMax, 1))
						return nil, nil
					}
					Download(illust, options)
					satisfy++
					ProcessNow++
					callEvent("UpdateProcess", 100*ProcessNow/max(ProcessMax, 1))
					return nil, nil
				})
			}
			P.Wait()
			println(len(c), " ", satisfy)
			for len(c) > 0 {
				if IsClosed {
					return
				}
				ss := <-c
				//log.Println(ss, " Download failed Now retrying")
				P.AddTask(func() (interface{}, error) {
					if a, b := JustDownload(ss, options); b {
						satisfy += a
					}
					return nil, nil
				})
			}
			P.Wait()
			InfoLog.Println(op.RankDate+" "+op.Rank+"'s artworks -> Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all))
			callEvent("UpdateTerminal", fmt.Sprintln(op.RankDate+" "+op.Rank+"'s artworks -> Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all)))
			satisfy = 0
			close(c)
			ProcessNow = 0
			callEvent("UpdateProcess", ProcessNow/max(ProcessMax, 1))
			callEvent("Pop")
		})

	}
}
