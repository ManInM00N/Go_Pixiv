package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ManInM00N/go-tool/statics"
	"github.com/bmaupin/go-epub"
	"github.com/tidwall/gjson"
	"io"
	. "main/configs"
	"main/internal/cache/DAO"
	. "main/internal/pixivlib/DAO"
	"main/pkg/utils"
	"net/http"
	url2 "net/url"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const (
	IllustInfo  = 1
	IllustPages = iota + 1
	AuthorArtworks
	AuthorInfo
	RankInfo
	FollowInfo
	GifPage
	NovelInfo
	SeriesInfo
	NovelSeriesList
	NovelText
	UserDashboard
	PicSource
	FollowNovelInfo
	Base = "https://www.pixiv.net/"
)

// 获取下载路径
func getDownloadPath(i *Illust, op *Option, setting *Settings) string {
	path := setting.PixivConf.Downloadposition

	switch {
	case op.Mode == ByRank:
		path = filepath.Join(path, op.Rank+op.RankDate)
	case i.IllustType == UgoiraType:
		path = filepath.Join(path, "GIF")
	case op.DiffAuthor || op.Mode == ByAuthor:
		path = filepath.Join(path, statics.Int64ToString(i.UserID))
	}

	return path
}

func Download(i *Illust, op *Option) bool {

	setting := NowSetting()
	client := GetClient()

	path := getDownloadPath(i, op, &setting)
	typePath := filepath.Join(path, i.AgeLimit)
	if err := os.MkdirAll(typePath, os.ModePerm); err != nil {
		utils.DebugLog.Println("Failed to create directory:", err)
		return false
	}

	Request, err := CreatePixivRequest(i.Source, &setting)
	if err != nil {
		utils.DebugLog.Println("Error creating request", err)
		return false
	}

	// 根据类型下载
	var success bool
	if i.IllustType <= 1 {
		success = downloadImages(i, op, typePath, Request, client, &setting)
	} else {
		success = downloadUgoira(i, path, Request, client, &setting)
	}

	if success {
		saveCache(i)
	}

	return success
}

// return url & referer
func GetUrlRefer(url, id string, num int) (string, string) {
	switch num {
	case IllustInfo:
		return Base + "ajax/illust/" + url, Base + "artworks/" + id
	case IllustPages:
		return Base + "ajax/illust/" + url + "/pages", Base + "artworks/" + id
	case AuthorArtworks:
		return Base + "ajax/user/" + url + "/profile/all", Base + "member.php?id=" + id
	case AuthorInfo:
		return Base + "ajax/user/" + url + "?full=1", Base
	case RankInfo:
		return Base + "ranking.php?format=json" + url, Base
	case FollowInfo:
		return Base + "ajax/follow_latest/illust?" + url, Base
	case FollowNovelInfo:
		return Base + "ajax/follow_latest/novel?" + url, Base
	case GifPage:
		return Base + "ajax/illust/" + id + "/ugoira_meta", Base
	case NovelInfo:
		return Base + "ajax/novel/" + id, Base
	case SeriesInfo:
		return Base + "ajax/novel/series/" + id, Base
	case NovelSeriesList:
		return Base + "ajax/novel/series_content/" + id, Base
	case NovelText:
		return Base + "novel/show.php?id=" + id, Base
	case UserDashboard:
		return Base + "ajax/user/extra", Base
	default:
		return Base + "ajax/user/extra", Base
	}
}

// 下载普通图片
func downloadImages(i *Illust, op *Option, typePath string, request *http.Request, client *http.Client, setting *Settings) bool {
	failTimes := 0

	for j := 0; j < i.Pages; j++ {
		imageFilename := statics.GetFileName(i.ImageUrl[j])
		imageFilepath := filepath.Join(typePath, imageFilename)

		// 检查文件是否已存在
		if shouldSkipDownload(imageFilepath, op) {
			time.Sleep(time.Millisecond * time.Duration(setting.PixivConf.Downloadinterval))
			continue
		}

		// 下载单个图片
		if !downloadSingleImage(i.ImageUrl[j], imageFilepath, request, client, setting) {
			failTimes++
			if failTimes > 2 {
				continue
			}
			j-- // 重试
			continue
		}

		failTimes = 0
		time.Sleep(time.Millisecond * time.Duration(setting.PixivConf.Downloadinterval))
	}

	return true
}

