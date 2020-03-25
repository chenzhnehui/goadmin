package admin

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
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"goadmin/utils"
	"runtime"
	"strconv"
	"time"
)

type IndexController struct {
	BaseController
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *IndexController) Index() {
	if this.Ctx.Input.IsPost() {
		client_id := this.GetString("client_id")
		if(client_id != ""){
			utils.WsJoinGroup(client_id,"system-admin-index")
		}
		go func() {
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
				utils.WsSendToGroup(utils.Message{Type: "system", Message: utils.JsonEncode(jsons)},"system-admin-index")
			}
		}()
		this.Data["json"] = map[string]interface{}{"code": 1, "msg": "信息收集中"}
		this.ServeJSON()
		return
	} else {
		this.Data["system"] = map[string]interface{}{"os": runtime.GOOS, "version": runtime.Version(), "cpu": runtime.GOMAXPROCS(0), "ip": this.Ctx.Input.IP()}
		this.TplName = this.TplNames
	}
}
