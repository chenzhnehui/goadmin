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
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils/pagination"
	"goadmin/models/groups"
	"goadmin/models/menu"
	"goadmin/utils"
	"strconv"
	"strings"
)

type MenuController struct {
	BaseController
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *MenuController) Lists() {
	condArr := make(map[string]string)
	condArr["title"] = this.GetString("title")
	condArr["pid"] = this.GetString("pid", "0")
	page, _ := this.GetInt("p")
	offset, err := beego.AppConfig.Int("pageoffset")
	if err != nil {
		offset = 15
	}

	list, count := menu.Lists(condArr, page, offset)
	this.Data["paginator"] = pagination.SetPaginator(this.Ctx, offset, count)
	if condArr["pid"] == "0" {
		grouplist, _ := groups.Lists(map[string]string{}, 1, 1000)
		for k, v := range list {
			for _, v1 := range grouplist {
				if utils.Equal(v["GroupsId"], v1["Id"]) {
					list[k]["GroupsName"] = v1["Title"]
					break
				}
			}
		}
	} else {
		pid, _ := strconv.Atoi(condArr["pid"])
		info := menu.GetInfo(pid)
		for k, _ := range list {
			list[k]["GroupsName"] = info.Title
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
func (this *MenuController) Add() {
	if this.Ctx.Input.IsPost() {
		pid, _ := this.GetInt("pid", 0)
		status, _ := this.GetInt("status", 0)
		sort, _ := this.GetInt("sort", 0)
		hide, _ := this.GetInt("hide", 0)
		groups_id, _ := this.GetInt("groups_id", 0)
		_, err := menu.Add(this.GetString("title"), pid, sort, this.GetString("url"), hide, groups_id, status)
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "添加成功", "url": this.Urls("lists", map[string]string{"pid": strconv.Itoa(pid)})}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		groupslist, _ := groups.Lists(map[string]string{}, 0, 1000)
		this.Data["groupslist"] = groupslist
		this.Data["pid"] = this.GetString("pid", "0")
		this.TplName = this.TplNames
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *MenuController) Edit() {
	if this.Ctx.Input.IsPost() {
		id, _ := this.GetInt("id", 0)
		pid, _ := this.GetInt("pid", 0)
		status, _ := this.GetInt("status", 0)
		sort, _ := this.GetInt("sort", 0)
		hide, _ := this.GetInt("hide", 0)
		groups_id, _ := this.GetInt("groups_id", 0)
		err := menu.Edit(id, this.GetString("title"), pid, sort, this.GetString("url"), hide, groups_id, status)
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "编辑成功", "url": this.Urls("lists", map[string]string{"pid": strconv.Itoa(pid)})}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		id, _ := this.GetInt("id")
		info := menu.GetInfo(id)
		this.Data["info"] = info
		groupslist, _ := groups.Lists(map[string]string{}, 0, 1000)
		fmt.Println(groupslist)
		for k, v := range groupslist {
			groupslist[k]["Select"] = false
			if int64(info.GroupsId) == v["Id"].(int64) {
				groupslist[k]["Select"] = true
			}
		}
		this.Data["groupslist"] = groupslist
		this.TplName = this.TplNames
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *MenuController) Del() {
	ids := strings.Split(this.GetString("id"), ",")
	if this.GetString("id") == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "msg": "您没有选择任何数据"}
	} else {
		err := menu.Del(ids)
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "删除成功"}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
	}
	this.ServeJSON()
	return
}
