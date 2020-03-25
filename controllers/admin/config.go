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
	"goadmin/models/config"
	"goadmin/utils"
	"strings"
)

type ConfigController struct {
	BaseController
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *ConfigController) Lists() {
	condArr := make(map[string]string)
	condArr["title"] = this.GetString("title")
	condArr["pid"] = this.GetString("pid", "0")
	condArr["type"] = this.GetString("type", "0")
	page, _ := this.GetInt("p")
	offset, err := beego.AppConfig.Int("pageoffset")
	if err != nil {
		offset = 15
	}

	list, count := config.Lists(condArr, page, offset)
	this.Data["paginator"] = pagination.SetPaginator(this.Ctx, offset, count)
	groupslist := config.GetConfig("CONFIG_GROUP_LIST")
	for k, v := range list {
		for k1, v1 := range groupslist.(map[string]interface{}) {
			if utils.Equal(k1, v["Group"]) {
				list[k]["GroupName"] = v1
				break
			}
		}
	}
	this.Data["list"] = list
	this.Data["count"] = count
	this.Data["params"] = condArr
	this.TplName = this.TplNames
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *ConfigController) Add() {
	if this.Ctx.Input.IsPost() {
		need, _ := this.GetInt("need", 0)
		status, _ := this.GetInt("status", 0)
		types, _ := this.GetInt("type", 0)
		group, _ := this.GetInt("group", 0)
		sort, _ := this.GetInt("sort", 0)
		extra := strings.Trim(this.GetString("extra", ""), " ")
		if utils.InArray(types, []int{4, 5, 6}) && extra == "" {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": "单选多选枚举类型请填写多选值"}
			this.ServeJSON()
			return
		}
		_, err := config.Add(this.GetString("name"), this.GetString("title"), need, status, types, group, sort, extra, this.GetString("remark"), this.GetString("value"))
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
func (this *ConfigController) Edit() {
	if this.Ctx.Input.IsPost() {
		id, _ := this.GetInt("id", 0)
		need, _ := this.GetInt("need", 0)
		status, _ := this.GetInt("status", 0)
		types, _ := this.GetInt("type", 0)
		group, _ := this.GetInt("group", 0)
		sort, _ := this.GetInt("sort", 0)
		extra := strings.Trim(this.GetString("extra", ""), " ")
		if utils.InArray(types, []int{4, 5, 6}) && extra == "" {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": "单选多选枚举类型请填写多选值"}
			this.ServeJSON()
			return
		}
		err := config.Edit(id, this.GetString("name"), this.GetString("title"), need, status, types, group, sort, extra, this.GetString("remark"), this.GetString("value"))
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "编辑成功", "url": this.Urls("lists")}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		id, _ := this.GetInt("id")
		info := config.GetInfo(id)
		this.Data["info"] = info
		this.TplName = this.TplNames
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *ConfigController) Del() {
	ids := strings.Split(this.GetString("id"), ",")
	if this.GetString("id") == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "msg": "您没有选择任何数据"}
	} else {
		err := config.Del(ids)
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
func (this *ConfigController) Groups() {
	if this.Ctx.Input.IsPost() {
		params := utils.Inputs(this.Ctx.Request.PostForm)
		for k, v := range params {
			go config.UpdateConfigVal(k, v)
		}
		utils.DelCache("Config.GetConfig.aconfigListCacheInfo")
		this.Data["json"] = map[string]interface{}{"code": 1, "msg": "操作成功"}
		this.ServeJSON()
		return
	} else {
		list, _ := config.Lists(map[string]string{"status": "1", "group": "-1", "order": "sort"}, 1, 1000)
		this.Data["list"] = list
		this.TplName = this.TplNames
	}
}
