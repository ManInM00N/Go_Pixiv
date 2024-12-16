package init

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	. "main/backend/src/DAO"
	"net/http"
	url2 "net/url"
	"os"
	"strconv"
	"time"

	"github.com/ManInM00N/go-tool/statics"
	"github.com/tidwall/gjson"
	"github.com/yuin/goldmark/util"
)

const (
	IllustInfo = 1 << iota
	AuthorInfo
	RankInfo
	FollowInfo
	GifPage
	NovelInfo
	UserDashboard = 0
)

// TODO: 作者全部作品下载OK
// TODO: 基础下载 OK   目录管理下载 OK  主要图片全部下载OK    并发下载OK
// TODO: 指针内存问题OK
// TODO: 图片下载完整  OK
func Download(i *Illust, op *Option) bool {
	var err error
	total := 0
	Request, err2 := http.NewRequest("GET", i.PreviewImageUrl, nil)
	clientcopy := GetClient()
	if err2 != nil {
		DebugLog.Println("Error creating request", err2)
		return false
	}
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	Request.Header.Set("referer", "https://www.pixiv.net")
	Cookie := &http.Cookie{
		Name:  "cookie",
		Value: Setting.Cookie,
	}
	Request.AddCookie(Cookie)
	Request.Header.Set("cookie", "PHPSESSID="+Setting.Cookie)
	var Response *http.Response
	defer func() {
		if Response != nil {
			Response.Body.Close()
		}
	}()
	_, err = os.Stat(Setting.Downloadposition)
	if err != nil {
		os.Mkdir(Setting.Downloadposition, os.ModePerm)
	}
	Path := Setting.Downloadposition
	if op.Mode == ByRank {
		Path = Path + "/" + op.Rank + op.RankDate
		_, err = os.Stat(Path)
		if err != nil {
			os.Mkdir(Path, os.ModePerm)
		}
	}
	if op.DiffAuthor {
		Path = Path + "/" + statics.Int64ToString(i.UserID)
		_, err = os.Stat(Path)
		if err != nil {
			os.Mkdir(Path, os.ModePerm)
		}
		// Path = AuthorFile
	}

	Type := Path + "/" + i.AgeLimit
	_, err = os.Stat(Type)
	if err != nil {
		os.Mkdir(Type, os.ModePerm)
	}
	failtimes := 0
	for j := 0; j < i.Pages; j++ {
		imagefilename := statics.GetFileName(i.ImageUrl[j])
		imagefilepath := Type + "/" + imagefilename
		img, err2 := os.Stat(imagefilepath)
		if err2 == nil {
			if op.Mode == ByPid {
				os.Remove(imagefilepath)
			} else if img.Size() != 0 {
				time.Sleep(time.Millisecond * time.Duration(Setting.Downloadinterval))
				continue
			}
		}
		Request.URL, _ = url2.Parse(i.ImageUrl[j])
		ok := true
		for k := 0; k < 10; k++ {
			Response, err = clientcopy.Do(Request)
			if k == 9 && err != nil {
				DebugLog.Println("Illust Resouce Request Error", err, Request.URL.String())
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
			time.Sleep(time.Millisecond * time.Duration(Setting.Downloadinterval))
		}
		if !ok {
			os.Remove(imagefilepath)
			continue
		}
		failtimes = 0
		f, err := os.OpenFile(imagefilepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			DebugLog.Println(i.Pid, "Download Failed", err, "retrying")
			os.Remove(imagefilepath)
			j--
			continue
		}
		bufWriter := bufio.NewWriter(f)
		_, err = io.Copy(bufWriter, Response.Body)
		if err != nil {
			DebugLog.Println(i.Pid, " Write Failed", err)
			return false
		}
		f.Close()
		bufWriter.Flush()
		total++
		time.Sleep(time.Millisecond * time.Duration(Setting.Downloadinterval))
	}
	return true
}

// return url & referer
func CheckMode(url, id string, num int) (string, string) {
	// switch num{
	// case 1:
	// case 2:
	// case
	// }
	if num == 1 { // illust page
		return "https://www.pixiv.net/ajax/illust/" + url, "https://www.pixiv.net/artworks/" + id
	} else if num == 2 { // author page
		return "https://www.pixiv.net/ajax/user/" + url + "/profile/all", "https://www.pixiv.net/member.php?id=" + id
	} else if num == 4 { // ranking page
		return "https://www.pixiv.net/ranking.php?format=json" + url, "https://www.pixiv.net/"
	} else if num == 8 { // follow page
		return "https://www.pixiv.net/ajax/follow_latest/illust?" + url, "https://www.pixiv.net/"
	}
	return "https://www.pixiv.net/ajax/user/extra", "https://www.pixiv.net/"
}

// TODO:下载作品主题信息json OK
func GetWebpageData(url, id string, num int) ([]byte, error) { // 请求得到作品json

	var response *http.Response
	var err error
	ur, ref := CheckMode(url, id, num)
	// println(ur, ref)
	Request, err := http.NewRequest("GET", ur, nil)
	if err != nil {
		DebugLog.Println("Error creating request", err)
		return nil, err
	}
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")
	Request.Header.Set("referer", ref)
	// Cookie := &http.Cookie{
	// 	Name:  "PHPSESSID",
	// 	Value: Setting.Cookie,
	// }
	// Request.AddCookie(Cookie)
	// if num != 4 {
	Request.Header.Set("Cookie", "PHPSESSID="+Setting.Cookie)
	// }

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
			DebugLog.Println("Request failed ", err)
			return nil, err
		}
		time.Sleep(time.Duration(Setting.Retryinterval) * time.Millisecond)

	}
	defer response.Body.Close()

	// webpageBytes, err3 := ioutil.ReadAll(response.Body)
	var buffer bytes.Buffer
	reader := bufio.NewReader(response.Body)
	_, err3 := io.Copy(&buffer, reader)
	webpageBytes := buffer.Bytes()
	if err3 != nil {
		DebugLog.Println("read failed", err3)
		return nil, err3
	}
	if response.StatusCode != http.StatusOK {
		DebugLog.Println("status code ", response.StatusCode, ur, string(webpageBytes))
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
	data, err2 := GetWebpageData(urltail, strid, 1)

	if err2 != nil {
		err = fmt.Errorf("GetWebpageData error %w", err2)
		DebugLog.Println("GetWebpageData error", err2)
		return nil, err
	}
	Results := gjson.ParseBytes(data)
	canbedownload := Results.Get("error").Bool()
	if canbedownload {
		// println(strid, len(strid))
		return nil, &NotGood{}
	}
	jsonmsg := gjson.ParseBytes(data).Get("body") // 读取json内作品及作者id信息
	// println(id, jsonmsg.Str)
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
	}
	for _, tag := range jsonmsg.Get("tags.tags.#.tag").Array() {
		i.Tags = append(i.Tags, tag.Str)
		if tag.Str == "R-18" {
			i.AgeLimit = "r18"
			break
		}
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
	pages, err2 := GetWebpageData(urltail+"/pages", strid, 1)
	if err2 != nil {
		err = fmt.Errorf("Get illustpage data error %w", err2)
		DebugLog.Println("get illustpage data error", err2)
		return nil, err
	}
	imagejson := gjson.ParseBytes(pages).Get("body").String()
	var imagedata []ImageData
	err2 = json.Unmarshal(util.StringToReadOnlyBytes(imagejson), &imagedata)
	if err2 != nil {
		err = fmt.Errorf("error decoding %w", err2)
		DebugLog.Println("Error decoding", err2)
		return nil, err
	}
	if len(imagedata) == 0 {
		DebugLog.Println("No images found ", i.Pid)
		return nil, fmt.Errorf("No images found", i.Pid)
	}
	for _, image := range imagedata {
		i.ImageUrl = append(i.ImageUrl, image.URLs.Original)
	}

	return i, err
}

