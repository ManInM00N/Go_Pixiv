package init

import (
	"github.com/gin-gonic/gin"
	"main/backend/src/DAO"
)

func UpdateSetting(c *gin.Context) {
	var set DAO.Settings
	c.BindJSON(&set)
	DebugLog.Println(set.MsgDetail())
	Setting.UpdateSettings(set)
	UpdateSettings()
}
func GetSetting(c *gin.Context) {
	c.JSON(200, gin.H{
		"setting": Setting,
	})
}
