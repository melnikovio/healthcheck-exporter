package healthcheck

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
)

type WsClient struct {
	connection []wsConnection
}

type wsConnection struct {
	url  string
	time int64
}

func NewWsClient() *WsClient {
	wc := WsClient{}

	return &wc
}

func (ws *WsClient) addUrl(url string) {
	////
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Error(fmt.Sprintf("Error listening ws connection: %s", err.Error()))
	}

	if c != nil {
		go func() {
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Info(fmt.Sprintf("Error reading ws connection message: %s", err.Error()))
				}
				log.Info(fmt.Sprintf("Received ws connection message"))
				log.Trace(fmt.Sprintf("Received ws connection message: %s", message))

				for i := 0; i < len(ws.connection); i++ {
					if ws.connection[i].url == url {
						ws.connection[i].time = int64(time.Now().Unix())
					}
				}
			}
		}()
	}
	////
}

func (ws *WsClient) LastMessageTime(url string) int64 {
	for i := 0; i < len(ws.connection); i++ {
		if ws.connection[i].url == url {
			return ws.connection[i].time
		}
	}

	ws.addUrl(url)

	return 0
}

func (ws *WsClient) DifferenceLastMessageTime(url string) int64 {
	for i := 0; i < len(ws.connection); i++ {
		if ws.connection[i].url == url {
			return int64(time.Now().Unix()) - ws.connection[i].time
		}
	}

	ws.addUrl(url)

	return 0
}