func GetAuthor(id int64) (map[string]gjson.Result, error) {
	data, err := GetWebpageData(strconv.FormatInt(id, 10), strconv.FormatInt(id, 10), 2)
	if err != nil {
		return nil, err
	}
	jsonmsg := gjson.ParseBytes(data).Get("body")
	ss := jsonmsg.Get("illusts").Map()
	return ss, nil
}

func GetRank(option *Option) (gjson.Result, error) {
	option.Msg()
	// DebugLog.Println("https://www.pixiv.net/ranking.php?format=json" + option.Suffix)
	data, err := GetWebpageData(option.Suffix, "", 4)
	if err != nil {
		// println("get failed: ", err.Error())
		return gjson.Result{}, err
	}
	// arr := gjson.ParseBytes(data).Get("contents.#.illust_id")
	arr := gjson.ParseBytes(data).Get("contents")
	return arr, nil
}

func GetFollow(option *Option) (gjson.Result, error) {
	option.Msg()
	// https://www.pixiv.net/ajax/follow_latest/illust?&mode=all&p=1
	InfoLog.Println("https://www.pixiv.net/ajax/follow_latest/illust?" + option.Suffix)
	data, err := GetWebpageData(option.Suffix, "", 8)
	if err != nil {
		DebugLog.Println("get failed: ", err.Error())
		return gjson.Result{}, err
	}
	arr := gjson.ParseBytes(data).Get("body")
	// DebugLog.Println(arr.Get("thumbnails"))
	//for i := range arr {
	//	println(arr[i].Int())
	//}
	return arr, nil
}

func JustDownload(pid string, mode *Option) (int, bool) {
	illust, err := work(statics.StringToInt64(pid), mode)
	if ContainMyerror(err) {
		DebugLog.Println(err)
		if !mode.OnlyPreview {
			return 0, true
		}
	}
	if illust == nil {
		DebugLog.Println(pid, " Download failed")
		return 0, false
	}
	if mode.ShowSingle {
		InfoLog.Println(pid + " Start download")
		defer InfoLog.Println(pid + " Finished download")
	}
	if !Download(illust, mode) {
		DebugLog.Println(pid, " Download failed")
		return 0, false
	}
	return 1, true
}
