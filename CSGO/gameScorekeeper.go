/**
* @Author: D-S
* @Date: 2020/3/20 10:33 上午
 */

package CSGO

import (
	"encoding/json"
	"fmt"
	"game-test/library/log"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

//爬虫部分
var (
	//分页参数
	pageLimit = 50

	//代理
	requestCount = 0
	proxyList    = []string{
		"",
		"http://47.91.246.62:59073",
		"http://47.56.193.197:59073",
	}

	httpHeader = http.Header{
		"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36"},
	}
)

const (
	//wsURL      = "wss://api.gamescorekeeper.com/v1/liveapi/63201" //15.30
	wsURL      = "ws://127.0.0.1:50051/api/v1/ws"                 //15.30
	wsURL1     = "wss://api.gamescorekeeper.com/v1/liveapi/63158" //19.00
	wsURL2     = "wss://api.gamescorekeeper.com/v1/liveapi/63170" //20.00
	wsURLLOL   = "wss://api.gamescorekeeper.com/v1/liveapi/59853"
	wsURLDota2 = "wss://api.gamescorekeeper.com/v1/liveapi/62925"
	tokenConst = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjaHVuemhpIiwiaXNzIjoiR2FtZVNjb3Jla2VlcGVyIiwianRpIjoyNDY2MjI1NTE0OTgxMDcwODkzLCJjdXN0b21lciI6dHJ1ZX0.ewFZvscy0sx5-s3BP295qtpBcZtV4q9TEHUeKwJxJsI"
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

var once = sync.Once{}

func WS() {
	var dialer *websocket.Dialer
	log.Info("wss url", zap.String("websocket", wsURL))
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		log.Error("ws error", zap.Error(err))
		return
	}
	//var stop  =  make(chan bool)
	for {
		messType, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("ws error", zap.Error(err))
			return
		}
		conn.WriteMessage(websocket.PingMessage, nil)
		_ = messType
		//result := Message{}
		//if err := json.Unmarshal(message, &result); err != nil {
		//	log.Error("marshal err", zap.Error(err))
		//	return
		//}
		log.Info("===数据===")
		fmt.Println(string(message))
		//if result.Type == "auth" {
		//	log.Info("auth 权限验证")
		//	token := new(Token)
		//	token.Token = tokenConst
		//	b, _ := json.Marshal(token)
		//	if err := conn.WriteMessage(messType, b); err != nil {
		//		log.Error("ws error", zap.Error(err))
		//		return
		//	}
		//}
		once.Do(func() {
			go test(conn)
		})
	}
}

func test(conn *websocket.Conn) {
	for {
		b, err := json.Marshal("ping")
		if err != nil {
			log.Error(err.Error())
		}
		log.Info("send ping")
		err1 := conn.WriteMessage(websocket.PingMessage, b)
		if err1 != nil {
			log.Error("ping error ", zap.Any("err", zap.Error(err)))
			//stop <- true
			return
		}
		<-time.After(30 * time.Second)
	}
}
func Ray57() {
	log.Info(fmt.Sprintf("%s开始连接websocket", "ray57"))

	u := "wss://cfsocket.raybet.ai/socketcluster/"
	conn, _, err := websocket.DefaultDialer.Dial(u, httpHeader)
	if err != nil {
		log.Error("Ray57 websocket Dial err ", zap.Error(err))
		return
	}

	//发订阅消息
	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"event":"#handshake","data":{"authToken":null},"cid":1}`))
	if err != nil {
		log.Error("Ray57 websocket handshake err ", zap.Error(err))
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"event":"#subscribe","data":{"channel":"match"},"cid":2}`))
	if err != nil {
		log.Error("Ray57 websocket handshake err ", zap.Error(err))
		return
	}

	//处理接收
	for {
		select {
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Error("Ray57 websocket ReadMessage err ", zap.Error(err))
				return
			}

			recv := string(message)

			//处理心跳消息
			if recv == "#1" {
				_ = conn.WriteMessage(websocket.TextMessage, []byte(`#2`))
				continue
			}
			result := new(Ray57WS)
			err = json.Unmarshal(message, result)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			//log.Info(fmt.Sprintf("%s 收到websocket消息:%s", "ray57", recv))
			log.Info("json", zap.Any("json", result))
		}
	}

}

type Ray57WS struct {
	Event string `json:"event"`
	Data  struct {
		Channel string `json:"channel"`
		Data    struct {
			Source string `json:"source"`
			Rates  []struct {
				Id         int32  `json:"id"`
				MatchId    int64  `json:"match_id"`
				Rate       string `json:"odds"`
				LastUpdate string `json:"last_update"`
				Status     int64  `json:"status"`
			} `json:"odds"`
		} `json:"data"`
	} `json:"data"`
}
