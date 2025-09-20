package handler

import (
	"github.com/ManInM00N/go-tool/statics"
	. "main/internal/pixivlib/DAO"
	"main/internal/taskQueue"
)

func Download_By_Pid(text string, callEvent func(name string, data ...interface{})) {
	text = statics.CatchNumber(text)
	println(text)
	if text == "" {
		return
	}
	taskQueue.SinglePool.AddTask(func() (interface{}, error) {
		op := NewOption(WithMode(ByPid), WithLikeLimit(0), WithR18(true), WithShowSingle(true), WithDiffAuthor(false))
		JustDownload(text, op, callEvent)
		return nil, nil
	})
}
