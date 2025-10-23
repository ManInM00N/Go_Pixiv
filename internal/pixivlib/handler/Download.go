package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ManInM00N/go-tool/statics"
	"github.com/ManInM00N/nicogif"
	"github.com/bmaupin/go-epub"
	"github.com/klauspost/compress/zip"
	"github.com/tidwall/gjson"
	"image"
	_ "image/jpeg"
	_ "image/png"
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

// TODO: 作者全部作品下载OK
// TODO: 基础下载 OK   目录管理下载 OK  主要图片全部下载OK    并发下载OK
// TODO: 指针内存问题OK
// TODO: 图片下载完整  OK
func Download(i *Illust, op *Option) bool {
	//if i.IllustType == UgoiraType {
	//	return true
	//}

	var err error
	total := 0
	// create Request
	Request, err2 := http.NewRequest("GET", i.Source, nil)
	if err2 != nil {
		utils.DebugLog.Println("Error creating request", err2)
		return false
	}
	setting := NowSetting()
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	Request.Header.Set("referer", "https://www.pixiv.net")
	Request.Header.Set("cookie", "PHPSESSID="+setting.Cookie)
	var Response *http.Response
	clientcopy := GetClient()

	Path := setting.Downloadposition
	if op.Mode == ByRank {
		Path = filepath.Join(Path, op.Rank+op.RankDate)
	} else if i.IllustType == UgoiraType {
		Path = filepath.Join(Path, "GIF")
	} else {
		if op.DiffAuthor || op.Mode == ByAuthor {
			Path = filepath.Join(Path, statics.Int64ToString(i.UserID))
		}
	}

	Type := filepath.Join(Path, i.AgeLimit)
	if has, _ := utils.FileExists(Type); !has {
		os.MkdirAll(Type, os.ModePerm)
	}

	cache := DAO.Cache{
		DownloadID: statics.Int64ToString(i.Pid),
		Type:       "Illust",
	}
	failtimes := 0
	if i.IllustType <= 1 {
		for j := 0; j < i.Pages; j++ {
			imagefilename := statics.GetFileName(i.ImageUrl[j])
			imagefilepath := filepath.Join(Type, imagefilename)
			img, err2 := os.Stat(imagefilepath)
			if err2 == nil {
				if op.Mode == ByPid {
					os.Remove(imagefilepath)
				} else if img.Size() != 0 {
					time.Sleep(time.Millisecond * time.Duration(setting.Downloadinterval))
					continue
				}
			}
			Request.URL, _ = url2.Parse(i.ImageUrl[j])
			ok := true
			for k := 0; k < 10; k++ {
				Response, err = clientcopy.Do(Request)
				if k == 9 && err != nil {
					utils.DebugLog.Println("Illust Resouce Request Error", err, Request.URL.String())
					ok = false
					j--
					failtimes++
					if failtimes > 2 {
						j++
					}
					break
				} else if err == nil {
					break
				}
				time.Sleep(time.Millisecond * time.Duration(setting.Downloadinterval))
			}
			if !ok {
				os.Remove(imagefilepath)
				continue
			}
			failtimes = 0
			f, err := os.OpenFile(imagefilepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				utils.DebugLog.Println(i.Pid, "Download Failed", err, "retrying")
				os.Remove(imagefilepath)
				j--
				if f != nil {
					f.Close()
				}
				continue
			}
			bufWriter := bufio.NewWriter(f)

			if _, err = io.Copy(bufWriter, Response.Body); err != nil {
				utils.DebugLog.Println(i.Pid, " Write Failed", err)
				return false
			}
			bufWriter.Flush()
			Response.Body.Close()
			f.Close()
			total++
			time.Sleep(time.Millisecond * time.Duration(setting.Downloadinterval))
		}
	} else {
		for k := 0; k < 10; k++ {
			Response, err = clientcopy.Do(Request)
			if k == 9 && err != nil {
				utils.DebugLog.Println("Illust Resouce Request Error", err, Request.URL.String())
				return false
			} else if err == nil {
				break
			}
			time.Sleep(time.Millisecond * time.Duration(setting.Downloadinterval))
		}
		defer func() {
			if Response != nil {
				Response.Body.Close()
			}
		}()
		body, err := io.ReadAll(Response.Body)
		if err != nil {
			utils.DebugLog.Println("Ugoira Response read failed", err.Error())
			return false
		}

		reader := bytes.NewReader(body)
		zipReader, err := zip.NewReader(reader, int64(len(body)))
		if err != nil {
			utils.DebugLog.Println("Zip reader failed", err.Error())
			return false
		}

		opt := gifencoder.EncodeOptions{
			Width:   int(i.Width),
			Height:  int(i.Height),
			Repeat:  0,
			Quality: 10,
			Dither:  gifencoder.DitherFloydSteinberg,
		}
		var Delay []int
		for _, v := range i.Frames {
			Delay = append(Delay, v.Delay)
		}
		opt.Delays = Delay
		var frames []image.Image
		for _, file := range zipReader.File {
			rc, err := file.Open()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			img, _, err := image.Decode(rc)
			rc.Close()
			if err == nil {
				frames = append(frames, img)
				opt.Width = img.Bounds().Dx()
				opt.Height = img.Bounds().Dy()
			} else {
				fmt.Println(err.Error())
			}
		}
		body = nil      // ⭐ 添加
		reader = nil    // ⭐ 添加
		zipReader = nil // ⭐ 添加

		data, err := gifencoder.EncodeGIFWithOptions(frames, opt)
		frames = nil
		if err != nil {
			utils.DebugLog.Printf("GIF encode failed %s", err.Error())
			return false
		}
		Path = filepath.Join(Path, statics.Int64ToString(i.Pid)+".gif")
		f, err := os.OpenFile(Path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
		bufWriter := bufio.NewWriter(f)
		defer func() {
			bufWriter.Flush()
			f.Close()
		}()
		_, err = bufWriter.Write(data)
		data = nil
		if err != nil {
			utils.DebugLog.Printf("GIF write file failed %s", err.Error())
			return false
		}
		runtime.GC()
		cache.Type = "Ugoria"
	}
	cache.CreatedAt = time.Now()
	DAO.Db.FirstOrCreate(&cache, DAO.Cache{DownloadID: cache.DownloadID})
	return true
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

func DownloadNovel(id string) bool {
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
	Path := filepath.Join(Setting.Downloadposition, "novel")
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
	time.Sleep(time.Millisecond * time.Duration(Setting.Downloadinterval))
	client := GetClient()
	Request, err := http.NewRequest("GET", novel.CoverUrl, nil)
	if err != nil {
		utils.DebugLog.Println("Error creating request", err)
		return false
	}
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")
	Request.Header.Set("Cookie", "PHPSESSID="+Setting.Cookie)
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

	var response *http.Response
	var err error
	ur, ref := GetUrlRefer(url, id, num)
	// println(ur, ref)
	Request, err := http.NewRequest("GET", ur, nil)
	if err != nil {
		utils.DebugLog.Println("Error creating request", err)
		return nil, err
	}
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")
	Request.Header.Set("referer", ref)
	Request.Header.Set("Cookie", "PHPSESSID="+Setting.Cookie)

	clientcopy := GetClient()
	for i := 0; i < 10; i++ {
		response, err = clientcopy.Do(Request)
		if err == nil {
			if response.StatusCode == 429 {
				println("429")
				time.Sleep(time.Duration(Setting.Retry429) * time.Millisecond)
				// i--
				continue
			}
			break
		}
		if i == 9 && err != nil {
			utils.DebugLog.Println("Request failed ", err)
			return nil, err
		}
		time.Sleep(time.Duration(Setting.Retryinterval) * time.Millisecond)

	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	// webpageBytes, err3 := ioutil.ReadAll(response.Body)
	var buffer bytes.Buffer
	reader := bufio.NewReader(response.Body)
	_, err3 := io.Copy(&buffer, reader)
	if err3 != nil {
		utils.DebugLog.Println("read failed", err3)
		return nil, err3
	}
	webpageBytes := buffer.Bytes()
	if response.StatusCode != http.StatusOK {
		utils.DebugLog.Println("status code ", response.StatusCode, ur, string(webpageBytes))
		if response.StatusCode == 429 {
			time.Sleep(time.Duration(Setting.Retry429) * time.Millisecond)
			return nil, &TooFastRequest{S: "TooMuchRequest in a short period", Err: errors.New("TooMuchRequest")}
		}
	}
	return webpageBytes, nil
}

// TODO: 作品信息json请求   OK
// TODO: 多页下载 OK

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
	} else if illust.IllustType == UgoiraType {

		//id, _ := shortid.Generate()
		//identify := statics.Int64ToString(illust.Pid) + id
		//UgoiraMap.Set(identify, false)
		//callEvent("downloadugoira", illust.Pid, illust.Width, illust.Height, illust.Frames, illust.Source, identify)
		//UgoiraDownloadWait(identify)
		//return 1, true
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
