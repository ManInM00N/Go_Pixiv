package init

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ManInM00N/go-tool/statics"
	_ "image/png"

	. "main/backend/src/DAO"
	"strconv"
	"time"
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

	TaskPool.Add(func() {
		if IsClosed {
			return
		}
		c := make(chan string, 2000)
		all, err := GetAuthor(statics.StringToInt64(text))
		WaitingTasks--
		if err != nil {
			DebugLog.Println("Error getting author", err)
			if WaitingTasks > 0 {
				QueueTaskMsg = "There are " + fmt.Sprintf("%d", WaitingTasks) + " waiting tasks"
			} else {
				QueueTaskMsg = "There is no tasks waiting"
			}
			callEvent("UpdateQueueNow", QueueTaskMsg)

			return
		}
		if WaitingTasks > 0 {
			QueueTaskMsg = "There are " + fmt.Sprintf("%d", WaitingTasks) + " waiting tasks"
		} else {
			QueueTaskMsg = "There is no tasks waiting"
		}
		NowTaskMsg = text + " are downloading:"
		callEvent("UpdateTaskNow", NowTaskMsg)
		callEvent("UpdateQueueNow", QueueTaskMsg)
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
				//process.Refresh()

				return nil, nil
			})
		}
		P.Wait()
		NowTaskMsg = "Now Recheck " + text
		callEvent("UpdateQueueNow", QueueTaskMsg)
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
		callEvent("UpdateTaskNow", QueueTaskMsg)
	})
	if WaitingTasks > 0 {
		QueueTaskMsg = "There are " + fmt.Sprintf("%d", WaitingTasks) + " waiting tasks"
	} else {
		QueueTaskMsg = "There is no tasks waiting"
	}
	callEvent("UpdateQueueNow", QueueTaskMsg)

}
func WindowInit() {
	app := app.New()
	Appwindow = app.NewWindow("GO Pixiv")
	authorId := widget.NewEntry()
	illustId := widget.NewEntry()
	illustLabel := widget.NewLabel("Download by IllustId")
	authorLabel := widget.NewLabel("Download all Illusts by AuthorId")
	button1 := widget.NewButton("Download", func() {})
	button1.OnTapped = func() {
		text := illustId.Text
		text = statics.CatchNumber(text)
		if text == "" {
			return
		}
		SinglePool.AddTask(func() (interface{}, error) {
			op := NewOption(WithMode(ByPid), WithLikeLimit(0), WithR18(true), WithShowSingle(true))
			JustDownload(text, op)
			return nil, nil
		})
		illustId.SetText("")
	}
	container.New(layout.NewStackLayout())
	process := widget.NewProgressBar()
	process.Min = 0
	waitingtasks := 0
	gridLayout := layout.NewGridLayoutWithColumns(3)
	waitingtasksLabel := widget.NewLabel("There is no tasks waiting")
	waitingtasksLabel.TextStyle.Bold = true
	waitingtasksLabel.TextStyle.TabWidth = 16
	TasknameLabel := widget.NewLabel("No Task in queue")
	TasknameLabel.TextStyle.Bold = true
	TasknameLabel.TextStyle.TabWidth = 16
	Process := container.New(gridLayout,
		TasknameLabel,
		process,
		waitingtasksLabel,
	)
	button2 := widget.NewButton("Download", func() {})
	button2.OnTapped = func() {
		text := authorId.Text
		text = statics.CatchNumber(text)
		if text == "" {
			return
		}
		authorId.SetText("")
		button2.Disabled()
		InfoLog.Println(text + " pushed TaskQueue")
		waitingtasks++

		TaskPool.Add(func() {
			if IsClosed {
				return
			}
			c := make(chan string, 2000)
			all, err := GetAuthor(statics.StringToInt64(text))
			waitingtasks--
			if err != nil {
				DebugLog.Println("Error getting author", err)
				if waitingtasks > 0 {
					waitingtasksLabel.SetText("There are " + fmt.Sprintf("%d", waitingtasks) + " waiting tasks")
				} else {
					waitingtasksLabel.SetText("There is no tasks waiting")
				}
				waitingtasksLabel.Refresh()
				return
			}
			if waitingtasks > 0 {
				waitingtasksLabel.SetText("There are " + fmt.Sprintf("%d", waitingtasks) + " waiting tasks")
			} else {
				waitingtasksLabel.SetText("There is no tasks waiting")
			}
			waitingtasksLabel.Refresh()
			TasknameLabel.SetText(text + " are downloading:")
			TasknameLabel.Refresh()
			process.Max = float64(len(all))
			process.Value = 0
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
						process.Value++
						process.Refresh()

						return nil, nil
					}
					Download(illust, options)
					satisfy++
					process.Value++
					process.Refresh()

					return nil, nil
				})
			}
			P.Wait()
			TasknameLabel.SetText("Now Recheck " + text)
			TasknameLabel.Refresh()
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
			TasknameLabel.SetText("No Task in queue")
			process.SetValue(0)
			process.Refresh()
		})
		if waitingtasks > 0 {
			waitingtasksLabel.SetText("There are " + fmt.Sprintf("%d", waitingtasks) + " waiting tasks")
		} else {
			waitingtasksLabel.SetText("There is no tasks waiting")
		}
		waitingtasksLabel.Refresh()

	}
	date := widget.NewEntry()
	date.PlaceHolder = "rank date"
	choices := widget.NewSelect([]string{"daily", "weekly", "monthly", "male", "female", "rookie", "original", "daily_r18", "weekly_r18", "male_r18", "female_r18"}, func(string) {
	})
	choices.SetSelectedIndex(7)
	button3 := widget.NewButton("Download TodayR18", func() {})
	button3.OnTapped = func() {
		waitingtasks++
		tt := time.Now().Add(time.Hour * -24)
		//hasrank := time.Date(		,tt.Location())
		//println(tt.Date())

		for i := int64(1); i < int64(3); i++ {
			temp := i
			TaskPool.Add(func() {

				if IsClosed {
					return
				}
				page := temp
				dd := date.Text
				modett := int64(choices.SelectedIndex())
				var op *Option

				if len(dd) != 8 {
					op = NewOption(WithType(0), WithRankmode(modett), WithDate(fmt.Sprintf("%04d%02d%02d", tt.Year(), tt.Month(), tt.Day())), WithR18(true), WithLikeLimit(Setting.LikeLimit), WithPage(strconv.FormatInt(page, 10)))

				} else {
					op = NewOption(WithType(0), WithRankmode(modett), WithDate(dd), WithR18(true), WithLikeLimit(Setting.LikeLimit), WithPage(strconv.FormatInt(page, 10)))
				}

				InfoLog.Println(op.RankDate + " " + op.Rank + "Rank pushed queue")
				println(page)

				c := make(chan string, 2000)
				all, err := GetRank(op)
				waitingtasks--
				if err != nil {
					DebugLog.Println("Error getting Rank", err)
					if waitingtasks > 0 {
						waitingtasksLabel.SetText("There are " + fmt.Sprintf("%d", waitingtasks) + " waiting tasks")
					} else {
						waitingtasksLabel.SetText("There is no tasks waiting")
					}
					waitingtasksLabel.Refresh()
					return
				}
				if waitingtasks > 0 {
					waitingtasksLabel.SetText("There are " + fmt.Sprintf("%d", waitingtasks) + " waiting tasks")
				} else {
					waitingtasksLabel.SetText("There is no tasks waiting")
				}
				waitingtasksLabel.Refresh()
				TasknameLabel.SetText(op.RankDate + " " + op.Rank + " are downloading:")
				TasknameLabel.Refresh()
				process.Max = float64(len(all))
				process.Value = 0
				InfoLog.Println(op.RankDate + " " + op.Rank + "'s artworks Start download")
				satisfy := 0
				options := NewOption(WithMode(ByAuthor), WithR18(Setting.Agelimit), WithLikeLimit(Setting.LikeLimit), WithDiffAuthor(false), WithDate(op.RankDate), WithRankmode(modett))
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
							process.Value++
							process.Refresh()

							return nil, nil
						}
						Download(illust, options)
						satisfy++
						process.Value++
						process.Refresh()

						return nil, nil
					})
				}
				P.Wait()
				TasknameLabel.SetText("Now Recheck " + op.RankDate + " " + op.Rank)
				TasknameLabel.Refresh()
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
				satisfy = 0
				close(c)
				TasknameLabel.SetText("No Task in queue")
				process.SetValue(0)
				process.Refresh()
			})

		}

	}
	r18 := widget.NewCheck("R-18", func(i bool) {
	})
	r18.SetChecked(Setting.Agelimit)
	r18.Refresh()
	likelimit := widget.NewLabel("likelimit")
	readlikelimit := widget.NewEntry()
	cookieLabel := widget.NewLabel("cookie")
	readcookie := widget.NewEntry()
	readlikelimit.SetText(statics.Int64ToString(Setting.LikeLimit))
	readcookie.SetText(Setting.Cookie)
	readcookie.Refresh()
	Likelimit := container.New(layout.NewGridWrapLayout(fyne.Size{Width: 100, Height: 38}), likelimit, readlikelimit)
	Cookie := container.New(layout.NewGridWrapLayout(fyne.Size{Width: 100, Height: 38}), cookieLabel, readcookie)
	save := widget.NewButton("Save Settings", func() {
		Setting.Agelimit = r18.Checked
		to := readlikelimit.Text
		if !statics.AllNum(to) {
			to = "0"
		}
		Setting.LikeLimit = statics.StringToInt64(to)
		UpdateSettings()
	})

	setting := container.New(layout.NewGridWrapLayout(fyne.Size{Width: 400, Height: 50}), r18, Likelimit, Cookie, save)

	content := container.New(layout.NewGridLayoutWithColumns(3), illustLabel, illustId, button1, authorLabel, authorId, button2, date, choices, button3)
	//stackqueue := container.NewScroll()
	all := container.NewVBox(content, Process, setting)
	icon, _ := fyne.LoadResourceFromPath("assets/icon.ico")
	app.SetIcon(icon)
	Appwindow.SetIcon(icon)
	Appwindow.SetContent(all)
	Appwindow.Resize(fyne.Size{Width: 500, Height: 350})
	//Appwindow.SetCloseIntercept(func() {
	//
	//	//app.Quit()
	//})
}
