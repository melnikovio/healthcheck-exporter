package healthcheck

import (
	"fmt"
	"github.com/healthcheck-exporter/cmd/exporter"
	"github.com/sacOO7/gowebsocket"
	log "github.com/sirupsen/logrus"
	"time"
)

type WsClient struct {
	connections map[string]*WsConnection
	prometheus *exporter.Exporter
}

type WsConnection struct {
	urls  map[string]*Url
}

type Url struct {
	url  string
	Time int64
}

func NewWsClient(prometheus *exporter.Exporter) *WsClient {
	connection := make(map[string]*WsConnection)
	wc := WsClient{
		connections: connection,
		prometheus: prometheus,
	}

	return &wc
}

func (ws *WsClient) getWsConnection(jobId string) *WsConnection {
	if ws.connections[jobId] == nil {
		connection := WsConnection{
			urls: make (map[string]*Url),
		}
		ws.connections[jobId] = &connection
	}

	return ws.connections[jobId]
}

func (ws *WsClient) getWsUrl(jobId string, urlAddress string) *Url {
	connection := ws.getWsConnection(jobId)
	if connection.urls[urlAddress] == nil {
		url := Url{
			url:  urlAddress,
			Time: time.Now().Unix(),
		}
		connection.urls[urlAddress] = &url
		ws.addUrl(&url, jobId)
	}

	return connection.urls[urlAddress]
}

func (wsClient *WsClient) addUrl(url *Url, jobId string) {
	log.Info(fmt.Sprintf("Registering url: %s", url.url))
	socket := gowebsocket.New(url.url)

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Fatal("Received connect error - ", err)
	}

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server");
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println(jobId + ": Received message - " + message)
		wsClient.prometheus.IncCounter(jobId)
		url.Time = time.Now().Unix()
	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Received ping - " + data)
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Received pong - " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		return
	}

	socket.Connect()

	//c, _, err := websocket.DefaultDialer.Dial(url.url, nil)
	//defer c.Close()
	//
	//if err != nil {
	//	log.Error(fmt.Sprintf("Error listening ws connection (%s): %s", url, err.Error()))
	//}
	//
	//if c != nil {
	//	go func() {
	//		//c.SetReadLimit(100)
	//		for {
	//			_, _, err := c.ReadMessage()
	//
	//			//_, _, err := c.NextReader()
	//			if err != nil {
	//				log.Error(fmt.Sprintf("Error reading ws message (%s): %s", url, err.Error()))
	//				break
	//			}
	//
	//			ws.prometheus.IncCounter(jobId)
	//
	//			log.Info(fmt.Sprintf("%s: Received ws connection message", jobId))
	//			//log.Info(fmt.Sprintf("Received ws connection message: %s", message))
	//
	//			url.Time = time.Now().Unix()
	//		}
	//	}()
	//}
}

func (ws *WsClient) TimeLastMessage(jobId string, url string) int64 {
	return ws.getWsUrl(jobId, url).Time
}

func (ws *WsClient) TimeDifferenceWithLastMessage(jobId string, url string) int64 {
	return time.Now().Unix() - ws.getWsUrl(jobId, url).Time
}
