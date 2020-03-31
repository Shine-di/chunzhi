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
	"sync"
	"time"
)

const (
	wsURL = "wss://api.gamescorekeeper.com/v1/liveapi/63201"  //15.30
	wsURL1 = "wss://api.gamescorekeeper.com/v1/liveapi/63158" //19.00
	wsURL2 = "wss://api.gamescorekeeper.com/v1/liveapi/63170" //20.00
	wsURLLOL = "wss://api.gamescorekeeper.com/v1/liveapi/59853"
	wsURLDota2 = "wss://api.gamescorekeeper.com/v1/liveapi/62925"
	tokenConst = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjaHVuemhpIiwiaXNzIjoiR2FtZVNjb3Jla2VlcGVyIiwianRpIjoyNDY2MjI1NTE0OTgxMDcwODkzLCJjdXN0b21lciI6dHJ1ZX0.ewFZvscy0sx5-s3BP295qtpBcZtV4q9TEHUeKwJxJsI"
)

//{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjaHVuemhpIiwiaXNzIjoiR2FtZVNjb3Jla2VlcGVyIiwianRpIjoyNDY2MjI1NTE0OTgxMDcwODkzLCJjdXN0b21lciI6dHJ1ZX0.ewFZvscy0sx5-s3BP295qtpBcZtV4q9TEHUeKwJxJsI"}
//{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjaHVuemhpIiwiaXNzIjoiR2FtZVNjb3Jla2VlcGVyIiwianRpIjoyNDY2MjI1NTE0OTgxMDcwODkzLCJjdXN0b21lciI6dHJ1ZX0.ewFZvscy0sx5-s3BP295qtpBcZtV4q9TEHUeKwJxJsI"}
type Token struct {
	Token string 	`json:"token"`
}

type Message struct {
	Type string `json:"type"`
	Payload  interface{}  `json:"payload"`
	SortIndex  int64  `json:"sortIndex"`
}

var once  =  sync.Once{}
func WS() {
	var dialer *websocket.Dialer
	log.Info("wss url",zap.String("websocket",wsURL))
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
		result := Message{}
 		if err :=  json.Unmarshal(message,&result);err != nil{
 			log.Error("marshal err",zap.Error(err))
			return
		}
		log.Info("===数据===")
		fmt.Println(string(message))
		if result.Type == "auth" {
			log.Info("auth 权限验证")
			token := new(Token)
			token.Token = tokenConst
			b,_ := json.Marshal(token)
			if err := conn.WriteMessage(messType,b);err != nil{
				log.Error("ws error", zap.Error(err))
				return
			}
		}
		once.Do(func() {
			go test(conn)
		})
	}
}

func test(conn *websocket.Conn) {
	for  {
		b,err:= json.Marshal("ping")
		if err != nil {
			log.Error(err.Error())
		}
		log.Info("send ping")
		err1 := conn.WriteMessage(1,b)
		if err1 != nil {
			log.Error("ping error ", zap.Any("err", zap.Error(err)))
			//stop <- true
			return
		}
		<- time.After(30 * time.Second)
	}
}