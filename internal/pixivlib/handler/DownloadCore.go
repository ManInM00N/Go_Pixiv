package handler

import (
	Map "github.com/ManInM00N/go-tool/sync"
	"github.com/tidwall/gjson"
	_ "image/png"
	. "main/internal/pixivlib/DAO"
	"main/pkg/utils"
	"strconv"
	"time"
)

var (
	Logo             []byte
	NowTaskMsg       = ""
	QueueTaskMsg     = ""
	ProcessMax       = int64(0)
	ProcessNow       = int64(0)
	RankLoadingNow   = false
	FollowLoadingNow = false
	UgoiraMap        = Map.NewRWMap[string, bool]()
)

func UgoiraDownloadWait(indentify string) {
	var v, ok bool
	for {
		v, ok = UgoiraMap.Get(indentify)
		if ok && v == true {
			UgoiraMap.Delete(indentify)
			break
		} else if ok == false {
			break
		}
		time.Sleep(time.Second)
	}
}

func GetAuthorArtworks(id int64) (map[string]gjson.Result, error) {
	data, err := GetWebpageData(strconv.FormatInt(id, 10), strconv.FormatInt(id, 10), AuthorArtworks)
	if err != nil {
		return nil, err
	}
	jsonmsg := gjson.ParseBytes(data).Get("body")
	ss := jsonmsg.Get("illusts").Map()
	return ss, nil
}

func GetAuthorInfo(id int64) (map[string]gjson.Result, error) {
	data, err := GetWebpageData(strconv.FormatInt(id, 10), strconv.FormatInt(id, 10), AuthorInfo)
	if err != nil {
		return nil, err
	}
	jsonmsg := gjson.ParseBytes(data).Get("body").Map()
	return jsonmsg, nil
}

func GetRank(option *Option) (gjson.Result, error) {
	option.Msg()
	// DebugLog.Println("https://www.pixiv.net/ranking.php?format=json" + option.Suffix)
	data, err := GetWebpageData(option.Suffix, "", RankInfo)
	if err != nil {
		// println("get failed: ", err.Error())
		return gjson.Result{}, err
	}
	// arr := gjson.ParseBytes(data).Get("contents.#.illust_id")
	arr := gjson.ParseBytes(data).Get("contents")
	utils.DebugLog.Println(arr)
	return arr, nil
}

func GetFollow(option *Option) (gjson.Result, error) {
	option.Msg()
	// https://www.pixiv.net/ajax/follow_latest/illust?&mode=all&p=1
	utils.InfoLog.Println("https://www.pixiv.net/ajax/follow_latest/illust?" + option.Suffix)
	var data []byte
	var err error
	if option.FollowType == "illust" {
		data, err = GetWebpageData(option.Suffix, "", FollowInfo)
	} else if option.FollowType == "novel" {
		data, err = GetWebpageData(option.Suffix, "", FollowNovelInfo)
	} else {

	}
	if err != nil {
		utils.DebugLog.Println("get failed: ", err.Error())
		return gjson.Result{}, err
	}
	arr := gjson.ParseBytes(data).Get("body")
	return arr, nil
}

func GetNovel(id string) (gjson.Result, error) {
	data, err := GetWebpageData(id, id, NovelInfo)
	if err != nil {
		utils.DebugLog.Println("get novel html failed: ", err.Error())
		return gjson.Result{}, err
	}
	v := gjson.ParseBytes(data).Get("body")
	return v, nil
}

func GetSeriesInfo(id string) (gjson.Result, error) {
	data, err := GetWebpageData(id, id, SeriesInfo)
	if err != nil {
		utils.DebugLog.Println("get novel html failed: ", err.Error())
		return gjson.Result{}, err
	}
	v := gjson.ParseBytes(data).Get("body")
	return v, nil
}

func GetNovelSeries(id string) (gjson.Result, error) {
	data, err := GetWebpageData(id, id, NovelSeriesList)
	if err != nil {
		utils.DebugLog.Println("get novel html failed: ", err.Error())
		return gjson.Result{}, err
	}
	v := gjson.ParseBytes(data).Get("body").Get("page").Get("seriesContents")
	return v, nil
}
