package menu

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
	"goadmin/models/groups"
	"goadmin/models/user"
	"goadmin/utils"
	"strconv"
	"strings"
	"time"
)

type Menu struct {
	Id         int
	Title      string
	Pid        int
	Sort       int
	Url        string
	Hide       int
	GroupsId   int
	Status     int
	CreateTime string
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func init() {
	orm.RegisterModel(new(Menu))
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
	if condArr["pid"] != "" {
		cond = cond.And("pid__in", strings.Split(condArr["pid"], ","))
	}
	if condArr["id"] != "" {
		cond = cond.And("id__in", strings.Split(condArr["id"], ","))
	}
	if condArr["hide"] != "" {
		cond = cond.And("hide", condArr["hide"])
	}
	if condArr["title"] != "" {
		cond = cond.And("title__icontains", condArr["title"])
	}
	start := (page - 1) * offset
	var lists []orm.Params
	o := orm.NewOrm()
	qs := o.QueryTable("menu")
	qs = qs.SetCond(cond)
	count, _ := qs.Count()
	qs.Limit(offset, start).OrderBy("-sort").Values(&lists)
	return lists, count
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Add(title string, pid int, sort int, url string, hide int, groups_id int, status int) (int64, error) {
	o := orm.NewOrm()
	var menu Menu
	menu.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	menu.GroupsId = groups_id
	menu.Pid = pid
	menu.Title = title
	menu.Sort = sort
	menu.Url = url
	menu.Hide = hide
	menu.Status = status

	valid := validation.Validation{}
	valid.Required(menu.Title, "Title").Message("标题不能为空")
	valid.Required(menu.Url, "Url").Message("链接不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return 0, errors.New(err.Message)
		}
	}
	id, err := o.Insert(&menu)
	return id, err
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Edit(id int, title string, pid int, sort int, url string, hide int, groups_id int, status int) error {
	o := orm.NewOrm()
	menu := Menu{Id: id}
	menu.GroupsId = groups_id
	menu.Pid = pid
	menu.Title = title
	menu.Sort = sort
	menu.Url = url
	menu.Hide = hide
	menu.Status = status
	menu.CreateTime = time.Now().Format("2006-01-02 15:04:05")

	valid := validation.Validation{}
	valid.Required(menu.Title, "Title").Message("标题不能为空")
	valid.Required(menu.Url, "Url").Message("链接不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
	_, err := o.Update(&menu)
	return err
}

/**
获取信息
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func GetInfo(id int) Menu {
	o := orm.NewOrm()
	var menu Menu
	o.QueryTable("menu").Filter("id", id).One(&menu)
	return menu
}

/**
删除数据
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func Del(ids []string) error {
	o := orm.NewOrm()
	var menu Menu
	err := o.QueryTable("menu").Filter("pid__in", ids).One(&menu)
	if err == orm.ErrNoRows {
		_, errs := o.QueryTable("menu").Filter("id__in", ids).Delete()
		return errs
	} else {
		return errors.New("请先删除子菜单")
	}
}

/**
返回左侧菜单
admin_id 存在则要检查权限
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func GetMenuGroup(url string, adminid interface{}) interface{} {
	admin_id := adminid.(int)
	groupslist, _ := groups.Lists(map[string]string{"hide": "0"}, 1, 1000)
	menulist, _ := Lists(map[string]string{"hide": "0"}, 1, 1000)
	j := 0
	for _, v := range groupslist {
		if admin_id > 0 && !CheckPrivilege("", "-"+strconv.FormatInt(v["Id"].(int64), 10), admin_id) {
			continue
		}
		groupslist[j] = v
		if v["Url"].(string) == url {
			groupslist[j]["Active"] = "active"
		} else {
			groupslist[j]["Active"] = ""
		}
		var menus []interface{}
		for _, v1 := range menulist {
			if v["Id"] == v1["GroupsId"] {
				if admin_id > 0 && !CheckPrivilege(v1["Url"].(string), "", admin_id) {
					continue
				}
				v1["Active"] = ""
				v1["Sonmenu"] = ""
				v1["Son"] = make([]interface{}, 0)
				if v1["Url"].(string) == url {
					v1["Active"] = "active"
					groupslist[j]["Active"] = "active open"
				}
				for _, v2 := range menulist {
					if v2["Pid"] == v1["Id"] && v2["Url"].(string) == url {
						v1["Active"] = "active"
						groupslist[j]["Active"] = "active open"
						v1["Sonmenu"] = v2["Title"]
						break
					}
				}
				menus = append(menus, v1)
			}
		}
		groupslist[j]["Son"] = menus
		j++
	}
	return groupslist[:j]
}

/**
分配权限获取所有权限菜单列表
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func GetPrivilegeMenu(menu_ids string) interface{} {
	menuIdsArr := strings.Split(menu_ids, ",")
	groupslist, _ := groups.Lists(map[string]string{}, 1, 1000)
	menulist, _ := Lists(map[string]string{}, 1, 1000)
	for k, v := range groupslist {
		groupslist[k]["Check"] = false
		if utils.InArray("-"+strconv.FormatInt(v["Id"].(int64), 10), menuIdsArr) {
			groupslist[k]["Check"] = true
		}
		var menus []interface{}
		for _, v1 := range menulist {
			var menuson []interface{}
			if v["Id"] == v1["GroupsId"] {
				v1["Check"] = false
				if utils.InArray(strconv.FormatInt(v1["Id"].(int64), 10), menuIdsArr) {
					v1["Check"] = true
				}
				for _, v2 := range menulist {
					if v2["Pid"] == v1["Id"] {
						v2["Check"] = false
						if utils.InArray(strconv.FormatInt(v2["Id"].(int64), 10), menuIdsArr) {
							v2["Check"] = true
						}
						menuson = append(menuson, v2)
					}
				}
				v1["Twoson"] = menuson
				menus = append(menus, v1)
			}
		}
		groupslist[k]["Son"] = menus
	}
	return groupslist
}

/**
检查用户是否有权限
args[0] menulist
args[1] groupslist
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func CheckPrivilege(url string, groups_id string, admin_id int) bool {

	userInfo, _ := user.GetInfoById(admin_id)
	if userInfo.Supper == 1 {
		return true
	}
	if userInfo.AccessId == "0" || userInfo.AccessId == "" {
		return false
	}
	var menuIds []string
	o := orm.NewOrm()
	var accessListCacheInfo interface{}
	utils.GetCache("Menu.CheckPrivilege.accessList.id."+fmt.Sprintf("%d", admin_id), &accessListCacheInfo)
	if accessListCacheInfo == nil {
		var accessList []access.Access
		o.QueryTable("access").Filter("status", 1).Filter("id__in", strings.Split(userInfo.AccessId, ",")).All(&accessList)
		if accessList == nil {
			return false
		}
		utils.SetCache("Menu.CheckPrivilege.accessList.id."+fmt.Sprintf("%d", admin_id), accessList, 0)
		utils.GetCache("Menu.CheckPrivilege.accessList.id."+fmt.Sprintf("%d", admin_id), &accessListCacheInfo)
	}
	for _, v := range accessListCacheInfo.([]interface{}) {
		val := v.(map[string]interface{})
		menuIds = append(menuIds, strings.Split(val["MenuIds"].(string), ",")...)
	}
	if groups_id != "" {
		return utils.InArray(groups_id, menuIds)
	}
	if url != "" {
		var menuListCacheInfo interface{}
		menuIdsStr := make(map[string]string)
		menuIdsStr["id"] = strings.Join(menuIds, ",")
		utils.GetCache("Menu.CheckPrivilege.menuList.id."+fmt.Sprintf("%s", menuIdsStr["id"]), &menuListCacheInfo)
		if menuListCacheInfo == nil {
			menuListCacheInfo, _ = Lists(menuIdsStr, 1, 1000)
			utils.SetCache("Menu.CheckPrivilege.menuList.id."+fmt.Sprintf("%s", menuIdsStr["id"]), menuListCacheInfo, 0)
			utils.GetCache("Menu.CheckPrivilege.menuList.id."+fmt.Sprintf("%s", menuIdsStr["id"]), &menuListCacheInfo)
		}
		for _, v := range menuListCacheInfo.([]interface{}) {
			val := v.(map[string]interface{})
			if url == val["Url"].(string) {
				return true
			}
		}

		var groupslistCacheInfo interface{}
		utils.GetCache("Menu.CheckPrivilege.groupslistCacheInfo", &groupslistCacheInfo)
		if groupslistCacheInfo == nil {
			groupList, _ := groups.Lists(map[string]string{}, 1, 1000)
			utils.SetCache("Menu.CheckPrivilege.groupslistCacheInfo", groupList, 0)
			utils.GetCache("Menu.CheckPrivilege.groupslistCacheInfo", &groupslistCacheInfo)
		}
		for _, val := range groupslistCacheInfo.([]interface{}) {
			v := val.(map[string]interface{})
			if v["Url"].(string) != "" && v["Url"].(string) == url {
				return utils.InArray("-"+strconv.FormatInt(int64(v["Id"].(float64)), 10), menuIds)
			}
		}
	}

	return false

}
