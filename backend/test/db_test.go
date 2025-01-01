package test

import (
	"fmt"
	"main/backend/src/DAO"
	. "main/backend/src/init"
	"testing"
)

func TestDb(t *testing.T) {
	CacheInit()
	db := DAO.GetDb()
	var cnt int64
	db.Model(&DAO.Cache{}).Where("download_id = ?", "125675405").Count(&cnt)
	fmt.Println(cnt)
}