// 下载 Ugoira（GIF）
func downloadUgoira(i *Illust, path string, request *http.Request, client *http.Client, setting *Settings) bool {
	// 获取数据
	var response *http.Response
	var err error
	for k := 0; k < 10; k++ {
		response, err = client.Do(request)
		if err == nil {
			break
		}
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
		if k == 9 {
			utils.DebugLog.Println("Ugoira request error:", err)
			return false
		}
		time.Sleep(time.Millisecond * time.Duration(setting.PixivConf.Downloadinterval))
	}

	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		utils.DebugLog.Println("Ugoira response read failed:", err)
		return false
	}

	data, err := ExtractAndProcessFrames(body, i.Frames, int(i.Width), int(i.Height))

	if err != nil {
		utils.DebugLog.Println("GIF encode failed:", err)
		return false
	}

	// 保存 GIF
	gifPath := filepath.Join(path, statics.Int64ToString(i.Pid)+".gif")
	if err := utils.WriteBytesToFile(gifPath, data); err != nil {
		utils.DebugLog.Println("GIF write failed:", err)
		return false
	}

	data = nil
	runtime.GC()
	debug.FreeOSMemory() // 归还内存给OS

	return true
}

// 检查是否应该跳过下载
func shouldSkipDownload(filepath string, op *Option) bool {
	img, err := os.Stat(filepath)
	if err != nil {
		return false
	}

	if op.Mode == ByPid {
		os.Remove(filepath)
		return false
	}

	return img.Size() != 0
}

// 下载单个图片（修复内存泄漏）
func downloadSingleImage(url, filepath string, request *http.Request, client *http.Client, setting *Settings) bool {
	request.URL, _ = url2.Parse(url)

	// 重试逻辑
	var response *http.Response
	var err error
	for k := 0; k < 10; k++ {
		response, err = client.Do(request)

		// ⭐ 关键修复：确保每次请求的 Response.Body 都被关闭
		if err != nil {
			if response != nil && response.Body != nil {
				response.Body.Close()
			}
			if k == 9 {
				utils.DebugLog.Println("Image request error:", err, url)
				return false
			}
			time.Sleep(time.Millisecond * time.Duration(setting.PixivConf.Downloadinterval))
			continue
		}
		break
	}

	// ⭐ 使用 defer 确保 Body 被关闭
	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()

	// 写入文件
	if err := utils.WriteToFile(filepath, response.Body); err != nil {
		utils.DebugLog.Println("Write failed:", err)
		os.Remove(filepath)
		return false
	}

	return true
}

