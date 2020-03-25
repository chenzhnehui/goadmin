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
	"github.com/astaxie/beego"
	"goadmin/utils"
)

type PublicController struct {
	beego.Controller
}

/**
设置模板样式
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func (this *PublicController) Setskin() {
	if this.Ctx.Input.IsPost() {
		ace_skin := this.GetString("ace_skin")
		if ace_skin != "" {
			this.Ctx.SetCookie("ace_skin", ace_skin, 2592000)
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "设置成功"}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": "样式不存在"}
		}
		this.ServeJSON()
		return
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *PublicController) Test() {

	body := utils.Http("https://upload-images.jianshu.io/upload_images/5637367-6f1b8d1f04b4c02a.png", nil)

	this.Ctx.WriteString(body)
}
