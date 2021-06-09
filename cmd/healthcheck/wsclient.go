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
	connection := make([]wsConnection, 0)
	wc := WsClient{
		connection: connection,
	}

	return &wc
}

func (ws *WsClient) addUrl(jobId string, url string) {
	////
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Error(fmt.Sprintf("Error listening ws connection (%s): %s", url, err.Error()))
	}

	if c != nil {
		go func() {
			ws.connection = append(ws.connection, wsConnection{
				url:  url,
				time: time.Now().Unix(),
			})
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Info(fmt.Sprintf("Error reading ws connection message: %s", err.Error()))
					ws.addUrl(jobId, url)
					break
				}
				log.Info(fmt.Sprintf("%s: Received ws connection message", jobId))
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

func (ws *WsClient) LastMessageTime(jobId string, url string) int64 {
	for i := 0; i < len(ws.connection); i++ {
		if ws.connection[i].url == url {
			return ws.connection[i].time
		}
	}

	ws.addUrl(jobId, url)

	return 0
}

func (ws *WsClient) DifferenceLastMessageTime(jobId string, url string) int64 {
	for i := 0; i < len(ws.connection); i++ {
		if ws.connection[i].url == url {
			return time.Now().Unix() - ws.connection[i].time
		}
	}

	ws.addUrl(jobId, url)

	return 0
}
