package home

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
	"goadmin/utils"
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
		if client_id != "" {
			utils.WsJoinGroup(client_id, "system-home-index")
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "加群完成"}
			this.ServeJSON()
			return
		}
	} else {
		this.TplName = this.TplNames
	}
}
