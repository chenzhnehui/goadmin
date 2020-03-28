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
	"encoding/json"
	"github.com/astaxie/beego"
	"goadmin/controllers"
	"goadmin/models/category"
	"goadmin/models/menu"
	"goadmin/models/operate"
	"goadmin/utils"
)

type BaseController struct {
	controllers.AppController
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *BaseController) Prepare() {
	this.Inits()
	ctx := this.Ctx
	ctx.Input.SetData("WEB_NAME", beego.AppConfig.String("webname"))
	urls := this.Urls("")
	adminId := this.GetSession("adminId")
	adminInfo := this.GetSession("adminInfo")
	if (adminId == nil || adminId == "") && urls != this.Urls("user/login") {
		this.Redirect(this.Urls("user/login"), 302)
	}

	if ctx.GetCookie("ace_skin") == "" { //设置网页样式模板
		ctx.SetCookie("ace_skin", "no-skin", 2592000)
		ctx.Input.SetData("ACE_SIKN", "no-skin")
	} else {
		ctx.Input.SetData("ACE_SIKN", ctx.GetCookie("ace_skin"))
	}

	//建立和文章模型一样的分类控制方式和权限,把控制器放在对应 controllerList数组中，复制 article控制器和列表即可

	controllerList := map[string]string{"article": "1"} //如果加商品控制器，可使用 map[string]string{"goods","types"} ，types是绑定的对应分类属性，然后在菜单里面添加对应权限，路由，控制器和模板views

	//跳过不检查权限的控制其方法
	ignoreUrl := []string{this.Urls("user/login"), this.Urls("user/logout")}
	if utils.InArray(urls, ignoreUrl) {
		return
	}
	if adminId != nil { //登录成功
		category_id, _ := this.GetInt("category_id", 0)
		if this.Ctx.Input.IsGet() {
			menu_groups_list := menu.GetMenuGroup(urls, adminId)
			for k, v := range controllerList {
				menu_groups_list = category.GetArticleMenuGroup(menu_groups_list, adminId, category_id, this.ACTION_NAME, this.MODULE_NAME+"/"+k, v)
			}
			ctx.Input.SetData("menu_groups_list", menu_groups_list) //过滤权限
			ctx.Input.SetData("adminInfo", adminInfo)
			if !menu.CheckPrivilege(urls, "", adminId.(int)) {
				ctx.Abort(401, "401")
			}
			for k, _ := range controllerList {
				if utils.Equal(this.CONTROLLER_NAME, k) {
					if !category.IsArticlePrivilege(adminId, category_id) {
						ctx.Abort(401, "401")
					}
				}
			}
			ctx.Input.SetData("RequestFormParam", utils.Inputs(this.Ctx.Request.Form))
		} else {
			operateParamJson, _ := json.Marshal(&ctx.Request.Form)
			go operate.Add(adminId.(int), 0, urls, string(operateParamJson))
			if !menu.CheckPrivilege(urls, "", adminId.(int)) {
				ctx.Output.JSON(map[string]interface{}{"code": 0, "msg": "暂无权限"}, false, false)
				return
			}
			for k, _ := range controllerList {
				if utils.Equal(this.CONTROLLER_NAME, k) {
					if !category.IsArticlePrivilege(adminId, category_id) {
						ctx.Output.JSON(map[string]interface{}{"code": 0, "msg": "暂无权限"}, false, false)
						return
					}
				}
			}
			//ctx.Output.JSON(map[string]interface{}{"code": 0, "msg": "该平台为共享案例平台，请勿操作数据，谢谢！"}, false, false)
			//return
		}
		//文章分类权限
	}

}
