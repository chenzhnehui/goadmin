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
	"goadmin/models/category"
	"goadmin/models/model"
	"goadmin/utils"
	"strconv"
	"strings"
)

type CategoryController struct {
	BaseController
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *CategoryController) Lists() {
	condArr := make(map[string]string)
	condArr["title"] = this.GetString("title")
	list, _ := category.Lists(condArr, 1, 1000)
	this.Data["treelist"] = utils.GetTree(list, "Pid", "Id")
	this.TplName = this.TplNames
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *CategoryController) Add() {
	if this.Ctx.Input.IsPost() {
		pid, _ := this.GetInt("pid", 0)
		types, _ := this.GetInt("types", 1)
		status, _ := this.GetInt("status", 0)
		sort, _ := this.GetInt("sort", 0)
		hide, _ := this.GetInt("hide", 0)
		model_id, _ := this.GetInt("model_id", 0)
		_, err := category.Add(this.GetString("title"), pid, sort, this.GetString("url"), hide, model_id, status, types)
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "添加成功", "url": this.Urls("lists", map[string]string{"pid": strconv.Itoa(pid)})}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		pid, _ := this.GetInt("pid", 0)
		var pidInfo category.Category
		if pid >= 0 {
			pidInfo = category.GetInfo(pid)
			this.Data["pidInfo"] = pidInfo
		}
		modellist, _ := model.Lists(map[string]string{"status": "1"}, 0, 1000)
		for k, v := range modellist {
			modellist[k]["Selected"] = ""
			if utils.Equal(v["Id"], pidInfo.ModelId) {
				modellist[k]["Selected"] = "Selected"
			}
		}
		this.Data["modellist"] = modellist

		this.Data["pid"] = pid
		this.TplName = this.TplNames
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *CategoryController) Edit() {
	if this.Ctx.Input.IsPost() {
		id, _ := this.GetInt("id", 0)
		pid, _ := this.GetInt("pid", 0)
		status, _ := this.GetInt("status", 0)
		sort, _ := this.GetInt("sort", 0)
		hide, _ := this.GetInt("hide", 0)
		model_id, _ := this.GetInt("model_id", 0)
		types, _ := this.GetInt("types", 1)
		err := category.Edit(id, this.GetString("title"), pid, sort, this.GetString("url"), hide, model_id, status, types)
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "编辑成功", "url": this.Urls("lists", map[string]string{"pid": strconv.Itoa(pid)})}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		id, _ := this.GetInt("id")
		info := category.GetInfo(id)
		this.Data["info"] = info

		if info.Pid >= 0 {
			this.Data["pidInfo"] = category.GetInfo(info.Pid)
		}
		modellist, _ := model.Lists(map[string]string{"status": "1"}, 0, 1000)
		fmt.Println(modellist)
		for k, v := range modellist {
			modellist[k]["Select"] = false
			if int64(info.ModelId) == v["Id"].(int64) {
				modellist[k]["Select"] = true
			}
		}
		this.Data["modellist"] = modellist
		this.TplName = this.TplNames
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *CategoryController) Del() {
	ids := strings.Split(this.GetString("id"), ",")
	if this.GetString("id") == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "msg": "您没有选择任何数据"}
	} else {
		err := category.Del(ids)
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "删除成功"}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
	}
	this.ServeJSON()
	return
}
