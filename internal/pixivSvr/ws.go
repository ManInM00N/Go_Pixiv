package pixivSvr

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"main/pkg/utils"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	WS *websocket.Conn
)

type DownloadRequest struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Period string `json:"period"`
	Time   string `json:"time"`
}

func (a *DownloadRequest) Msg() string {
	return fmt.Sprintf(`Type: %s, ID: %s, Period: %s, Time: %s`, a.Type, a.ID, a.Period, a.Time)
}
func UpdateProgress(c *gin.Context) {
	var err error
	WS, err = upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.DebugLog.Println(err)
		return
	}
	defer WS.Close()

	for {
		_, msg, err := WS.ReadMessage()
		if err != nil {
			utils.DebugLog.Println("Read Error", err)
			break
		}
		var rq DownloadRequest
		err = json.Unmarshal(msg, &rq)
		if err != nil {
			utils.DebugLog.Println("Write back", err)
			break
		} else {
			utils.InfoLog.Println(rq.Msg())
		}
	}
}

func Transform(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		userInput := string(msg)

		err = ws.WriteMessage(websocket.TextMessage, []byte(userInput+" txt"))
		if err != nil {
			log.Println("Write back", err)
			break
		}
	}
}
