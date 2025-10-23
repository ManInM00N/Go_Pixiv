package test

import (
	. "main/init"
	_ "main/init"
	"main/internal/pixivlib/DAO"
	. "main/internal/pixivlib/handler"
	"os"
	"testing"
)

func TestDownloadUgoira(t *testing.T) {
	text := "104107248"

	op := DAO.NewOption(DAO.WithMode(DAO.ByPid), DAO.WithLikeLimit(0), DAO.WithR18(true), DAO.WithShowSingle(true), DAO.WithDiffAuthor(false))
	_, success := JustDownload(text, op, func(name string, data ...interface{}) {

	})
	if !success {
		t.Errorf("Download ugoira failed")
	}
}

func TestDownloadNovel(t *testing.T) {
	id := "23608144"
	success := DownloadNovel(id)
	if !success {
		t.Errorf("Download ugoira failed")
	}
}

func TestMain(m *testing.M) {
	CacheInit()
	code := m.Run()

	// 退出测试
	os.Exit(code)
}
