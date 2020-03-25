package attribute

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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/goinggo/mapstructure"
	"goadmin/models/category"
	"goadmin/models/menu"
	"goadmin/models/model"
	"goadmin/utils"
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Attribute struct {
	Id           int
	Name         string
	Title        string
	Field        string
	Type         int
	Value        string
	Remark       string
	IsShow       int
	Extra        string
	ModelId      int
	Need         int
	Status       int
	ValidateRule string
	ErrorInfo    string
	CreateTime   string
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func init() {
	orm.RegisterModel(new(Attribute))
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
	cond = cond.And("status", 1)
	if condArr["title"] != "" {
		cond = cond.And("title__icontains", condArr["title"])
	}
	if condArr["model_id"] != "" {
		cond = cond.And("model_id", condArr["model_id"])
	}
	if condArr["id"] != "" {
		cond = cond.And("id__in", strings.Split(condArr["id"], ","))
	}
	start := (page - 1) * offset
	var lists []orm.Params
	o := orm.NewOrm()
	qs := o.QueryTable("attribute")
	qs = qs.SetCond(cond)
	count, _ := qs.Count()
	qs.Limit(offset, start).OrderBy("id").Values(&lists)
	return lists, count
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Add(Name string, Title string, Field string, Type int, Value string, Remark string, IsShow int, Extra string, ModelId int, Need int, Status int, ValidateRule string, ErrorInfo string) (int64, error) {
	o := orm.NewOrm()
	var attribute Attribute
	attribute.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	attribute.Name = strings.ToLower(strings.Trim(Name, " "))
	attribute.Title = Title
	attribute.Field = Field
	attribute.Type = Type
	attribute.Remark = Remark
	attribute.IsShow = IsShow
	attribute.ModelId = ModelId
	attribute.Need = Need
	attribute.Status = Status
	attribute.ValidateRule = ValidateRule
	attribute.ErrorInfo = ErrorInfo
	attribute.Extra = strings.Replace(Extra, "\r", "", -1)
	attribute.Value = strings.Replace(Value, "\r", "", -1)
	valid := validation.Validation{}
	valid.Required(attribute.Title, "Title").Message("标题不能为空")
	valid.Required(attribute.Name, "Name").Message("名称不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return 0, errors.New(err.Message)
		}
	}
	modelInfo := model.GetInfo(ModelId)
	if utils.InArray(attribute.Name, GetTableField(modelInfo.Name)) {
		return 0, errors.New("该标识数据字段已存在")
	}
	if modelInfo.Extend > 0 {
		modelBase := model.GetInfo(modelInfo.Extend)
		if utils.InArray(attribute.Name, GetTableField(modelBase.Name)) {
			return 0, errors.New("该标识数据字段在基表中已存在")
		}
	}

	o.Begin()
	id, err1 := o.Insert(&attribute)
	sql := "alter table " + modelInfo.Name + " add column " + attribute.Name + " " + attribute.Field + " default  '" + attribute.Value + "'  COMMENT '" + attribute.Title + "' ;"
	_, err2 := o.Raw(sql).Exec()
	if err1 == nil && err2 == nil {
		o.Commit()
		utils.DelCache("Attribute.GetAttributeConfig.configAttributeListCacheInfo")
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
func GetTableField(table string) []string {
	var field []string
	o := orm.NewOrm()
	var tables orm.ParamsList
	sql := "select COLUMN_NAME from information_schema.COLUMNS where table_name = '" + table + "';"
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
func Edit(id int, Name string, Title string, Field string, Type int, Value string, Remark string, IsShow int, Extra string, ModelId int, Need int, Status int, ValidateRule string, ErrorInfo string) error {
	o := orm.NewOrm()
	attribute := Attribute{Id: id}
	attribute.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	attribute.Name = strings.ToLower(strings.Trim(Name, " "))
	attribute.Title = Title
	attribute.Field = Field
	attribute.Type = Type
	attribute.Remark = Remark
	attribute.IsShow = IsShow
	attribute.Extra = strings.Replace(Extra, "\r", "", -1)
	attribute.Value = strings.Replace(Value, "\r", "", -1)
	attribute.ModelId = ModelId
	attribute.Need = Need
	attribute.Status = Status
	attribute.ValidateRule = ValidateRule
	attribute.ErrorInfo = ErrorInfo
	valid := validation.Validation{}
	valid.Required(attribute.Title, "Title").Message("标题不能为空")
	valid.Required(attribute.Name, "Name").Message("名称不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
	modelInfo := model.GetInfo(ModelId)
	attributeInfo := GetInfo(id)
	attributeNewInfo := GetInfoByNameModelId(ModelId, attribute.Name)
	if attributeInfo.Name != attribute.Name && attributeNewInfo.Name != "" {
		return errors.New("该字段标识已存在")
	}
	utils.DelCache("Attribute.GetAttributeConfig.configAttributeListCacheInfo")
	o.Begin()
	_, err1 := o.Update(&attribute)
	sql := "alter table " + modelInfo.Name + " change  " + attributeInfo.Name + " " + attribute.Name + " " + attribute.Field + " default  '" + attribute.Value + "'  COMMENT '" + attribute.Title + "' ;"

	_, err2 := o.Raw(sql).Exec()
	if err1 == nil && err2 == nil {
		o.Commit()
		model.DelGetFeildSortCache(modelInfo)
		return nil
	} else {
		return o.Rollback()
	}
}

/**
获取信息
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func GetInfo(id int) Attribute {
	var attribute Attribute
	o := orm.NewOrm()
	o.QueryTable("attribute").Filter("id", id).One(&attribute)
	return attribute
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetInfoByNameModelId(model_id int, name string) Attribute {
	var attribute Attribute
	o := orm.NewOrm()
	o.QueryTable("attribute").Filter("model_id", model_id).Filter("name", name).One(&attribute)
	return attribute
}

//获取字段排序
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetFeildSort(modelid interface{}, fieldgrop interface{}) interface{} {
	var attributeGetFeildSortCacheInfo interface{}
	utils.GetCache("Attribute.GetFeildSort.attributeGetFeildSortCacheInfo"+fmt.Sprintf("%d,%s", modelid, fieldgrop), &attributeGetFeildSortCacheInfo)
	if attributeGetFeildSortCacheInfo != nil {
		return attributeGetFeildSortCacheInfo
	}

	model_id := utils.GetInt(modelid)
	field_grop := fieldgrop.(string)
	modelInfo := model.GetInfo(model_id)
	modelSort := utils.AnalysisStr(modelInfo.FieldSort, "|")
	var ids []string
	if modelSort != nil && modelSort != "" {
		for _, val := range modelSort.([]map[string]interface{}) {
			if val["key"] == field_grop {
				for _, v2 := range strings.Split(val["val"].(string), ",") {
					ids = append(ids, strings.Replace(v2, "\r", "", -1))
				}
			}
		}
	}
	var lists []Attribute
	if len(ids) == 0 {
		return lists
	}

	var attributeList []Attribute
	o := orm.NewOrm()
	o.QueryTable("attribute").Filter("status", 1).Filter("is_show", 1).Filter("model_id__in", []string{strconv.Itoa(model_id), "1"}).Filter("id__in", ids).All(&attributeList)
	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		for _, v := range attributeList {
			if id == v.Id {
				lists = append(lists, v)
			}
		}

	}
	utils.SetCache("Attribute.GetFeildSort.attributeGetFeildSortCacheInfo"+fmt.Sprintf("%d,%s", modelid, fieldgrop), lists, 0)
	return lists
}

//获取值
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func ListArrtibuteVal(val string, types int, extra string) interface{} {
	if !utils.InArray(types, []int{4, 5, 6, 8, 9, 11}) {
		return val
	}
	if val == "" {
		return val
	}

	if utils.InArray(types, []int{4, 5, 6}) {
		if extra == "" {
			return val
		}
		var strArr []string
		for _, v := range strings.Split(val, ",") {
			for _, val1 := range utils.AnalysisStr(extra, "\n").([]map[string]interface{}) {
				if utils.Equal(v, val1["key"]) {
					strArr = append(strArr, val1["val"].(string))
				}
			}
		}
		return strings.Join(strArr, ",")
	}
	if types == 9 {
		var jsons map[string]string
		json.Unmarshal([]byte(val), &jsons)
		return template.HTML("<a href='" + jsons["url"] + "' target='_blank'>" + jsons["name"] + "</a>")
	}
	if utils.InArray(types, []int{8, 11}) {
		var strPicArr template.HTML
		for _, v := range strings.Split(val, ",") {
			strPicArr = "&nbsp;" + template.HTML("<a href='"+v+"' target='_blank'><img src='"+v+"' style='width:50px;height:25px;'></a>")
		}
		return strPicArr
	}
	return val
}

//封装list列表信息
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func ListArrtibute(list []orm.Params, fieldKey []string, modelInfo model.Model, adminid interface{}, controller string, module string) ([]orm.Params, string) {
	var attributeList []Attribute
	var modelIds []int
	modelIds = append(modelIds, modelInfo.Id)
	if modelInfo.Extend > 0 {
		modelIds = append(modelIds, modelInfo.Extend)
	}
	o := orm.NewOrm()
	o.QueryTable("attribute").Filter("status", 1).Filter("is_show", 1).Filter("model_id__in", modelIds).Filter("name__in", fieldKey).All(&attributeList)

	var listGrids []string
	var fieldOperate []interface{}
	for _, v := range strings.Split(modelInfo.ListGrid, "\n") {
		str := strings.Split(v, ":")
		strOne := strings.ToLower(str[0])
		if v != "" && len(str[0]) > 2 && strOne[:2] == "__" {
			fields := make(map[string]interface{})
			strArr := strings.Split(strOne, "__")
			fields["field"] = strArr[1]
			fields["key"] = strOne
			strTitle := strings.Split(str[1], "|")
			fields["title"] = strTitle[0]
			fields["funcs"] = strTitle[1:]
			fieldOperate = append(fieldOperate, fields)
			listGrids = append(listGrids, fields["key"].(string)+":"+fields["title"].(string))
		} else {
			listGrids = append(listGrids, v)
		}
	}
	editPrivilege := menu.CheckPrivilege("/"+module+"/"+controller+"/edit", "", adminid.(int)) //编辑权限
	delPrivilege := menu.CheckPrivilege("/"+module+"/"+controller+"/del", "", adminid.(int))   //删除权限
	for k, v := range list {
		for k1, v1 := range v {
			for _, v2 := range fieldOperate {
				val2 := v2.(map[string]interface{})
				if utils.Equal(k1, val2["field"]) { //自定义字段函数处理
					var fieldValArr template.HTML
					if utils.InArray("EDIT", val2["funcs"]) && editPrivilege {
						fieldValArr += template.HTML("<a href='/" + module + "/" + controller + "/edit?" + k1 + "=" + v1.(string) + "&category_id=" + list[k]["category_id"].(string) + "'>编辑</a>")
					}
					if utils.InArray("DELETE", val2["funcs"]) && delPrivilege {
						if fieldValArr != "" {
							fieldValArr += " | "
						}
						fieldValArr += template.HTML("<a href='javascript:;' class='del-confirm' data-id='" + v1.(string) + "' data-url='/" + module + "/" + controller + "/del'>删除</a>")
					}
					list[k][val2["key"].(string)] = fieldValArr
				}
			}

			if utils.Equal("status", k1) {
				if utils.Equal(0, v1) {
					list[k][k1] = template.HTML("<span class='label label-sm label-warning'>禁用</span>")
				} else {
					list[k][k1] = template.HTML("<span class='label label-sm label-success'>正常</span>")
				}
				continue
			}
			for _, v2 := range attributeList {
				if utils.Equal(k1, v2.Name) {
					list[k][k1] = ListArrtibuteVal(v1.(string), v2.Type, v2.Extra)
					break
				}
			}
		}
	}
	return list, strings.Join(listGrids, "\n")
}

//获取全部字段
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetArrtibute(model_id int) []map[string]interface{} {
	modelInfo := model.GetInfo(model_id)

	modelSort := utils.AnalysisStr(modelInfo.FieldSort, "|")
	var ids []string
	if modelSort != nil && modelSort != "" {
		for _, v := range modelSort.([]map[string]interface{}) {
			for _, v2 := range strings.Split(v["val"].(string), ",") {
				ids = append(ids, strings.Replace(v2, "\r", "", -1))
			}
		}
	}
	lists, _ := Lists(map[string]string{"model_id": strconv.Itoa(model_id)}, 1, 1000)
	var arrList []map[string]interface{}
	for _, v := range lists {
		if utils.InArray(strconv.Itoa(int(v["Id"].(int64))), ids) {
			continue
		}
		arrList = append(arrList, v)
	}

	if modelInfo.Extend == 1 && modelInfo.Id != 1 {
		extlists, _ := Lists(map[string]string{"model_id": "1"}, 1, 1000)
		for _, v := range extlists {
			if utils.InArray(strconv.Itoa(int(v["Id"].(int64))), ids) {
				continue
			}
			arrList = append(arrList, v)
		}
	}
	return arrList
}

//获取单独字段属性
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetAttributeByName(key string, categoryid interface{}, args ...interface{}) interface{} {
	modelInfo := GetModelInfoByCategoryId(utils.GetInt(categoryid))
	var configAttributeListCacheInfo interface{}
	utils.GetCache("Attribute.GetAttributeConfig.configAttributeListCacheInfo", &configAttributeListCacheInfo)
	if configAttributeListCacheInfo == nil {
		var attributeList []Attribute
		orm.NewOrm().QueryTable("attribute").OrderBy("model_id").Filter("status", 1).All(&attributeList)
		if attributeList == nil {
			return ""
		}
		utils.SetCache("Attribute.GetAttributeConfig.configAttributeListCacheInfo", attributeList, -1)
		utils.GetCache("Attribute.GetAttributeConfig.configAttributeListCacheInfo", &configAttributeListCacheInfo)
	}
	var infos map[string]interface{}
	for _, v := range configAttributeListCacheInfo.([]interface{}) {
		val := v.(map[string]interface{})
		if val["Name"] == key {
			if modelInfo.Extend == 0 && utils.Equal(val["ModelId"], modelInfo.Id) {
				infos = val
				break
			} else {
				if utils.Equal(val["ModelId"], modelInfo.Id) || utils.Equal(val["ModelId"], 1) {
					infos = val
					break
				}
			}

		}
	}
	arrList := make(map[string]interface{})
	if len(args) > 0 {
		if utils.Equal("Extra", args[0]) {
			for _, v1 := range strings.Split(infos["Extra"].(string), "\n") {
				strs := strings.Split(v1, ":")
				arrList[strs[0]] = strs[1]
			}
			return arrList
		} else {
			return infos[args[0].(string)]
		}
	} else {
		if utils.Equal(10, infos["Type"]) {
			for _, v1 := range strings.Split(infos["Value"].(string), "\n") {
				strs := strings.Split(v1, ":")
				arrList[strs[0]] = strs[1]
			}
			return arrList
		}
		return infos
	}
}

//通过分类获取模型信息
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetModelInfoByCategoryId(category_id int) model.Model {
	var getModelInfoByCategoryIdCacheInfo interface{}
	utils.GetCache("Attribute.getModelInfoByCategoryIdCacheInfo"+fmt.Sprintf("%d", category_id), &getModelInfoByCategoryIdCacheInfo)
	if getModelInfoByCategoryIdCacheInfo == nil {
		categoryInfo := category.GetInfo(category_id)
		modelInfo := model.GetInfo(categoryInfo.ModelId)
		utils.SetCache("Attribute.getModelInfoByCategoryIdCacheInfo"+fmt.Sprintf("%d", category_id), modelInfo, 0)
		utils.GetCache("Attribute.getModelInfoByCategoryIdCacheInfo"+fmt.Sprintf("%d", category_id), &getModelInfoByCategoryIdCacheInfo)
	}
	var models model.Model
	mapstructure.Decode(getModelInfoByCategoryIdCacheInfo, &models)
	return models
}

/**
删除数据
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func Del(ids []string) error {
	o := orm.NewOrm()
	utils.DelCache("Attribute.GetAttributeConfig.configAttributeListCacheInfo")
	o.Begin()
	var lists []Attribute
	o.QueryTable("attribute").Filter("id__in", ids).All(&lists)
	var err error
	for _, v := range lists {
		modelInfo := model.GetInfo(v.ModelId)
		model.DelGetFeildSortCache(modelInfo)
		_, err = o.Raw("alter table " + modelInfo.Name + " drop column `" + v.Name + "`;").Exec()
		if err != nil {
			return o.Rollback()
		}
	}

	_, err = o.QueryTable("attribute").Filter("id__in", ids).Delete()
	if err == nil {
		o.Commit()
		return nil
	} else {
		return o.Rollback()
	}
}

//校验模型数据
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func checkModelData(post map[string]string) error {
	_, ok := post["category_id"]
	if !ok {
		return errors.New("分类ID不存在")
	}
	category_id, _ := strconv.Atoi(post["category_id"])
	modelInfo := GetModelInfoByCategoryId(category_id)
	if modelInfo.Name == "" {
		return errors.New("模型不存在")
	}
	var attributeId []string
	for _, v := range utils.AnalysisStr(modelInfo.FieldSort, "|").([]map[string]interface{}) {
		attributeId = append(attributeId, v["val"].(string))
	}
	var attributeList []Attribute
	orm.NewOrm().QueryTable("attribute").Filter("id__in", attributeId).All(&attributeList)

	valid := validation.Validation{}
	for _, v := range attributeList {
		errorInfo := v.ErrorInfo
		if v.ErrorInfo == "" {
			errorInfo = v.Title + "必须！"
		}
		if v.Need > 0 { //必填校验
			valid.Required(post[v.Name], v.Name).Message(errorInfo)
		}
		if post[v.Name] != "" && v.ValidateRule != "" { //存在就校验
			rules := strings.Split(v.ValidateRule, "|")
			switch rules[0] {
			case "Required":
				valid.Required(post[v.Name], v.Name).Message(errorInfo)
			case "Min":
				if len(rules) == 2 {
					val, _ := strconv.Atoi(rules[1])
					valid.Min(post[v.Name], val, v.Name).Message(errorInfo)
				}
			case "Max":
				if len(rules) == 2 {
					val, _ := strconv.Atoi(rules[1])
					valid.Max(post[v.Name], val, v.Name).Message(errorInfo)
				}
			case "Range":
				if len(rules) == 3 {
					val1, _ := strconv.Atoi(rules[1])
					val2, _ := strconv.Atoi(rules[2])
					valid.Range(post[v.Name], val1, val2, v.Name).Message(errorInfo)
				}
			case "MinSize":
				if len(rules) == 2 {
					val, _ := strconv.Atoi(rules[1])
					valid.MinSize(post[v.Name], val, v.Name).Message(errorInfo)
				}
			case "MaxSize":
				if len(rules) == 2 {
					val, _ := strconv.Atoi(rules[1])
					valid.MaxSize(post[v.Name], val, v.Name).Message(errorInfo)
				}
			case "Length":
				if len(rules) == 2 {
					val, _ := strconv.Atoi(rules[1])
					valid.Length(post[v.Name], val, v.Name).Message(errorInfo)
				}
			case "Alpha":
				valid.Alpha(post[v.Name], v.Name).Message(errorInfo)
			case "Numeric":
				valid.Numeric(post[v.Name], v.Name).Message(errorInfo)
			case "AlphaNumeric":
				valid.AlphaNumeric(post[v.Name], v.Name).Message(errorInfo)
			case "Match ":
				if len(rules[0]) == 2 {
					regexpStr, _ := regexp.CompilePOSIX(rules[1])
					valid.Match(post[v.Name], regexpStr, v.Name).Message(errorInfo)
				}
			case "NoMatch ":
				if len(rules[0]) == 2 {
					regexpStr, _ := regexp.CompilePOSIX(rules[1])
					valid.NoMatch(post[v.Name], regexpStr, v.Name).Message(errorInfo)
				}
			case "AlphaDash":
				valid.AlphaDash(post[v.Name], v.Name).Message(errorInfo)
			case "Email ":
				valid.Email(post[v.Name], v.Name).Message(errorInfo)
			case "IP":
				valid.IP(post[v.Name], v.Name).Message(errorInfo)
			case "Base64":
				valid.Base64(post[v.Name], v.Name).Message(errorInfo)
			case "Mobile":
				valid.Mobile(post[v.Name], v.Name).Message(errorInfo)
			case "Tel":
				valid.Tel(post[v.Name], v.Name).Message(errorInfo)
			case "Phone":
				valid.Phone(post[v.Name], v.Name).Message(errorInfo)
			case "ZipCode":
				valid.ZipCode(post[v.Name], v.Name).Message(errorInfo)
			}
		}

	}
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
	return nil
}

//获取模型列表
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func ModelLists(condArr map[string]string, page int, offset int, adminid interface{}, controller string, module string) ([]orm.Params, int64, model.Model, string, string) {
	if page < 1 {
		page = 1
	}
	if offset < 1 {
		offset, _ = beego.AppConfig.Int("pageoffset")
	}
	start := (page - 1) * offset
	category_id, _ := strconv.Atoi(condArr["category_id"])
	categoryInfo := category.GetInfo(category_id)
	modelInfo := model.GetInfo(categoryInfo.ModelId)
	var fieldKey []string
	var fieldOperate []string
	for _, v := range strings.Split(modelInfo.ListGrid, "\n") {
		str := strings.Split(v, ":")
		strOne := strings.ToLower(str[0])
		if len(strOne) > 1 && strOne[:2] == "__" {
			strArr := strings.Split(strOne, "__")
			fieldOperate = append(fieldOperate, strArr[1])
			continue
		}
		fieldKey = append(fieldKey, strOne)
	}
	if !utils.InArray("id", fieldKey) {
		fieldKey = append(fieldKey, "id")
	}
	if !utils.InArray("category_id", fieldKey) {
		fieldKey = append(fieldKey, "category_id")
	}
	//把 __id__加入数据库查询
	for _, v := range fieldOperate {
		if !utils.InArray(v, fieldKey) {
			fieldKey = append(fieldKey, v)
		}
	}

	var searchKey []map[string]string //搜索字段
	var searchKeyList []string        //搜索字段显示
	for _, v1 := range strings.Split(modelInfo.SearchKey, "\n") {
		str := strings.Split(v1, ":")
		strOne := strings.ToLower(str[0])
		strOneArr := strings.Split(strOne, "|")
		if len(strOneArr) > 0 {
			searchKeyField := make(map[string]string)
			if len(strOneArr) == 1 { //搜索添加判断，是like 还是其他方式
				searchKeyField["field"] = strOneArr[0]
				searchKeyField["condition"] = " = '" + strOneArr[0] + "'"
			} else {
				searchKeyField["field"] = strOneArr[0]
				searchKeyField["condition"] = " " + strOneArr[1]
			}
			searchKey = append(searchKey, searchKeyField)
		}
		if len(strOneArr) > 0 && len(str) > 1 {
			searchKeyList = append(searchKeyList, strOneArr[0]+":"+str[1])
		}
	}

	var sql string
	var sqlCount string
	where := " where " + modelInfo.Name + ".category_id = " + strconv.Itoa(category_id)

	if modelInfo.Extend == 0 { //独立模型
		for k, v := range condArr {
			for _, v1 := range searchKey {
				if v != "" && utils.Equal(k, v1["field"]) {
					where += " and " + k + " " + strings.Replace(v1["condition"], k, v, -1)
					break
				}
			}
		}
		sql = "select " + strings.Join(fieldKey, ",") + " from " + modelInfo.Name + where + " order by id desc" + " limit " + strconv.Itoa(start) + "," + strconv.Itoa(offset)
		sqlCount = "select count(*) from " + modelInfo.Name + where
	} else {
		modelBase := model.GetInfo(1)
		documentField := GetTableField(modelBase.Name)
		var fieldKeyJoin []string
		for _, v := range fieldKey {
			if utils.InArray(v, documentField) {
				fieldKeyJoin = append(fieldKeyJoin, modelBase.Name+"."+v)
			} else {
				fieldKeyJoin = append(fieldKeyJoin, modelInfo.Name+"."+v)
			}
		}
		for k, v := range condArr {
			for _, v1 := range searchKey {
				if v != "" && utils.Equal(k, v1["field"]) {
					if utils.InArray(v, documentField) {
						where += " and " + modelBase.Name + "." + k + " " + strings.Replace(v1["condition"], k, v, -1)
					} else {
						where += " and " + modelInfo.Name + "." + k + " " + strings.Replace(v1["condition"], k, v, -1)
					}
					break
				}
			}
		}

		baseSql := " from " + modelBase.Name + " left join " + modelInfo.Name + " on " + modelBase.Name + ".id=" + modelInfo.Name + ".id" + where
		sql = "select " + strings.Join(fieldKeyJoin, ",") + baseSql + "  order by " + modelBase.Name + ".id desc" + " limit " + strconv.Itoa(start) + "," + strconv.Itoa(offset)
		sqlCount = "select count(*) " + baseSql

	}
	var lists []orm.Params
	var counts orm.ParamsList
	o := orm.NewOrm()
	o.Raw(sql).Values(&lists)
	o.Raw(sqlCount).ValuesFlat(&counts)
	lists, listGrid := ListArrtibute(lists, fieldKey, modelInfo, adminid, controller, module) //封装字段显示

	count, _ := strconv.ParseInt(counts[0].(string), 10, 64)
	return lists, count, modelInfo, listGrid, strings.Join(searchKeyList, "\n")
}

//添加模型数据
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func AddModelData(post map[string]string) error {
	var err error
	err = checkModelData(post)
	if err != nil {
		return err
	}
	_, idOk := post["id"]
	if idOk {
		delete(post, "id")
	}
	_, statusOk := post["status"]
	if statusOk {
		delete(post, "status")
	}
	post["create_time"] = time.Now().Format("2006-01-02 15:04:05")
	o := orm.NewOrm()
	o.Begin()
	category_id, _ := strconv.Atoi(post["category_id"])
	modelInfo := GetModelInfoByCategoryId(category_id)
	modelInfoTable := GetTableField(modelInfo.Name)

	var keys []string
	var values []string
	for k, v := range post {
		if utils.InArray(k, modelInfoTable) {
			keys = append(keys, k)
			values = append(values, "'"+utils.SqlEscape(v)+"'")
		}
	}
	if modelInfo.Extend == 0 { //单表
		_, err = o.Raw("insert into " + modelInfo.Name + " (" + strings.Join(keys, ",") + ") values(" + strings.Join(values, ",") + ")").Exec()
		if err != nil {
			return o.Rollback()
		}
	} else {
		var baseKeys []string
		var baseValues []string
		modelBase := model.GetInfo(modelInfo.Extend)
		modelBaseTable := GetTableField(modelBase.Name)
		for k, v := range post {
			if utils.InArray(k, modelBaseTable) {
				baseKeys = append(baseKeys, k)
				baseValues = append(baseValues, "'"+utils.SqlEscape(v)+"'")
			}
		}
		_, err = o.Raw("insert into " + modelBase.Name + " (" + strings.Join(baseKeys, ",") + ") values(" + strings.Join(baseValues, ",") + ")").Exec()
		if err != nil {
			return o.Rollback()
		}

		var ids orm.ParamsList
		o.Raw("select LAST_INSERT_ID();").ValuesFlat(&ids)
		keys = append(keys, "id")
		values = append(values, ids[0].(string))
		_, err = o.Raw("insert into " + modelInfo.Name + " (" + strings.Join(keys, ",") + ") values(" + strings.Join(values, ",") + ")").Exec()
		if err != nil {
			return o.Rollback()
		}
	}
	o.Commit()
	return nil
}

//编辑模型数据
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func EditModelData(post map[string]string) error {
	var err error
	err = checkModelData(post)
	if err != nil {
		return err
	}
	_, idOk := post["id"]
	if !idOk || post["id"] == "" {
		return errors.New("ID不存在")
	}
	_, statusOk := post["status"]
	if statusOk {
		delete(post, "status")
	}
	o := orm.NewOrm()
	o.Begin()
	category_id, _ := strconv.Atoi(post["category_id"])
	modelInfo := GetModelInfoByCategoryId(category_id)
	modelInfoTable := GetTableField(modelInfo.Name)

	var values []string
	for k, v := range post {
		if utils.InArray(k, modelInfoTable) {
			values = append(values, " "+k+" = '"+utils.SqlEscape(v)+"'")
		}
	}
	_, err = o.Raw("update " + modelInfo.Name + " set " + strings.Join(values, ",") + " where id='" + utils.SqlEscape(post["id"]) + "' and category_id='" + utils.SqlEscape(post["category_id"]) + "'").Exec()
	if err != nil {
		return o.Rollback()
	}
	if modelInfo.Extend > 0 { //单表
		var baseValues []string
		modelBase := model.GetInfo(modelInfo.Extend)
		modelBaseTable := GetTableField(modelBase.Name)
		for k, v := range post {
			if utils.InArray(k, modelBaseTable) {
				baseValues = append(baseValues, " "+k+" = '"+utils.SqlEscape(v)+"'")
			}
		}
		_, err = o.Raw("update " + modelBase.Name + " set " + strings.Join(baseValues, ",") + " where id='" + utils.SqlEscape(post["id"]) + "' and category_id='" + utils.SqlEscape(post["category_id"]) + "'").Exec()
		if err != nil {
			return o.Rollback()
		}
	}
	o.Commit()
	return nil
}

//获取模型数据
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetModelData(id int, category_id int) map[string]string {
	o := orm.NewOrm()
	params := make(map[string]string)
	modelInfo := GetModelInfoByCategoryId(category_id)

	sql := "select * from " + modelInfo.Name + " where id='" + strconv.Itoa(id) + "' and category_id='" + strconv.Itoa(category_id) + "' limit 1"
	var resultList []orm.Params
	o.Raw(sql).Values(&resultList)
	if modelInfo.Extend > 0 {
		modelBase := model.GetInfo(modelInfo.Extend)
		sql := "select * from " + modelBase.Name + " where id='" + strconv.Itoa(id) + "' and category_id='" + strconv.Itoa(category_id) + "' limit 1"
		var resultBaseList []orm.Params
		o.Raw(sql).Values(&resultBaseList)
		for k, v := range resultBaseList[0] {
			params[k] = v.(string)
		}
	}
	for k, v := range resultList[0] {
		_, ok := params[k]
		if !ok {
			params[k] = v.(string)
		}
	}
	return params
}

//删除模型数据 types 类型 delete，setstatusyes，setstatusno
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func SetModelData(ids []string, category_id int, types string) error {
	o := orm.NewOrm()
	modelInfo := GetModelInfoByCategoryId(category_id)
	var err error
	o.Begin()
	switch types {
	case "delete":
		_, err = o.Raw("delete from " + modelInfo.Name + " where category_id=" + strconv.Itoa(category_id) + " and id in (" + strings.Join(ids, ",") + ")").Exec()
	case "setstatusyes":
		_, err = o.Raw("update " + modelInfo.Name + " set status='1' where category_id=" + strconv.Itoa(category_id) + " and id in (" + strings.Join(ids, ",") + ")").Exec()
	case "setstatusno":
		_, err = o.Raw("update " + modelInfo.Name + " set status='0' where category_id=" + strconv.Itoa(category_id) + " and id in (" + strings.Join(ids, ",") + ")").Exec()
	}
	if err != nil {
		return o.Rollback()
	}
	if modelInfo.Extend > 0 {
		modelBase := model.GetInfo(modelInfo.Extend)
		switch types {
		case "delete":
			_, err = o.Raw("delete from " + modelBase.Name + " where category_id=" + strconv.Itoa(category_id) + " and id in (" + strings.Join(ids, ",") + ")").Exec()
		case "setstatusyes":
			_, err = o.Raw("update " + modelBase.Name + " set status='1' where category_id=" + strconv.Itoa(category_id) + " and id in (" + strings.Join(ids, ",") + ")").Exec()
		case "setstatusno":
			_, err = o.Raw("update " + modelBase.Name + " set status='0' where category_id=" + strconv.Itoa(category_id) + " and id in (" + strings.Join(ids, ",") + ")").Exec()
		}
		if err != nil {
			return o.Rollback()
		}
	}
	o.Commit()
	return nil
}
