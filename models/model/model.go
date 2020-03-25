package model

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
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"goadmin/utils"
	"strings"
	"time"
)

type Model struct {
	Id           int
	Name         string
	Title        string
	Extend       int
	Relation     int
	FieldSort    string
	FieldGroup   string
	TemplateList string
	TemplateAdd  string
	TemplateEdit string
	ListGrid     string
	SearchKey    string
	Status       int
	CreateTime   string
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func init() {
	orm.RegisterModel(new(Model))
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Lists(condArr map[string]string, page int, offset int) ([]orm.Params, int64) {
	if page < 1 {
		page = 1
	}
	if offset < 1 {
		offset, _ = beego.AppConfig.Int("pageoffset")
	}
	cond := orm.NewCondition()
	if condArr["title"] != "" {
		cond = cond.And("title__icontains", condArr["title"])
	}
	if condArr["status"] != "" {
		cond = cond.And("status", condArr["status"])
	}
	start := (page - 1) * offset
	var lists []orm.Params
	o := orm.NewOrm()
	qs := o.QueryTable("model")
	qs = qs.SetCond(cond)
	count, _ := qs.Count()
	qs.Limit(offset, start).OrderBy("id").Values(&lists)
	return lists, count
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Add(Name string, Title string, Extend int, Relation int) (int64, error) {
	o := orm.NewOrm()
	var model Model
	model.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	model.Name = strings.ToLower(strings.Trim(Name, " "))
	model.Title = Title
	model.Extend = Extend
	model.FieldGroup = "1:基础"
	model.SearchKey = "id:请输入ID编号\nstatus:状态"
	model.ListGrid = "id:编号\nstatus:状态\ncreate_time:创建时间\n__id__:操作|EDIT|DELETE"
	model.Relation = Relation
	model.Status = 1

	var tables orm.ParamsList
	o.Raw("show tables").ValuesFlat(&tables)
	for _, v := range tables {
		if v.(string) == model.Name {
			return 0, errors.New("该标识数据表已存在")
		}
	}
	valid := validation.Validation{}
	valid.Required(model.Title, "Title").Message("标题不能为空")
	valid.Required(model.Name, "Name").Message("名称不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return 0, errors.New(err.Message)
		}
	}
	o.Begin()
	id, err1 := o.Insert(&model)
	sql := "CREATE TABLE IF NOT EXISTS `" + model.Name + "`( `id` INT UNSIGNED AUTO_INCREMENT,`category_id` int(10) NOT NULL default 0,`status` int(1) NOT NULL default 1, `create_time` datetime NOT NULL, PRIMARY KEY ( `id` ),KEY `category_id` (`category_id`) USING BTREE,KEY `status` (`status`) USING BTREE) comment='" + model.Title + "' ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	_, err2 := o.Raw(sql).Exec()
	if err1 == nil && err2 == nil {
		o.Commit()
		return id, nil
	} else {
		err := o.Rollback()
		return id, err
	}
}

//获取表字段
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetTableField(table []string) []string {
	var field []string
	o := orm.NewOrm()
	var tables orm.ParamsList
	sql := "select COLUMN_NAME from information_schema.COLUMNS where table_name in (" + strings.Join(table, ",") + ");"
	o.Raw(sql).ValuesFlat(&tables)
	for _, v := range tables {
		field = append(field, v.(string))
	}
	return field
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Edit(id int, post map[string]string) error {
	o := orm.NewOrm()
	for k, v := range post {
		post[k] = strings.Replace(v, "\r", "", -1)
	}
	params := orm.Params{}
	valid := validation.Validation{}
	valid.Required(post["title"], "Title").Message("标题不能为空")
	valid.Required(post["name"], "Name").Message("名称不能为空")
	valid.Required(post["field_group"], "field_group").Message("表单显示不能为空")
	valid.Required(post["list_grid"], "list_grid").Message("表单列表定义不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
	info := GetInfo(id)
	var tables []string
	tables = append(tables, "'"+info.Name+"'")
	if info.Extend >= 0 {
		infoBase := GetInfo(info.Extend)
		tables = append(tables, "'"+infoBase.Name+"'")
	}
	fields := GetTableField(tables) //表字段
	for _, v := range strings.Split(post["list_grid"], "\n") {
		str := strings.Split(v, ":")
		strOne := strings.ToLower(str[0])
		if len(strOne) > 1 && strOne[:2] == "__" {
			strArr := strings.Split(strOne, "__")
			if !utils.InArray(strArr[1], fields) {
				return errors.New("列表定义字段 " + strOne + "对应表字段不存在")
			}
		} else {
			if !utils.InArray(strOne, fields) {
				return errors.New("列表定义字段 " + strOne + "对应表字段不存在")
			}
		}
	}

	if post["search_key"] != "" {
		for _, v1 := range strings.Split(post["search_key"], "\n") {
			str := strings.Split(v1, ":")
			strOne := strings.ToLower(str[0])
			strOneArr := strings.Split(strOne, "|")
			if len(strOneArr[0]) > 0 && strOneArr[0] != "" {
				if !utils.InArray(strOneArr[0], fields) {
					return errors.New("搜索字段 " + strOneArr[0] + "对应表字段不存在")
				}
			}
		}
	}

	params["title"] = post["title"]
	params["name"] = strings.ToLower(strings.Trim(post["name"], " "))
	params["extend"] = post["extend"]
	params["relation"] = post["relation"]
	params["field_sort"] = post["field_sort"]
	params["field_group"] = post["field_group"]
	params["template_list"] = post["template_list"]
	params["template_add"] = post["template_add"]
	params["template_edit"] = post["template_edit"]
	params["list_grid"] = post["list_grid"]
	params["search_key"] = post["search_key"]
	params["status"] = post["status"]

	_, err := o.QueryTable("model").Filter("id", id).Update(params)
	if info.Name != params["name"].(string) {
		o.Raw("alter table " + info.Name + " rename to " + params["name"].(string) + ";").Exec()
	}
	if info.Title != params["title"].(string) {
		o.Raw("alter table " + params["name"].(string) + " comment '" + params["title"].(string) + "';").Exec()
	}
	DelGetFeildSortCache(info, params["field_group"])
	return err
}

/**
获取信息
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func GetInfo(id int) Model {
	var model Model
	o := orm.NewOrm()
	o.QueryTable("model").Filter("id", id).One(&model)
	return model
}

//删除缓存
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func DelGetFeildSortCache(modelInfo Model, args ...interface{}) {
	var FieldGroup string
	if len(args) > 0 {
		FieldGroup = args[0].(string)
	} else {
		FieldGroup = modelInfo.FieldGroup
	}
	FieldGroupStr := utils.AnalysisStr(FieldGroup, ",")
	if FieldGroupStr != nil && FieldGroupStr != "" {
		for _, v := range FieldGroupStr.([]map[string]interface{}) {
			utils.DelCache("Attribute.GetFeildSort.attributeGetFeildSortCacheInfo" + fmt.Sprintf("%d,%s", modelInfo.Id, v["key"].(string)))
		}
	}
}

/**
删除数据
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func Del(ids []string) error {
	o := orm.NewOrm()
	if utils.InArray("1", ids) {
		return errors.New("基础模型不能删除")
	}
	o.Begin()
	var lists []Model
	var err error
	o.QueryTable("model").Filter("id__in", ids).All(&lists)
	for _, v := range lists {
		DelGetFeildSortCache(v)
		_, err := o.Raw("drop table " + v.Name + ";").Exec()
		if err != nil {
			return o.Rollback()
		}
	}

	_, err = o.QueryTable("model").Filter("id__in", ids).Delete()
	if err == nil {
		o.Commit()
		return nil
	} else {
		return o.Rollback()
	}
}
