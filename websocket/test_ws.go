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

const (
	//wsURL      = "wss://api.gamescorekeeper.com/v1/liveapi/63201" //15.30
	wsURL    = "ws://127.0.0.1:50052/api/v1/ws/" //15.30
	wsURLDEV = "ws://47.114.175.98:1325/ws"      //15.30
	//wsURLRelease = "ws://47.89.37.38:30023/api/v1/ws/" //15.30
	wsURLRelease   = "ws://47.57.152.73:30023/api/v1/ws/"             //15.30
	wsURLReleaseV2 = "wss://stream.dawnbyte.com/ws"                   //15.30
	wsURLDevV2     = "ws://47.114.175.98:1325/ws"                     //15.30
	pdev           = "ws://47.114.175.98:8082/ws"                     //15.30
	wsURL1         = "wss://api.gamescorekeeper.com/v1/liveapi/63158" //19.00
	wsURL2         = "wss://api.gamescorekeeper.com/v1/liveapi/63170" //20.00
	wsURLLOL       = "wss://api.gamescorekeeper.com/v1/liveapi/59853"
	wsURLDota2     = "wss://api.gamescorekeeper.com/v1/liveapi/62925"
	tokenConst     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjaHVuemhpIiwiaXNzIjoiR2FtZVNjb3Jla2VlcGVyIiwianRpIjoyNDY2MjI1NTE0OTgxMDcwODkzLCJjdXN0b21lciI6dHJ1ZX0.ewFZvscy0sx5-s3BP295qtpBcZtV4q9TEHUeKwJxJsI"
)

//{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjaHVuemhpIiwiaXNzIjoiR2FtZVNjb3Jla2VlcGVyIiwianRpIjoyNDY2MjI1NTE0OTgxMDcwODkzLCJjdXN0b21lciI6dHJ1ZX0.ewFZvscy0sx5-s3BP295qtpBcZtV4q9TEHUeKwJxJsI"}
//{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjaHVuemhpIiwiaXNzIjoiR2FtZVNjb3Jla2VlcGVyIiwianRpIjoyNDY2MjI1NTE0OTgxMDcwODkzLCJjdXN0b21lciI6dHJ1ZX0.ewFZvscy0sx5-s3BP295qtpBcZtV4q9TEHUeKwJxJsI"}
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