func DownloadNovel(id string) bool {
	setting := NowSetting()
	body, err := GetNovel(id)
	if err != nil {
		utils.DebugLog.Println("GetWebpageData error", err)
		return false
	}
	novel := Novel{
		Id:          body.Get("id").String(),
		Content:     body.Get("content").String(),
		Title:       body.Get("title").String(),
		UserId:      body.Get("userId").String(),
		UserName:    body.Get("userName").String(),
		CoverUrl:    body.Get("coverUrl").String(),
		SeriesId:    int(body.Get("seriesNavData.seriesId").Int()),
		SeriesTitle: body.Get("seriesNavData.title").String(),
	}
	for _, v := range body.Get("tags.tags.#.tag").Array() {
		novel.Tags = append(novel.Tags, v.String())
	}
	novel.R18 = utils.HasR18(&novel.Tags)
	Path := filepath.Join(setting.PixivConf.Downloadposition, "novel")
	if novel.R18 {
		Path = filepath.Join(Path, "r18")
	} else {
		Path = filepath.Join(Path, "all-age")
	}
	if novel.SeriesId != 0 {
		Path = filepath.Join(Path, statics.IntToString(novel.SeriesId))
	}
	os.MkdirAll(Path, os.ModePerm)
	title := utils.Cut(novel.Id + novel.Title)
	f, err := os.OpenFile(filepath.Join(Path, title+".txt"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		utils.DebugLog.Println("open file failed", err)
		return false
	}
	bufWriter := bufio.NewWriter(f)
	_, err = bufWriter.WriteString(novel.Content)
	if f != nil {
		f.Close()
	}
	if err != nil {
		utils.DebugLog.Println("write file failed ", err)
		return false
	}
	e := epub.NewEpub(title)
	time.Sleep(time.Millisecond * time.Duration(setting.PixivConf.Downloadinterval))
	client := GetClient()
	Request, err := http.NewRequest("GET", novel.CoverUrl, nil)
	if err != nil {
		utils.DebugLog.Println("Error creating request", err)
		return false
	}
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")
	Request.Header.Set("Cookie", "PHPSESSID="+setting.PixivConf.Cookie)
	Request.Header.Set("Referer", "https://www.pixiv.net/")
	res, err := client.Do(Request)
	if err != nil {
		utils.DebugLog.Println("novel request failed ", err)
		return false
	}
	reader := bufio.NewReader(res.Body)
	imgpath := filepath.Join(Path, novel.Id+"_cover."+utils.GetFileType(novel.CoverUrl))
	f, err = os.OpenFile(imgpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		utils.DebugLog.Println(imgpath, err)
		return false
	}
	bufWriter = bufio.NewWriter(f)
	_, err = io.Copy(bufWriter, reader)
	if err != nil {
		utils.DebugLog.Println("write cover failed ", err)
		return false
	}
	f.Close()
	coverimg, _ := e.AddImage(imgpath, "cover.jpg")
	e.SetCover(coverimg, "")
	novelTextHTML := "<p>" + strings.ReplaceAll(novel.Content, "\n\n", "</p><p>") + "</p>"

	_, err = e.AddSection(novelTextHTML, novel.Title, "", "")
	if err != nil {
		utils.DebugLog.Printf("无法添加章节: %v\n", err)
		return false
	}
	err = e.Write(filepath.Join(Path, title+".epub"))
	if err != nil {
		utils.DebugLog.Printf("无法保存 %s EPUB文件: %v\n", novel.Id, err)
		return false
	}
	os.Remove(imgpath)
	return true
}

func GetWebpageData(url, id string, num int) ([]byte, error) { // 请求得到作品json

	ur, _ := GetUrlRefer(url, id, num)
	setting := NowSetting()

	Request, err := CreatePixivRequest(ur, &setting)
	if err != nil {
		utils.DebugLog.Println("Error creating request:", err)
		return nil, err
	}

	clientcopy := GetClient()
	var response *http.Response
	for i := 0; i < 10; i++ {
		response, err = clientcopy.Do(Request)
		if err == nil {
			if response.StatusCode == 429 {
				if response.Body != nil {
					response.Body.Close()
				}
				time.Sleep(time.Duration(setting.PixivConf.Retry429) * time.Millisecond)
				continue
			}
			break
		}

		if response != nil && response.Body != nil {
			response.Body.Close()
		}

		if i == 9 {
			utils.DebugLog.Println("Request failed:", err)
			return nil, err
		}

		time.Sleep(time.Duration(setting.PixivConf.Retryinterval) * time.Millisecond)
	}
	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, response.Body); err != nil {
		utils.DebugLog.Println("Read failed:", err)
		return nil, err
	}

	webpageBytes := buffer.Bytes()

	if response.StatusCode != http.StatusOK {
		utils.DebugLog.Println("status code ", response.StatusCode, ur, string(webpageBytes))
		if response.StatusCode == 429 {
			time.Sleep(time.Duration(setting.PixivConf.Retry429) * time.Millisecond)
			return nil, &TooFastRequest{S: "TooMuchRequest in a short period", Err: errors.New("TooMuchRequest")}
		}
	}
	return webpageBytes, nil
}

