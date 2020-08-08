/**
 * @author: D-S
 * @date: 2020/8/7 5:47 下午
 */

package websocket

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type binder struct {
	Mux    sync.Mutex
	ConMap map[string]*websocket.Conn
}

var (
	sBinder = binder{
		Mux:    sync.Mutex{},
		ConMap: make(map[string]*websocket.Conn),
	}
	sOnce     = sync.Once{}
	sUpGrader = websocket.Upgrader{
		//跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	pushChannel = make(chan []byte, 10000)
)

func Push(data []byte) {
	pushChannel <- data
}

type WSService struct {
	w  http.ResponseWriter
	r  *http.Request
	Ip string // 连接Ip
}

func (wss *WSService) Start() error {

	con, err := sUpGrader.Upgrade(wss.w, wss.r, nil)
	if err != nil {
		return err
	}
	//保存连接
	if err := wss.saveCon(con); err != nil {
		return err
	}
	// 读取心跳
	wss.readMessage(con)
	return nil
}

func (wss *WSService) readMessage(con *websocket.Conn) {
	for {
		_, message, err := con.ReadMessage()
		if err != nil {
			fmt.Println(fmt.Sprintf("ip %v read message err %v ", wss.Ip, err))
			return
		}
		fmt.Println(fmt.Sprintf("ip %v message %v ", wss.Ip, message))
		if err := con.WriteMessage(websocket.TextMessage, []byte(`pong`)); err != nil {
			fmt.Println(fmt.Sprintf("ip %v send pong err %v ", wss.Ip, err))
		}
	}
}

func (wss *WSService) saveCon(con *websocket.Conn) error {
	defer sBinder.Mux.Unlock()
	sBinder.Mux.Lock()
	if wss.Ip == "" {
		return errors.New("ip is nil")
	}
	c, ok := sBinder.ConMap[wss.Ip]
	if ok {
		c.Close()
	}
	sBinder.ConMap[wss.Ip] = con
	return nil
}

func toSend() {
	for {
		select {
		case data, ok := <-pushChannel:
			if ok {
				if len(sBinder.ConMap) == 0 {
					continue
				}
				sendMessage(data)
			}
		}
	}
}

func sendMessage(message []byte) {
	defer sBinder.Mux.Unlock()
	sBinder.Mux.Lock()
	for ip, con := range sBinder.ConMap {
		if err := con.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println(fmt.Sprintf("ip %v seng message err %v ", ip, err))
			continue
		}
		fmt.Println(fmt.Sprintf("ip %v seng message success %v ", ip, string(message)))
	}
}
