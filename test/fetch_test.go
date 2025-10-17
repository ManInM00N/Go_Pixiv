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

var ct = `
{"seriesContents": [
                {
                    "id": "26166607",
                    "userId": "59261594",
                    "series": {
                        "id": 14662748,
                        "viewableType": 0,
                        "contentOrder": 4
                    },
                    "title": "肉戏追加1：勇者和爱人的十六岁",
                    "commentHtml": "感谢老板投喂的kfc疯狂星期四啊，也是吃上饭了，老板让我抓紧更新，我就把这篇搓出来了",
                    "tags": [
                        "R-18",
                        "中文/中国语/Chinese/中國語/中国語",
                        "纯爱",
                        "后宫",
                        "隐秘做爱/调教/榨精/精液孕肚/激烈性爱",
                        "子宫奸/子宫交/潮吹/连续高潮",
                        "勇者/黛绮丝",
                        "西幻",
                        "内射/中出/受孕/发情/破处"
                    ],
                    "restrict": 0,
                    "xRestrict": 1,
                    "isOriginal": true,
                    "textLength": 3682,
                    "characterCount": 3682,
                    "wordCount": 1629,
                    "useWordCount": false,
                    "readingTime": 552,
                    "bookmarkCount": 409,
                    "url": "https://i.pximg.net/c/150x150_80/novel-cover-master/img/2025/10/03/22/06/02/sci14662748_7a3bf9bc931bb9226fd24c16d2b8dd03_master1200.jpg",
                    "uploadTimestamp": 1760194828,
                    "reuploadTimestamp": 1760211288,
                    "isBookmarkable": true,
                    "bookmarkData": null,
                    "aiType": 1
                },
                {
                    "id": "26156304",
                    "userId": "59261594",
                    "series": {
                        "id": 14662748,
                        "viewableType": 0,
                        "contentOrder": 3
                    },
                    "title": "第三章，偶遇一见钟情的红龙女王与难得一见魅魔萝莉圣女，狠狠喂饱从来没有吃到过精液的纯洁魅魔萝莉，让强行以逆种付位压迫自己的龙娘露出好看（淫荡）的高潮脸吧",
                    "commentHtml": "我看到一张图啊，说是想看什么异种族娘，发现精灵名列前茅，拿这下不得不尝了，你们还想看什么异种族娘",
                    "tags": [
                        "R-18",
                        "中文/中国语/Chinese/中國語/中国語",
                        "纯爱/后宫/西幻",
                        "隐秘做爱/调教/榨精/精液孕肚/激烈性爱",
                        "征服/母狗/肉便器/巨乳肥臀",
                        "子宫奸/子宫交/潮吹/连续高潮/种付位",
                        "中出/内射/发情/白给/逆种付位",
                        "大量精液/精液喷出/受孕/窒息/口爆/深喉",
                        "吞精/吸精/口交/巨乳/乳头玩弄",
                        "勇者/女神/魅魔/龙娘"
                    ],
                    "restrict": 0,
                    "xRestrict": 1,
                    "isOriginal": true,
                    "textLength": 15318,
                    "characterCount": 15318,
                    "wordCount": 7636,
                    "useWordCount": false,
                    "readingTime": 2297,
                    "bookmarkCount": 940,
                    "url": "https://i.pximg.net/c/150x150_80/novel-cover-master/img/2025/10/03/22/06/02/sci14662748_7a3bf9bc931bb9226fd24c16d2b8dd03_master1200.jpg",
                    "uploadTimestamp": 1760108404,
                    "reuploadTimestamp": 1760119936,
                    "isBookmarkable": true,
                    "bookmarkData": null,
                    "aiType": 1
                },
                {
                    "id": "26119985",
                    "userId": "59261594",
                    "series": {
                        "id": 14662748,
                        "viewableType": 0,
                        "contentOrder": 2
                    },
                    "title": "第二章，女神归心，深夜告白之后的激情做爱，手交乳交，但是得到主人的宠爱后就迫不及待挑衅的母狗是不是太欠教训了一点，必须得狠狠用子宫内射惩罚这只母狗了，享受完女神的晨勃口交后，再开始新的冒险吧",
                    "commentHtml": "  首先是中秋快乐<br /><br />  因为是第一次写长系列，所以希望大家多多提意见呢，第二篇本来要写一到两个新角色的，但是后面决定还是先把朱桑若攻略了再写，先把锅里的吃干抹净。",
                    "tags": [
                        "R-18",
                        "中文/中国语/Chinese/中國語/中国語",
                        "纯爱/后宫/西幻",
                        "隐秘做爱/调教/榨精/精液孕肚/激烈性爱",
                        "子宫奸/子宫交/潮吹/连续高潮/种付位",
                        "口交/深喉/吞精/吸精/巨乳/乳头玩弄/乳交",
                        "发情/内射/白给/中出/受孕",
                        "大量精液/精液喷出/子宫内射/亲吻",
                        "征服/母狗/肉便器/巨乳肥臀",
                        "勇者/女神"
                    ],
                    "restrict": 0,
                    "xRestrict": 1,
                    "isOriginal": true,
                    "textLength": 13285,
                    "characterCount": 13285,
                    "wordCount": 6534,
                    "useWordCount": false,
                    "readingTime": 1992,
                    "bookmarkCount": 1070,
                    "url": "https://i.pximg.net/c/150x150_80/novel-cover-master/img/2025/10/03/22/06/02/sci14662748_7a3bf9bc931bb9226fd24c16d2b8dd03_master1200.jpg",
                    "uploadTimestamp": 1759762846,
                    "reuploadTimestamp": 1759856281,
                    "isBookmarkable": true,
                    "bookmarkData": null,
                    "aiType": 1
                },
                {
                    "id": "26091301",
                    "userId": "59261594",
                    "series": {
                        "id": 14662748,
                        "viewableType": 0,
                        "contentOrder": 1
                    },
                    "title": "第一章，获得了催眠APP系统后就和好感度爆表的青梅一起出门冒险，什么？第一场战斗就打光明女神！什么？一拳秒了！什么？女神变成我的肉便器了！",
                    "commentHtml": "这个本来要当作万粉福利写的，提前写了，那万粉福利只能重新想了<br /><br />番外一已发<strong><a href=\"https://www.pixiv.net/novel/show.php?id=26100939\">novel/26100939</a></strong><br /><br />封面来自网络，侵删",
                    "tags": [
                        "R-18",
                        "中文/中国语/Chinese/中國語/中国語",
                        "纯爱/后宫/西幻",
                        "隐秘做爱/调教/榨精/精液孕肚/激烈性爱",
                        "子宫奸/子宫交/潮吹/连续高潮/骑乘位/种付位",
                        "口交/深喉/巨乳/乳头玩弄/吞精",
                        "发情/内射/败给/中出/受孕",
                        "大量精液/精液喷出/白给",
                        "征服/母狗/肉便器/巨乳肥臀",
                        "勇者/女神"
                    ],
                    "restrict": 0,
                    "xRestrict": 1,
                    "isOriginal": true,
                    "textLength": 18981,
                    "characterCount": 18981,
                    "wordCount": 9018,
                    "useWordCount": false,
                    "readingTime": 2847,
                    "bookmarkCount": 3276,
                    "url": "https://i.pximg.net/c/150x150_80/novel-cover-master/img/2025/10/04/00/00/24/ci26091301_afa92cc294b4590c76fe4b368f5ab446_master1200.jpg",
                    "uploadTimestamp": 1759503624,
                    "reuploadTimestamp": 1760022020,
                    "isBookmarkable": true,
                    "bookmarkData": null,
                    "aiType": 1
                }
            ]
        }
`

func TestParseJsonArray(t *testing.T) {
	data := gjson.ParseBytes([]byte(ct)).Get("seriesContents")
	var tmp []DAO.FollowData
	all := data.Array()
	err := json.Unmarshal([]byte(data.Raw), &tmp)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(tmp, len(all))
	for _, value := range all {
		var tt DAO.FollowData
		json.Unmarshal([]byte(value.Raw), &tt)
		t.Log(tt.ID)
		t.Log(1)
		//t.Log(tmp.IllustContentType.Sexual)
	}
}
