/**
* @Author: D-S
* @Date: 2020/3/31 11:33 下午
 */

package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"game-test/library/log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type WebSocketCon struct {
	Con         *websocket.Conn `json:"con"`
	Token       string          `json:"token"`
	ClientId    string          `json:"clientId"`
	ConnectTime time.Time       `json:"connectTime"`
	LastTime    time.Time       `json:"lastTime"`
}

type Binder struct {
	Mux           sync.Mutex              `json:"mux"`
	WebSocketCons map[string]WebSocketCon `json:"webSocketCons"`
}

var Bindder = &Binder{
	Mux:           sync.Mutex{},
	WebSocketCons: make(map[string]WebSocketCon, 0),
}

var upGrader = websocket.Upgrader{
	//跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

func NewConnect(c *gin.Context) {
	id := c.GetHeader("tenant_id")
	if id == "" {
		log.Error("参数错误")
		return
	}
	clientId := "sdfsdfsd"
	con, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error("ws 连接错误 ")
		return
	}
	defer con.Close()
	wsc := WebSocketCon{
		Con:         con,
		ConnectTime: time.Now(),
	}

	if e, ok := Bindder.WebSocketCons[clientId]; ok {
		log.Info("ws在其他端登录，本次连接即将断开...")
		wsc.LastTime = e.ConnectTime
		e.Con.Close()
		return
	}
	Bindder.Mux.Lock()
	Bindder.WebSocketCons[clientId] = wsc
	Bindder.Mux.Unlock()

	log.Info("ws新连接")
	ProcessWsMsg(wsc, clientId)

}

func ProcessWsMsg(wsc WebSocketCon, clientId string) {
	SendWsMsg(clientId, "测试成功")
	for {
		msgType, data, err := wsc.Con.ReadMessage()
		if err != nil {
			log.Error("ws接收消息异常 即将关闭连接...")
			wsc.Con.Close()
			Bindder.Mux.Lock()
			delete(Bindder.WebSocketCons, clientId)
			Bindder.Mux.Unlock()
			return
		}

		switch msgType {
		case websocket.PingMessage:
			log.Info("ping")
			if err := wsc.Con.WriteMessage(websocket.PongMessage, nil); err != nil {
				log.Error("回复客户端心跳失败")
			}
		default:
			log.Info("ws收到新消息")
			log.Info("", zap.Any("data", data))
		}
	}
}

func SendWsMsg(clientId, message string) error {
	wsc, ok := Bindder.WebSocketCons[clientId]
	if !ok {
		log.Error("发送ws错误 ")
		return errors.New("无该用户连接")
	}

	if err := wsc.Con.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Error("发送ws错误 ")
		return err
	}
	log.Info("发送ws成功")
	return nil
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
