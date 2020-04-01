/**
* @Author: D-S
* @Date: 2020/3/31 11:37 下午
 */

package router

import (
	"game-test/websocket"
	"github.com/gin-gonic/gin"
)

func LoadRouter(engine *gin.Engine) {
	//权限校验
	//bAuth := engine.Group("/api/v1/", auth.AuthMiddleware())
	bAuth := engine.Group("/api/v1/")
	bAuth.GET("ws", websocket.NewConnect) //游戏列表
	//盘口结算推送  赔率推送
}
