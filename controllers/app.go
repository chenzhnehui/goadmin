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
	"goadmin/utils"
	"strings"
)

type AppController struct {
	beego.Controller
	CONTROLLER_NAME string
	ACTION_NAME     string
	MODULE_NAME     string
	TplNames        string
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *AppController) Inits(args ...interface{}) {
	ctx := this.Ctx
	controller, action := this.GetControllerAndAction()
	this.CONTROLLER_NAME = strings.Replace(strings.ToLower(controller), "controller", "", 1)
	this.ACTION_NAME = strings.ToLower(action)

	urlArr := strings.Split(this.Ctx.Input.URL(), "/")
	this.MODULE_NAME = "home"
	if len(args) > 0 {
		this.MODULE_NAME = args[0].(string)
	} else {
		if len(urlArr) > 1 && urlArr[1] != "" {
			this.MODULE_NAME = strings.ToLower(urlArr[1])
		}
	}
	this.TplNames = this.MODULE_NAME + "/" + this.CONTROLLER_NAME + "/" + this.ACTION_NAME + ".html"
	ctx.Input.SetData("CONTROLLER_NAME", this.CONTROLLER_NAME)
	ctx.Input.SetData("ACTION_NAME", this.ACTION_NAME)
	ctx.Input.SetData("MODULE_NAME", this.MODULE_NAME)
	ctx.Input.SetData("URLPARAM", map[string]string{"MODULE_NAME": this.MODULE_NAME, "CONTROLLER_NAME": this.CONTROLLER_NAME, "ACTION_NAME": this.ACTION_NAME})
}

//返回url
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *AppController) Urls(url string, args ...interface{}) string {
	param := map[string]string{"MODULE_NAME": this.MODULE_NAME, "CONTROLLER_NAME": this.CONTROLLER_NAME, "ACTION_NAME": this.ACTION_NAME}
	if len(args) > 0 {
		return utils.Urls(param, url, args)
	} else {
		return utils.Urls(param, url)
	}
}
