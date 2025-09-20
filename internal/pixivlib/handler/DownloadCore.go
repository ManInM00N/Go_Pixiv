package handler

import (
	_ "image/png"
)

var (
	Logo             []byte
	NowTaskMsg       = ""
	QueueTaskMsg     = ""
	ProcessMax       = int64(0)
	ProcessNow       = int64(0)
	RankLoadingNow   = false
	FollowLoadingNow = false
)
