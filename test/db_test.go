package test

import (
	"fmt"
	. "main/init"
	"main/internal/cache/DAO"
	"testing"
)

func TestDb(t *testing.T) {
	CacheInit()
	db := DAO.GetDb()
	var cnt int64
	db.Model(&DAO.Cache{}).Where("download_id = ?", "125675405").Count(&cnt)
	fmt.Println(cnt)
}