func work(id int64, mode *Option) (i *Illust, err error) { // 按作品id查找
	urltail := strconv.FormatInt(id, 10)
	strid := urltail
	err = nil
	data, err2 := GetWebpageData(urltail, strid, IllustInfo)
	if err2 != nil {
		err = fmt.Errorf("GetWebpageData error %w", err2)
		utils.DebugLog.Println("GetWebpageData error", err2)
		return nil, err
	}
	Results := gjson.ParseBytes(data)
	canbedownload := Results.Get("error").Bool()
	if canbedownload {
		return nil, NotFound
	}

	jsonmsg := gjson.ParseBytes(data).Get("body") // 读取json内作品及作者id信息
	i = &Illust{
		AgeLimit:    "all-age",
		Pid:         jsonmsg.Get("illustId").Int(),
		UserID:      jsonmsg.Get("userId").Int(),
		Caption:     jsonmsg.Get("alt").Str,
		CreatedTime: jsonmsg.Get("createDate").Str,
		Pages:       int(jsonmsg.Get("pageCount").Int()),
		Title:       jsonmsg.Get("illustTitle").Str,
		UserName:    jsonmsg.Get("userName").Str,
		Likecount:   int(jsonmsg.Get("bookmarkCount").Int()),
		IllustType:  int(jsonmsg.Get("illustType").Int()),
		Width:       jsonmsg.Get("width").Int(),
		Height:      jsonmsg.Get("height").Int(),
	}
	for _, tag := range jsonmsg.Get("tags.tags.#.tag").Array() {
		i.Tags = append(i.Tags, tag.Str)
	}
	if utils.HasR18(&i.Tags) {
		i.AgeLimit = "R-18"
	}
	if i.Likecount < mode.Likelimit {
		err = fmt.Errorf("%w", &NotGood{S: "LikeNotEnough", Err: errors.New("LikeNotEnough")})
	}
	if i.AgeLimit == "r18" && !mode.R18 {
		err = fmt.Errorf("%w", &AgeLimit{S: "AgeLimitExceed", Err: errors.New("AgeLimitExceed")})
	}
	if mode.OnlyPreview {
		i.PreviewImageUrl = jsonmsg.Get("urls.small").String()
		// DebugLog.Println(i.PreviewImageUrl)
		return i, err
	}
	if i.IllustType <= 1 {
		pages, err2 := GetWebpageData(urltail+"/pages", strid, IllustInfo)
		if err2 != nil {
			utils.DebugLog.Println("get illustpage data error", err2)
			return nil, err
		}
		imagejson := gjson.ParseBytes(pages).Get("body").String()
		var imagedata []ImageData

		err2 = json.Unmarshal([]byte(imagejson), &imagedata)
		if err2 != nil {
			utils.DebugLog.Println("Error decoding", err2, imagejson)
			return nil, err
		}
		if len(imagedata) == 0 {
			utils.DebugLog.Println("No images found ", i.Pid)
			return nil, fmt.Errorf("No images found %d", i.Pid)
		}
		for _, image := range imagedata {
			i.ImageUrl = append(i.ImageUrl, image.URLs.Original)
		}
	} else {
		data, err2 := GetWebpageData(urltail, strid, GifPage)
		if err2 != nil {
			utils.DebugLog.Println("get ugoira data error", err2)
			return nil, err
		}
		jsonbody := gjson.ParseBytes(data).Get("body")
		i.Source = jsonbody.Get("originalSrc").String()
		i.FileType = jsonbody.Get("mime_type").String()
		if err2 = json.Unmarshal([]byte(jsonbody.Get("frames").String()), &i.Frames); err2 != nil {
			return nil, err2
		}

	}

	return i, err
}

func JustDownload(pid string, mode *Option, callEvent func(name string, data ...interface{})) (int, bool) {
	illust, err := work(statics.StringToInt64(pid), mode)
	if ContainMyerror(err) {
		utils.DebugLog.Println(err)
		if !mode.OnlyPreview {
			return 0, true
		}
	}
	if errors.Is(err, NotFound) {
		utils.InfoLog.Println(pid, err.Error())
		callEvent("NotFound", pid+" "+err.Error())
		return 0, false
	}
	if illust == nil {
		utils.DebugLog.Println(pid, " Download failed")
		return 0, false
	}
	if mode.ShowSingle {
		utils.InfoLog.Println(pid + " Start download")
		defer utils.InfoLog.Println(pid + " Finished download")
	}
	if !Download(illust, mode) {
		utils.DebugLog.Println(pid, " Download failed")
		return 0, false
	}
	return 1, true
}

// 保存缓存
func saveCache(i *Illust) {
	cache := DAO.Cache{
		DownloadID: statics.Int64ToString(i.Pid),
		Type:       "Illust",
		CreatedAt:  time.Now(),
	}
	if i.IllustType == UgoiraType {
		cache.Type = "Ugoira"
	}
	DAO.Db.FirstOrCreate(&cache, DAO.Cache{DownloadID: cache.DownloadID})
}
