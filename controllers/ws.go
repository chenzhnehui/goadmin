package controllers

// +----------------------------------------------------------------------
// | GOadmin [ I CAN DO IT JUST IT ]
// +----------------------------------------------------------------------
// | Copyright (c) 2020~2030 http://www.woaishare.cn All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( http://www.apache.org/licenses/LICENSE-2.0 )
// +----------------------------------------------------------------------
// | Author: chenzhenhui <971688607@qq.com>
// +----------------------------------------------------------------------
// | 分享交流QQ群请加  1062428023
// +----------------------------------------------------------------------

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"goadmin/utils"
	"log"
	"net/http"
)

type WsController struct {
	beego.Controller
}

/**
websocket连接
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func (this *WsController) Get() {
	ws, err := utils.Upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		log.Fatal("Not a websocket connection")
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Fatal("Cannot setup WebSocket connection:", err)
		return
	}
	utils.ConnectInit(ws)
	this.StopRun()
}
