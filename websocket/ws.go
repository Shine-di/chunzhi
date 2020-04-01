/**
* @Author: D-S
* @Date: 2020/3/31 11:33 下午
 */

package websocket

import (
	"errors"
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
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewConnect(c *gin.Context) {
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
