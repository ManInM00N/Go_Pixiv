package init

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func UpdateProgress(c *gin.Context) {
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
			log.Println(err)
			break
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
			log.Println(err)
			break
		}
	}
}
