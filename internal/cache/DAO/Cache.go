package DAO

import (
	"time"

	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

var (
	C  *cache.Cache
	Db *gorm.DB
)

type Cache struct {
	DownloadID string `gorm:"type:string;primaryKey" json:"id"`
	Type       string `gorm:"type:string" json:"type"`
	CreatedAt  time.Time
}

func GetDb() *gorm.DB {
	return Db
}
