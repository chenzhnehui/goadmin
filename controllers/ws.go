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
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"goadmin/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

type WsController struct {
	beego.Controller
}
/**
处理socket消息
Type消息分类有如下
ping：心跳发送ping消息 ，Message{Type: "ping", Message: Now()}
pong：响应的ping消息
level：用户退出消息   Message{Clientid: client_id, Type: "level", Message: '{uid:绑定的uid,group:[加入的群]}'}
join:用户加入消息  Message{Clientid: client_id, Type: "join", Message: "用户加入"}

joingroup:加入群消息 Message{Clientid: client_id, Type: "joingroup", Message: group_id}
levelgroup：退出群消息 Message{Clientid: client_id, Type: "levelgroup", Message: group_id}
ungroup：解散群消息 Message{ Type: "ungroup", Message: group_id}


其他消息可以自定义发送和接收,所有的消息只能这里处理，不要其他地方对Msg进行消费
*/
func init()  {
	go func() {
		for {
			select { // 消息通道中有消息则执行，否则堵塞	// 哪个case可以执行，则转入到该case。都不可执行，则堵塞。
			case msg := <- utils.Msg:
				switch msg.Type { //消息类型
				case "pong": //收到ping消息的回复
				case "join":utils.WsSendToClient(msg.Clientid, msg)
				case "level": //用户退出消息
				case "joingroup":
					if(msg.Message == "system-home-index"){
						utils.WsSendToGroup(utils.Message{Type:"groupcount",Message:strconv.Itoa(utils.WsGetClientIdCountByGroup("system-home-index"))},"system-home-index")
					}
				case "levelgroup":
					if(msg.Message == "system-home-index"){
						utils.WsSendToGroup(utils.Message{Type:"groupcount",Message:strconv.Itoa(utils.WsGetClientIdCountByGroup("system-home-index"))},"system-home-index")
					}
				case "ungroup"://解散群消息
				case "chat":
					utils.WsSendToGroup(msg,"system-home-index",msg.Clientid)
				case "system":
					utils.WsSendToGroup(msg,"system-admin-index")
				default:
				}
			}
		}
	}()
	go getSysInfo()
}


func getSysInfo()  {
	for {
		time.Sleep(1 * time.Second)
		mem, _ := mem.VirtualMemory()
		jsons := make(map[string]interface{})
		jsons["mem"] = utils.JsonDecode(utils.JsonEncode(mem))

		diskslist, _ := disk.Partitions(true) //所有分区
		var disktotal float64
		var diskused float64
		for _, json := range diskslist {
			diskInfo, _ := disk.Usage(utils.JsonDecode(utils.JsonEncode(json), "mountpoint").(string)) //指定某路径的硬盘使用情况
			disktotal += utils.JsonDecode(utils.JsonEncode(diskInfo), "total").(float64) / 1048576
			diskused += utils.JsonDecode(utils.JsonEncode(diskInfo), "used").(float64) / 1048576
		}
		diskusedPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", (diskused/disktotal)*100), 64)
		disktotalGB := int(disktotal / 1024)
		diskusedGB, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", diskused/1024), 64)
		jsons["disk"] = map[string]interface{}{"usedPercent": diskusedPercent, "disktotalGB": disktotalGB, "diskusedGB": diskusedGB}
		cpuusedPercent, _ := cpu.Percent(time.Duration(time.Second), false)
		//cpuInfo,_ :=  cpu.Times(false)
		jsons["cpu"] = map[string]interface{}{"usedPercent": cpuusedPercent}
		//jsons["io"],_ = net.IOCounters(true) //获取网络读写字节／包的个数
		utils.Msg <- utils.Message{Type: "system", Message: utils.JsonEncode(jsons)}
	}
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
