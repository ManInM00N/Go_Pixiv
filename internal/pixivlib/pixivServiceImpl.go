package pixivlib

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/ManInM00N/go-tool/goruntine"
	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails/v3/pkg/application"
	"io"
	"main/internal/imageService"

	"main/configs"
	"main/internal/pixivlib/DAO"
	. "main/internal/pixivlib/handler"
	"main/internal/taskQueue"
	. "main/pkg/utils"
	"net/http"
	"time"
)

type Ctl struct {
	ctx    context.Context
	App    *application.App
	cancel func()
}

func NewCtl() *Ctl {
	ctx, cancel := context.WithCancel(context.Background())
	return &Ctl{ctx: ctx, cancel: cancel}
}
func (a *Ctl) RegisterService(app *application.App) {
	a.App = app
	app.RegisterService(application.NewService(a))
	t := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-a.ctx.Done():
				t.Stop()
				goto end
				//break
			case <-t.C:
				arr, worker := taskQueue.TaskPool.GetTaskStatistic()
				arr2, worker2 := taskQueue.P.GetTaskStatistic()
				a.App.Event.Emit("taskPoolInfos", arr, worker, arr2, worker2)
			}
		}
	end:
		DebugLog.Println("Infos end")
		fmt.Println("Ticker End", time.Now())
	}()
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

func (a *Ctl) DownloadByNovelId(text string, isSeries bool) bool {
	InfoLog.Println("Download Novel ", text)
	if isSeries {
		Download_By_SeriesId(text, func(name string, data ...interface{}) {
			a.App.Event.Emit(name, data)
		})
	} else {
		Download_By_NovelId(text)
	}
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
	Request.Header.Set("Cookie", "PHPSESSID="+configs.NowSetting().PixivConf.Cookie)
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

func (a *Ctl) RemoveTask(taskId string) bool {
	taskQueue.TaskPool.RemoveTaskByDeleteFunc(
		func(item goruntine.Task) bool {
			return (*item.GetInfo()).(*DAO.TaskInfo).ID == taskId
		})
	return true
}

const str = `
	const cookies = document.cookie;
	console.log(cookies)
    if ('cookieStore' in window) {
        const allCookies = await cookieStore.getAll();
        console.log(allCookies);
    }
`

func (a *Ctl) GetWebView2Cookies(windowName string) {
	// 获取 WebView2 控制器
	App := a.App
	windows := App.Window.GetAll()
	fmt.Println(windowName, len(windows))
	for _, window := range windows {
		window.ExecJS(str)
		fmt.Println(window.Name())
		//App.Browser.OpenURL("https://www.pixiv.net/users/114572298")
	}

}

func (a *Ctl) OpenInBrowser(url string) {
	err := a.App.Browser.OpenURL(url)
	if err != nil {
		DebugLog.Println(url, err.Error())
	}
}

func (a *Ctl) FetchRecentLog() {
	res, err := ReadRecentLogs(LogPositon, 10)
	if err != nil {
		DebugLog.Println(res, err.Error())
		return
	}
	a.App.Event.Emit("FetchLogs", res)
	return
}

func (a *Ctl) OpenFileFolder() string {

	tmp := &application.OpenFileDialogOptions{}
	tmp.CanChooseFiles = true
	tmp.CanChooseDirectories = true
	tmp.Title = "选择一个目录 Choose a folder"
	tmp.CanCreateDirectories = true
	dia := a.App.Dialog.OpenFileWithOptions(tmp)
	selection, err := dia.PromptForSingleSelection()
	if err != nil {
		return ""
	}
	return selection
}

func (a *Ctl) GetSauceNAOQuota() map[string]interface{} {
	res, err := imageService.SauceNaoService.GetAPIQuota()
	if err != nil {
		ErrorLog.Printf("%v", err)
	}
	return res
}
