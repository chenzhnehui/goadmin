package groups

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
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"goadmin/utils"
	"time"
)

type Groups struct {
	Id         int
	Title      string
	Sort       int
	Url        string
	Status     int
	Hide       int
	Icon       string
	CreateTime string
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func init() {
	orm.RegisterModel(new(Groups))
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
	if condArr["hide"] != "" {
		cond = cond.And("hide", condArr["hide"])
	}
	start := (page - 1) * offset
	var lists []orm.Params
	o := orm.NewOrm()
	qs := o.QueryTable("groups")
	qs = qs.SetCond(cond)
	count, _ := qs.Count()
	qs.Limit(offset, start).OrderBy("-sort").Values(&lists)
	return lists, count
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Add(title string, sort int, url string, status int, hide int, icon string) (int64, error) {
	o := orm.NewOrm()
	var groups Groups
	groups.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	groups.Title = title
	groups.Icon = icon
	groups.Sort = sort
	groups.Url = url
	groups.Status = status
	groups.Hide = hide
	valid := validation.Validation{}
	valid.Required(groups.Title, "Title").Message("标题不能为空")
	valid.Required(groups.Icon, "Icon").Message("图标不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return 0, errors.New(err.Message)
		}
	}
	id, err := o.Insert(&groups)
	utils.DelCache("Menu.CheckPrivilege.groupslistCacheInfo")
	return id, err
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Edit(id int, title string, sort int, url string, status int, hide int, icon string) error {
	o := orm.NewOrm()
	groups := Groups{Id: id}
	groups.Title = title
	groups.Icon = icon
	groups.Sort = sort
	groups.Url = url
	groups.Status = status
	groups.Hide = hide
	groups.CreateTime = time.Now().Format("2006-01-02 15:04:05")

	valid := validation.Validation{}
	valid.Required(groups.Title, "Title").Message("标题不能为空")
	valid.Required(groups.Icon, "Icon").Message("图标不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
	_, err := o.Update(&groups)
	utils.DelCache("Menu.CheckPrivilege.groupslistCacheInfo")
	return err
}

/**
获取信息
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func GetInfo(id int) Groups {
	o := orm.NewOrm()
	var groups Groups
	o.QueryTable("groups").Filter("id", id).One(&groups)
	return groups
}

/**
删除数据
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func Del(ids []string) error {
	o := orm.NewOrm()
	_, errs := o.QueryTable("groups").Filter("id__in", ids).Delete()
	utils.DelCache("Menu.CheckPrivilege.groupslistCacheInfo")
	return errs
}
