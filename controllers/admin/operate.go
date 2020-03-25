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
	"goadmin/models/operate"
	"strconv"
)

type OperateController struct {
	BaseController
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *OperateController) Lists() {
	adminId := this.GetSession("adminId").(int)
	if this.Ctx.Input.IsPost() {
		id, _ := this.GetInt("id", 0)
		if operate.Reads(adminId, id) == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "操作成功"}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": "操作失败"}
		}
		this.ServeJSON()
		return
	} else {
		condArr := make(map[string]string)
		condArr["user_id"] = strconv.Itoa(adminId)
		list, _ := operate.Lists(condArr, 1, 10)
		this.Data["json"] = list
	}
	this.ServeJSON()
	return
}
