package test

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"main/internal/pixivlib/DAO"
	"testing"
)

var js = `
{
    "contents": [
        {
            "title": "週３で性感マッサージに通うエリート婦警",
            "date": "2025年09月15日 17:34",
            "tags": [
                "R-18",
                "オリジナル",
                "漫画",
                "オリジナル漫画",
                "女性",
                "茶髪ポニーテール",
                "警察官",
                "女性警察官",
                "婦警",
                "オリジナル5000users入り"
            ],
            "url": "https://i.pximg.net/c/240x480/img-master/img/2025/09/15/17/34/06/135133400_p0_master1200.jpg",
            "illust_type": "0",
            "illust_book_style": "0",
            "illust_page_count": "5",
            "user_name": "ぽりうれたん",
            "profile_img": "https://i.pximg.net/user-profile/img/2017/10/17/01/44/21/13351425_5e525f93a2bdcdf383423e3a2df3da3d_50.jpg",
            "illust_content_type": {
                "sexual": 2,
                "lo": false,
                "grotesque": false,
                "violent": false,
                "homosexual": false,
                "drug": false,
                "thoughts": false,
                "antisocial": false,
                "religion": false,
                "original": true,
                "furry": false,
                "bl": false,
                "yuri": false
            },
            "illust_series": false,
            "illust_id": 135133400,
            "width": 4299,
            "height": 6071,
            "user_id": 19417472,
            "rank": 1,
            "yes_rank": 3,
            "rating_count": 2693,
            "view_count": 108495,
            "illust_upload_timestamp": 1757925246,
            "attr": "original",
            "is_masked": false,
            "is_bookmarked": false,
            "bookmarkable": true
        },
        {
            "title": "カラマネロ様だいすきメイちゃん",
            "date": "2025年09月15日 21:00",
            "tags": [
                "R-18",
                "メイ(トレーナー)",
                "ポケモン人間絵",
                "催眠",
                "巨根",
                "メイぱい",
                "体格差",
                "中出し",
                "ポケモン3000users入り"
            ],
            "url": "https://i.pximg.net/c/240x480/img-master/img/2025/09/15/21/00/09/135141929_p0_master1200.jpg",
            "illust_type": "0",
            "illust_book_style": "0",
            "illust_page_count": "2",
            "user_name": "超ジロー🔞",
            "profile_img": "https://i.pximg.net/user-profile/img/2023/06/04/22/58/25/24506125_37d7a87f0f871efecfe54e2c18545802_50.jpg",
            "illust_content_type": {
                "sexual": 2,
                "lo": false,
                "grotesque": false,
                "violent": false,
                "homosexual": false,
                "drug": false,
                "thoughts": false,
                "antisocial": false,
                "religion": false,
                "original": false,
                "furry": false,
                "bl": false,
                "yuri": false
            },
            "illust_series": false,
            "illust_id": 135141929,
            "width": 1166,
            "height": 1651,
            "user_id": 39575727,
            "rank": 2,
            "yes_rank": 4,
            "rating_count": 3572,
            "view_count": 70761,
            "illust_upload_timestamp": 1757937609,
            "attr": "",
            "is_masked": false,
            "is_bookmarked": false,
            "bookmarkable": true
        },
	]
}`

func TestParseJson(t *testing.T) {
	data := gjson.ParseBytes([]byte(js)).Get("contents")
	var tmp DAO.RankData
	all := data.Array()
	for _, value := range all {
		json.Unmarshal([]byte(value.Raw), &tmp)
		t.Log(tmp.IllustContentType.Sexual)
	}
}
