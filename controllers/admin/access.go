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
	"github.com/astaxie/beego/utils/pagination"
	"goadmin/models/access"
	"goadmin/models/category"
	"goadmin/models/menu"
	"strings"
)

type AccessController struct {
	BaseController
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *AccessController) Lists() {
	condArr := make(map[string]string)
	condArr["title"] = this.GetString("title")
	page, _ := this.GetInt("p")
	offset, err := beego.AppConfig.Int("pageoffset")
	if err != nil {
		offset = 15
	}

	list, count := access.Lists(condArr, page, offset)
	this.Data["paginator"] = pagination.SetPaginator(this.Ctx, offset, count)
	this.Data["list"] = list
	this.Data["count"] = count
	this.Data["params"] = condArr
	this.TplName = this.TplNames
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *AccessController) Add() {
	if this.Ctx.Input.IsPost() {
		status, _ := this.GetInt("status", 0)
		_, err := access.Add(this.GetString("title"), status, "0", "0")
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "添加成功", "url": this.Urls("lists")}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		this.TplName = this.TplNames
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *AccessController) Edit() {
	if this.Ctx.Input.IsPost() {
		id, _ := this.GetInt("id", 0)
		status, _ := this.GetInt("status", 0)
		err := access.Edit(id, this.GetString("title"), status, "", "")
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "编辑成功", "url": this.Urls("lists")}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		id, _ := this.GetInt("id")
		this.Data["info"] = access.GetInfo(id)
		this.TplName = this.TplNames
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *AccessController) Del() {
	ids := strings.Split(this.GetString("id"), ",")
	if this.GetString("id") == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "msg": "您没有选择任何数据"}
	} else {
		err := access.Del(ids)
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "删除成功"}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
	}
	this.ServeJSON()
	return
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *AccessController) Privilege() {
	if this.Ctx.Input.IsPost() {
		id, _ := this.GetInt("id", 0)
		menu_ids := strings.Join(this.GetStrings("auth"), ",")
		if menu_ids == "" {
			menu_ids = "0"
		}
		category_ids := strings.Join(this.GetStrings("category_ids"), ",")
		if category_ids == "" {
			category_ids = "0"
		}
		err := access.Edit(id, "", -1, menu_ids, category_ids)
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "分配成功，重新登录生效", "url": this.Urls("lists")}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		id, _ := this.GetInt("id")
		info := access.GetInfo(id)
		this.Data["info"] = info
		this.Data["privileges"] = menu.GetPrivilegeMenu(info.MenuIds)
		this.Data["categoryprivileges"] = category.GetPrivilegeCategory(info.CategoryIds)
		this.TplName = this.TplNames
	}
}
