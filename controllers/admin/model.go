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
	"goadmin/models/model"
	"goadmin/utils"
	"strconv"
	"strings"
)

type ModelController struct {
	BaseController
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *ModelController) Lists() {
	condArr := make(map[string]string)
	condArr["title"] = this.GetString("title")
	page, _ := this.GetInt("p")
	offset, err := beego.AppConfig.Int("pageoffset")
	if err != nil {
		offset = 15
	}

	list, count := model.Lists(condArr, page, offset)
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
func (this *ModelController) Add() {
	if this.Ctx.Input.IsPost() {
		extend, _ := this.GetInt("extend", 1)
		relation, _ := this.GetInt("relation", 1)
		_, err := model.Add(this.GetString("name"), this.GetString("title"), extend, relation)
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
func (this *ModelController) Edit() {
	if this.Ctx.Input.IsPost() {
		id, _ := this.GetInt("id", 0)
		err := model.Edit(id, utils.Inputs(this.Ctx.Request.PostForm))
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "编辑成功", "url": this.Urls("edit", map[string]string{"id": strconv.Itoa(id)})}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		id, _ := this.GetInt("id")
		this.Data["info"] = model.GetInfo(id)
		this.Data["attributelist"] = attribute.GetArrtibute(id)
		this.TplName = this.TplNames
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *ModelController) Del() {
	ids := strings.Split(this.GetString("id"), ",")
	if this.GetString("id") == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "msg": "您没有选择任何数据"}
	} else {
		err := model.Del(ids)
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "删除成功"}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
	}
	this.ServeJSON()
	return
}
