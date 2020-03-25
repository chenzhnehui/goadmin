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
	"strings"
)

type ArticleController struct {
	BaseController
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *ArticleController) Lists() {
	condArr := utils.Inputs(this.Ctx.Request.Form)
	page, _ := this.GetInt("p")
	offset, err := beego.AppConfig.Int("pageoffset")
	if err != nil {
		offset = 15
	}
	list, count, modelInfo, listGrid, SearchKey := attribute.ModelLists(condArr, page, offset, this.GetSession("adminId"), this.CONTROLLER_NAME, this.MODULE_NAME)
	this.Data["paginator"] = pagination.SetPaginator(this.Ctx, offset, count)
	this.Data["list"] = list
	this.Data["count"] = count
	this.Data["field_title"] = utils.AnalysisStr(listGrid, "\n")
	this.Data["search_field"] = utils.AnalysisStr(SearchKey, "\n")
	this.Data["params"] = condArr
	if modelInfo.TemplateList == "" {
		modelInfo.TemplateList = this.TplNames
	}
	this.TplName = modelInfo.TemplateList
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *ArticleController) Add() {
	if this.Ctx.Input.IsPost() {
		err := attribute.AddModelData(utils.Inputs(this.Ctx.Request.PostForm))
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "添加成功", "url": this.Urls("lists", map[string]string{"category_id": this.GetString("category_id")})}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		category_id, _ := this.GetInt("category_id")
		modelInfo := attribute.GetModelInfoByCategoryId(category_id)
		this.Data["info"] = modelInfo
		this.Data["category_id"] = category_id
		if modelInfo.TemplateAdd != "" {
			this.TplName = modelInfo.TemplateAdd
		} else {
			this.TplName = this.TplNames
		}
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *ArticleController) Edit() {
	if this.Ctx.Input.IsPost() {
		err := attribute.EditModelData(utils.Inputs(this.Ctx.Request.PostForm))
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "编辑成功", "url": this.Urls("lists", map[string]string{"category_id": this.GetString("category_id")})}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		id, _ := this.GetInt("id")
		category_id, _ := this.GetInt("category_id")
		modelInfo := attribute.GetModelInfoByCategoryId(category_id)
		this.Data["info"] = modelInfo
		this.Data["category_id"] = category_id
		this.Data["idInfo"] = attribute.GetModelData(id, category_id)
		if modelInfo.TemplateEdit != "" {
			this.TplName = modelInfo.TemplateEdit
		} else {
			this.TplName = this.TplNames
		}
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *ArticleController) Del() {
	ids := strings.Split(this.GetString("id"), ",")
	category_id, _ := this.GetInt("category_id", 0)
	if this.GetString("id") == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "msg": "您没有选择任何数据"}
	} else {
		err := attribute.SetModelData(ids, category_id, "delete")
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
func (this *ArticleController) Setstatusyes() {
	ids := strings.Split(this.GetString("id"), ",")
	category_id, _ := this.GetInt("category_id", 0)
	if this.GetString("id") == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "msg": "您没有选择任何数据"}
	} else {
		err := attribute.SetModelData(ids, category_id, "setstatusyes")
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "启用成功"}
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
func (this *ArticleController) Setstatusno() {
	ids := strings.Split(this.GetString("id"), ",")
	category_id, _ := this.GetInt("category_id", 0)
	if this.GetString("id") == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "msg": "您没有选择任何数据"}
	} else {
		err := attribute.SetModelData(ids, category_id, "setstatusno")
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "禁用成功"}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
	}
	this.ServeJSON()
	return
}
