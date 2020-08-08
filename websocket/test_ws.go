/**
 * @author: D-S
 * @date: 2020/7/25 5:03 下午
 */

package websocket

import (
	"fmt"
	"game-test/library/log"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Token struct {
	Token string `json:"token"`
}

type Message struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	SortIndex int64       `json:"sortIndex"`
}

var (
	once  = sync.Once{}
	once1 = sync.Once{}
	stop  = make(chan bool, 0)
)

type WS struct {
	URL      string
	Stop     chan bool
	Message  chan string
	Duration time.Duration
	Header   http.Header
	Proxy    string
}

func (e *WS) Start() {
	dialer := new(websocket.Dialer)
	if e.Proxy != "" {
		proxy, err := url.Parse(e.Proxy)
		if err != nil {
			log.Error(err.Error())
		}
		dialer.Proxy = http.ProxyURL(proxy)
	}
	conn, _, err := dialer.Dial(e.URL, e.Header)
	if err != nil {
		log.Error("ws error", zap.Error(err))
		return
	}
	go e.Heartbeat(conn)
	go e.ReadMessage(conn)
	for {
		select {
		case msg := <-e.Message:
			log.Info("==数据==")
			fmt.Println(msg)
		case stop := <-e.Stop:
			if stop {
				log.Error(err.Error())
				conn.Close()
				return
			}
		}
	}
}

func (e *WS) Heartbeat(conn *websocket.Conn) {
	for {
		err := conn.WriteMessage(websocket.TextMessage, []byte(`ping`))
		log.Info("send ping")
		if err != nil {
			stop <- true
			return
		}
		<-time.After(e.Duration)
	}
}
func (e *WS) ReadMessage(conn *websocket.Conn) {
	for {
		messType, message, err := conn.ReadMessage()
		if err != nil {
			stop <- true
			return
		}
		_ = messType
		e.Message <- string(message)
	}
}
