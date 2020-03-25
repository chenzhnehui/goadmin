package category

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
	"goadmin/models/access"
	"goadmin/models/user"
	"goadmin/utils"
	"strconv"
	"strings"
	"time"
)

type Category struct {
	Id         int
	Title      string
	Pid        int
	Sort       int
	Url        string
	Hide       int
	ModelId    int
	Types      int
	Status     int
	CreateTime string
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func init() {
	orm.RegisterModel(new(Category))
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
	if condArr["pid"] != "" {
		cond = cond.And("pid__in", strings.Split(condArr["pid"], ","))
	}
	if condArr["id"] != "" {
		cond = cond.And("id__in", strings.Split(condArr["id"], ","))
	}
	if condArr["hide"] != "" {
		cond = cond.And("hide", condArr["hide"])
	}
	if condArr["types"] != "" {
		cond = cond.And("types", condArr["types"])
	}
	if condArr["status"] != "" {
		cond = cond.And("status", condArr["status"])
	}
	if condArr["title"] != "" {
		cond = cond.And("title__icontains", condArr["title"])
	}
	var sort string
	if condArr["sort"] == "" {
		sort = "sort"
	} else {
		sort = condArr["sort"]
	}
	start := (page - 1) * offset
	var lists []orm.Params
	o := orm.NewOrm()
	qs := o.QueryTable("category")
	qs = qs.SetCond(cond)
	count, _ := qs.Count()
	qs.Limit(offset, start).OrderBy(sort).Values(&lists)
	return lists, count
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Add(title string, pid int, sort int, url string, hide int, model_id int, status int, types int) (int64, error) {
	o := orm.NewOrm()
	var category Category
	category.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	category.ModelId = model_id
	category.Pid = pid
	category.Title = title
	category.Sort = sort
	category.Url = url
	category.Hide = hide
	category.Status = status
	category.Types = types

	valid := validation.Validation{}
	valid.Required(category.Title, "Title").Message("标题不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return 0, errors.New(err.Message)
		}
	}
	id, err := o.Insert(&category)
	return id, err
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Edit(id int, title string, pid int, sort int, url string, hide int, model_id int, status int, types int) error {
	o := orm.NewOrm()
	category := Category{Id: id}
	category.ModelId = model_id
	category.Pid = pid
	category.Title = title
	category.Sort = sort
	category.Url = url
	category.Hide = hide
	category.Status = status
	category.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	category.Types = types
	valid := validation.Validation{}
	valid.Required(category.Title, "Title").Message("标题不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
	_, err := o.Update(&category)
	utils.DelCache("Attribute.getModelInfoByCategoryIdCacheInfo" + fmt.Sprintf("%d", id))
	return err
}

/**
获取信息
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func GetInfo(id int) Category {
	o := orm.NewOrm()
	var category Category
	o.QueryTable("category").Filter("id", id).One(&category)
	return category
}

//获取文章分类权限数据
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func IsArticlePrivilege(adminid interface{}, categoryid int) bool {
	category_id := strconv.Itoa(categoryid)
	admin_id := adminid.(int)
	userInfo, _ := user.GetInfoById(admin_id)
	if userInfo.Supper == 1 {
		return true
	}
	if category_id == "" || category_id == "0" {
		return false
	}
	var categoryIds []string
	if userInfo.AccessId != "0" && userInfo.AccessId != "" {
		o := orm.NewOrm()
		var accessListCacheInfo interface{}
		utils.GetCache("Menu.CheckPrivilege.accessList.id."+fmt.Sprintf("%d", admin_id), &accessListCacheInfo)
		if accessListCacheInfo == nil {
			var accessList []access.Access
			o.QueryTable("access").Filter("status", 1).Filter("id__in", strings.Split(userInfo.AccessId, ",")).All(&accessList)
			utils.SetCache("Menu.CheckPrivilege.accessList.id."+fmt.Sprintf("%d", admin_id), accessList, 0)
			utils.GetCache("Menu.CheckPrivilege.accessList.id."+fmt.Sprintf("%d", admin_id), &accessListCacheInfo)
		}
		for _, v := range accessListCacheInfo.([]interface{}) {
			val := v.(map[string]interface{})
			if val["CategoryIds"] != nil {
				for _, v1 := range strings.Split(val["CategoryIds"].(string), ",") {
					categoryIds = append(categoryIds, v1)
				}
			}
		}
		return utils.InArray(category_id, categoryIds)
	}
	return false
}

//获取文章分类权限数据
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetArticleMenuGroup(menugroup interface{}, adminid interface{}, category_id int, action string, controller string, types string) interface{} {
	admin_id := adminid.(int)
	userInfo, _ := user.GetInfoById(admin_id)
	listMap := map[string]string{"status": "1", "hide": "0"}
	if types != "" {
		listMap["types"] = types
	}
	var categoryIds []string
	categoryIds = append(categoryIds, "-1")
	if userInfo.Supper == 0 {
		if userInfo.AccessId != "0" && userInfo.AccessId != "" {
			o := orm.NewOrm()
			var accessListCacheInfo interface{}
			utils.GetCache("Menu.CheckPrivilege.accessList.id."+fmt.Sprintf("%d", admin_id), &accessListCacheInfo)
			if accessListCacheInfo == nil {
				var accessList []access.Access
				o.QueryTable("access").Filter("status", 1).Filter("id__in", strings.Split(userInfo.AccessId, ",")).All(&accessList)
				utils.SetCache("Menu.CheckPrivilege.accessList.id."+fmt.Sprintf("%d", admin_id), accessList, 0)
				utils.GetCache("Menu.CheckPrivilege.accessList.id."+fmt.Sprintf("%d", admin_id), &accessListCacheInfo)
			}
			for _, v := range accessListCacheInfo.([]interface{}) {
				val := v.(map[string]interface{})
				categoryIds = append(categoryIds, strings.Split(val["CategoryIds"].(string), ",")...)
			}
		}
		listMap["id"] = strings.Join(categoryIds, ",")
	}
	lists, _ := Lists(listMap, 0, 1000)
	var categoryList []orm.Params
	menus := make(map[int][]interface{})
	pidParent := 0
	for _, v := range lists {
		pid := int(v["Pid"].(int64))
		v["Url"] = "/" + controller + "/lists?category_id=" + strconv.Itoa(utils.GetInt(v["Id"]))
		v["Active"] = ""
		if utils.Equal(category_id, v["Id"]) {
			v["Active"] = "active"
			if action == "edit" {
				v["Sonmenu"] = "编辑"
			}
			if action == "add" {
				v["Sonmenu"] = "添加"
			}
			if pid > 0 {
				pidParent = pid
			}
		}
		if pid > 0 {
			menus[pid] = append(menus[pid], v)
		} else {
			categoryList = append(categoryList, v)
		}
	}
	for k, v := range categoryList {
		if utils.Equal(pidParent, v["Id"]) {
			categoryList[k]["Active"] = "active open"
		}
		categoryList[k]["Son"] = menus[int(v["Id"].(int64))]
	}

	for _, v := range menugroup.([]orm.Params) {
		article := false
		for _, v1 := range v["Son"].([]interface{}) {
			val1 := v1.(orm.Params)
			if val1["Url"] == "/"+controller+"/lists" {
				article = true
				break
			}
		}
		if article {
			var menu []interface{}
			for _, v1 := range v["Son"].([]interface{}) {
				val1 := v1.(orm.Params)
				if val1["Url"] != "/"+controller+"/lists" {
					menu = append(menu, v1)
				}
			}
			for _, v3 := range categoryList {
				menu = append(menu, v3)
			}
			v["Son"] = menu
		}
	}
	return menugroup
}

//获取分类权限
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetPrivilegeCategory(category_ids string) interface{} {
	lists, _ := Lists(map[string]string{"status": "1"}, 0, 1000)
	var categoryList []orm.Params
	categorys := make(map[int][]interface{})
	for _, v := range lists {
		pid := int(v["Pid"].(int64))
		v["Check"] = false
		if utils.InArray(strconv.FormatInt(v["Id"].(int64), 10), strings.Split(category_ids, ",")) {
			v["Check"] = true
		}
		if pid > 0 {
			categorys[pid] = append(categorys[pid], v)
		} else {
			categoryList = append(categoryList, v)
		}
	}
	for k, v := range categoryList {
		categoryList[k]["Son"] = categorys[int(v["Id"].(int64))]
	}
	return categoryList
}

/**
删除数据
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func Del(ids []string) error {
	o := orm.NewOrm()
	var category Category
	for _, id := range ids {
		utils.DelCache("Attribute.getModelInfoByCategoryIdCacheInfo" + fmt.Sprintf("%d", id))
	}
	err := o.QueryTable("category").Filter("pid__in", ids).One(&category)
	if err == orm.ErrNoRows {
		_, errs := o.QueryTable("category").Filter("id__in", ids).Delete()
		return errs
	} else {
		return errors.New("请先删除子分类")
	}
}
