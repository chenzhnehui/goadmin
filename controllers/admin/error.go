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
	"github.com/astaxie/beego/orm"
	"goadmin/models/menu"
	"goadmin/utils"
)

type ErrorController struct {
	beego.Controller
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (c *ErrorController) Error404() {
	c.Data["content"] = "page not found"
	c.TplName = "error/404.html"
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (c *ErrorController) Error401() {
	adminId := c.GetSession("adminId")
	url := c.Ctx.Input.URL()
	if utils.InArray(url, []string{"/admin", "/admin/", "/admin/index/index", "/admin/index", "/admin/index/", "/admin/index/index/"}) {
		url = "/admin/index/index"
		if adminId != nil && adminId != "" {
			menu_groups_list := menu.GetMenuGroup(url, adminId)
			for _, v := range menu_groups_list.([]orm.Params) {
				val := v
				if val["Son"] != nil && val["Son"] != "" {
					for _, v1 := range val["Son"].([]interface{}) {
						vals := v1.(orm.Params)
						c.Ctx.Redirect(302, vals["Url"].(string))
					}
				} else {
					if val["Url"] != nil && val["Url"] != "" {
						c.Ctx.Redirect(302, val["Url"].(string))
					}
				}
			}
		}
	}
	c.Data["content"] = "您无权限操作"
	c.TplName = "error/401.html"
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (c *ErrorController) Error501() {
	c.Data["content"] = "server error"
	c.TplName = "error/404.html"
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (c *ErrorController) ErrorDb() {
	c.Data["content"] = "database is now down"
	c.TplName = "error/404.html"
}
