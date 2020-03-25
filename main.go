package main

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
	"github.com/astaxie/beego/logs"
	"goadmin/controllers/admin"
	"goadmin/inits/template"
	_ "goadmin/models"
	_ "goadmin/routers"
	"goadmin/utils"
)

func main() {

	utils.InitCache()
	template.InitTemplateFunc()
	beego.ErrorController(&admin.ErrorController{})
	logs.SetLogger(logs.AdapterFile, `{"filename":"`+beego.AppConfig.String("logfile")+`"}`)
	logs.EnableFuncCallDepth(true)
	logs.Async()
	beego.Run()
}
