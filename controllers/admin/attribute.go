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
	"goadmin/models/attribute"
	"goadmin/utils"
	"strconv"
	"strings"
)

type AttributeController struct {
	BaseController
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *AttributeController) Lists() {
	condArr := make(map[string]string)
	condArr["title"] = this.GetString("title")
	condArr["model_id"] = this.GetString("model_id", "")
	page, _ := this.GetInt("p")
	offset, err := beego.AppConfig.Int("pageoffset")
	if err != nil {
		offset = 15
	}

	list, count := attribute.Lists(condArr, page, offset)
	this.Data["paginator"] = pagination.SetPaginator(this.Ctx, offset, count)
	this.Data["list"] = list
	this.Data["count"] = count
	this.Data["params"] = condArr
	this.Data["model_id"] = this.GetString("model_id")
	this.TplName = this.TplNames
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *AttributeController) Add() {
	if this.Ctx.Input.IsPost() {
		name := this.GetString("name")
		title := this.GetString("title")
		field := this.GetString("field")
		types, _ := this.GetInt("type", 1)
		remark := this.GetString("remark")
		is_show, _ := this.GetInt("is_show", 1)
		status, _ := this.GetInt("status", 1)
		model_id, _ := this.GetInt("model_id", 1)
		need, _ := this.GetInt("need", 0)
		validate_rule := this.GetString("validate_rule")
		error_info := this.GetString("error_info")
		extra := strings.Trim(this.GetString("extra", ""), " ")
		value := strings.Trim(this.GetString("value", ""), " ")
		if utils.InArray(types, []int{4, 5, 6}) && extra == "" {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": "单选多选枚举类型请填写多选值"}
			this.ServeJSON()
			return
		}
		_, err := attribute.Add(name, title, field, types, value, remark, is_show, extra, model_id, need, status, validate_rule, error_info)
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "添加成功", "url": this.Urls("lists", map[string]string{"model_id": strconv.Itoa(model_id)})}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		this.Data["model_id"] = this.GetString("model_id")
		this.TplName = this.TplNames
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *AttributeController) Edit() {
	if this.Ctx.Input.IsPost() {
		name := this.GetString("name")
		title := this.GetString("title")
		field := this.GetString("field")
		types, _ := this.GetInt("type", 1)
		remark := this.GetString("remark")
		is_show, _ := this.GetInt("is_show", 1)
		status, _ := this.GetInt("status", 1)
		id, _ := this.GetInt("id", 0)
		model_id, _ := this.GetInt("model_id", 1)
		need, _ := this.GetInt("need", 0)
		validate_rule := this.GetString("validate_rule")
		error_info := this.GetString("error_info")
		extra := strings.Trim(this.GetString("extra", ""), " ")
		value := strings.Trim(this.GetString("value", ""), " ")
		if utils.InArray(types, []int{4, 5, 6}) && extra == "" {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": "单选多选枚举类型请填写多选值"}
			this.ServeJSON()
			return
		}
		err := attribute.Edit(id, name, title, field, types, value, remark, is_show, extra, model_id, need, status, validate_rule, error_info)
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "编辑成功", "url": this.Urls("lists", map[string]string{"model_id": strconv.Itoa(model_id)})}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		id, _ := this.GetInt("id")
		this.Data["info"] = attribute.GetInfo(id)
		this.TplName = this.TplNames
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *AttributeController) Del() {
	ids := strings.Split(this.GetString("id"), ",")
	if this.GetString("id") == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "msg": "您没有选择任何数据"}
	} else {
		var err error
		if len(ids) == 1 {
			err = attribute.Del(ids)
		} else {
			err = nil
			go func() {
				attribute.Del(ids)
			}()
		}
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "删除成功"}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
	}
	this.ServeJSON()
	return
}
