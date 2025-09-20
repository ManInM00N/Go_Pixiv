package test

import (
	_ "main/init"
	. "main/internal/pixivlib/handler"
	"testing"
)

func TestDownloadUgoira(t *testing.T) {
	return
	// text := "124969828"
	// success := DownloadNovel(text)
	// if !success {
	// 	t.Errorf("Download ugoira failed")
	// }
}

func TestDownloadNovel(t *testing.T) {
	id := "23608144"
	success := DownloadNovel(id)
	if !success {
		t.Errorf("Download ugoira failed")
	}
}
