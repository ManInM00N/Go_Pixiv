package pixivlib

import (
	"bufio"
	"bytes"
	"context"
	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails/v3/pkg/application"
	"io"
	"main/configs"
	"main/internal/pixivlib/DAO"
	. "main/internal/pixivlib/handler"
	. "main/pkg/utils"
	"net/http"
	"time"
)

type Ctl struct {
	ctx context.Context
	App *application.App
}

func NewCtl() *Ctl {

	return &Ctl{}
}
func (a *Ctl) RegisterService(app *application.App) {
	a.App = app
	app.RegisterService(application.NewService(a))
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *Ctl) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *Ctl) DownloadByPid(text string) bool {
	InfoLog.Println("Download illust ", text)
	Download_By_Pid(text, func(name string, data ...interface{}) {
		a.App.Event.Emit(name, data)
	})
	return true
}

func (a *Ctl) DownloadByNovelId(text string) bool {
	InfoLog.Println("Download Novel ", text)
	Download_By_NovelId(text)
	return true
}

func (a *Ctl) ReturnString() string {
	return NowTaskMsg
}

func (a *Ctl) DownloadByAuthorId(text string) bool {
	Download_By_Author(text, func(name string, data ...interface{}) {
		a.App.Event.Emit(name, data)
	})
	return true
}

func (a *Ctl) DownloadByRank(text, Type string) bool {
	Download_By_Rank(text, Type, func(name string, data ...interface{}) {
		a.App.Event.Emit(name, data)
	})
	return true
}

func (a *Ctl) DownloadByFollowPage(page, Type string) bool {
	Download_By_FollowPage(page, Type, func(name string, data ...interface{}) {
		a.App.Event.Emit(name, data)
	})
	return true
}

func (a *Ctl) CheckLogin() bool {
	url, ref := GetUrlRefer("", "", 0)
	Request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		DebugLog.Println("Error creating request", err)
		a.App.Event.Emit("login", "False")
		return false
	}
	client := DAO.GetClient()
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")
	Request.Header.Set("referer", ref)
	Request.Header.Set("Cookie", "PHPSESSID="+configs.NowSetting().Cookie)
	var res *http.Response
	done := make(chan bool)
	go func() {
		for i := 0; i < 3; i++ {
			res, err = client.Do(Request)
			if err == nil {
				done <- true
				return
			}
		}
		done <- false
	}()
	select {
	case <-done:
		break
	case <-time.After(15 * time.Second):

		a.App.Event.Emit("login", "False", "Can't Connect to Pixiv.com timed out. Check your Proxy")
		return false
	}
	if err != nil {
		a.App.Event.Emit("login", "False", "Can't Connect to Pixiv.com . Check your Proxy")
		return false
	}
	var buffer bytes.Buffer
	reader := bufio.NewReader(res.Body)
	io.Copy(&buffer, reader)
	data := buffer.Bytes()
	Results := gjson.ParseBytes(data)
	canbedownload := Results.Get("error").Bool()
	DebugLog.Println("Loading succeed? :", !canbedownload)
	if canbedownload {
		a.App.Event.Emit("login", "False", "Make sure your Cookie not Expired and Existed")
	} else {
		a.App.Event.Emit("login", "True")
	}
	return true
}
