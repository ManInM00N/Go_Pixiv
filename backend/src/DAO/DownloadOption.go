package DAO

import (
	"fmt"
	"time"
)

type Option struct {
	Mode       int64
	R18        bool
	Likelimit  int64
	ShowSingle bool
	Suffix     string
	MinDate    string
	Rank       string
	RankDate   string
	DiffAuthor bool
}

const (
	ByPid    = int64(0)
	ByAuthor = int64(1)
	ByRank   = int64(2)
)

var (
	Rankmode = []string{"daily", "weekly", "monthly", "male", "female", "rookie", "original", "daily_r18", "weekly_r18", "male_r18", "female_r18"}

	ContentType = []string{"illust", "ugoira", "manga"}
)

func (Op *Option) Msg() {
	println(Op.Mode, Op.R18, Op.Likelimit, Op.ShowSingle, Op.Suffix, Op.MinDate)
}

type option func(*Option)

func NewOption(op ...option) *Option {
	Op := &Option{
		Mode:       ByPid,
		R18:        false,
		Likelimit:  0,
		ShowSingle: false,
		Suffix:     "",
		MinDate:    "19900101",
		Rank:       "daily",
		RankDate:   fmt.Sprintf("%04d%02d%02d", time.Now().Year(), time.Now().Month(), time.Now().Day()),
		DiffAuthor: true,
	}
	for _, O := range op {
		O(Op)
	}
	return Op
}
func WithR18(r18 bool) option {
	return func(o *Option) {
		o.R18 = r18
	}
}
func WithLikeLimit(num int64) option {
	return func(o *Option) {
		o.Likelimit = num
	}
}
func WithShowSingle(show bool) option {
	return func(o *Option) {
		o.ShowSingle = show
	}
}
func WithMode(mode int64) option {
	return func(o *Option) {
		o.Mode = mode
		//	int64(1) <<
	}
}

// Rankmode: Daily 0, Weekly 1, Monthly 2, Male 3, Female 4, Rookie 5
// Original 6, Daily_r18 7, Weekly_r18 8, male_r18 9, Female_r18 10
func WithRankmode(period int64) option {
	return func(o *Option) {
		o.Suffix += "&mode=" + Rankmode[period]
		o.Rank = Rankmode[period]
	}
}

// illust 0 ugoira 1 manga 2
func WithType(idx int64) option {
	return func(o *Option) {
		o.Suffix += "&content=" + ContentType[idx]
	}
}
func WithDate(date string) option {
	return func(o *Option) {
		o.Suffix += "&date=" + date
		o.RankDate = date
	}
}
func WithPage(page string) option {
	return func(o *Option) {
		o.Suffix += "&p=" + page
	}
}
func WithDiffAuthor(is bool) option {
	return func(o *Option) {
		o.DiffAuthor = is
	}
}

//func WithMinDate
