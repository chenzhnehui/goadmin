package template

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
	"goadmin/models/attribute"
	"goadmin/models/config"
	"goadmin/models/menu"
	"goadmin/utils"
)

func InitTemplateFunc() {
	beego.AddFuncMap("CheckPrivileges", CheckPrivileges)
	beego.AddFuncMap("GetConfigType", utils.GetConfigType)
	beego.AddFuncMap("GetConfigTypeField", utils.GetConfigTypeField)
	beego.AddFuncMap("GetConfig", config.GetConfig)
	beego.AddFuncMap("Equal", utils.Equal)
	beego.AddFuncMap("InString", utils.InString)
	beego.AddFuncMap("JsonEncode", utils.JsonEncode)
	beego.AddFuncMap("JsonDecode", utils.JsonDecode)
	beego.AddFuncMap("StrToArray", utils.StrToArray)
	beego.AddFuncMap("AnalysisStr", utils.AnalysisStr)
	beego.AddFuncMap("Operation", utils.Operation)
	beego.AddFuncMap("GetFeildSort", attribute.GetFeildSort)
	beego.AddFuncMap("GetAttributeByName", attribute.GetAttributeByName)
	beego.AddFuncMap("GetFields", utils.GetFields)
	beego.AddFuncMap("Urls", utils.Urls)

}

func CheckPrivileges(url string, admin_id interface{}) bool {
	if admin_id == nil {
		return false
	}
	return menu.CheckPrivilege(url, "", admin_id.(int))
}
