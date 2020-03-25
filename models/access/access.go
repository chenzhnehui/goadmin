package access

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
	"github.com/goinggo/mapstructure"
	"goadmin/utils"
	"time"
)

type Access struct {
	Id          int
	Title       string
	Status      int
	MenuIds     string
	CategoryIds string
	CreateTime  string
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func init() {
	orm.RegisterModel(new(Access))
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

	start := (page - 1) * offset
	var lists []orm.Params
	o := orm.NewOrm()
	qs := o.QueryTable("access")
	qs = qs.SetCond(cond)
	count, _ := qs.Count()
	qs.Limit(offset, start).OrderBy("-id").Values(&lists)
	return lists, count
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Add(title string, status int, menu_ids string, category_ids string) (int64, error) {
	o := orm.NewOrm()
	var access Access
	access.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	access.Title = title
	access.MenuIds = menu_ids
	access.CategoryIds = category_ids
	access.Status = status

	valid := validation.Validation{}
	valid.Required(access.Title, "Title").Message("标题不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return 0, errors.New(err.Message)
		}
	}
	id, err := o.Insert(&access)
	return id, err
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Edit(id int, title string, status int, menu_ids string, category_ids string) error {
	o := orm.NewOrm()
	params := orm.Params{}
	if title != "" {
		params["title"] = title
	}
	if menu_ids != "" {
		params["menu_ids"] = menu_ids
	}
	if category_ids != "" {
		params["category_ids"] = category_ids
	}
	if status >= 0 {
		params["status"] = status
	}
	_, err := o.QueryTable("access").Filter("id", id).Update(params)
	utils.DelCache("Access.GetInfo.id." + fmt.Sprintf("%d", id))
	return err
}

/**
获取信息
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func GetInfo(id int) Access {
	var access Access
	o := orm.NewOrm()
	o.QueryTable("access").Filter("id", id).One(&access)
	var cacheInfo interface{}
	utils.GetCache("Access.GetInfo.id."+fmt.Sprintf("%d", id), &cacheInfo)
	if cacheInfo == nil {
		o := orm.NewOrm()
		o.QueryTable("access").Filter("id", id).One(&access)
		utils.SetCache("Access.GetInfo.id."+fmt.Sprintf("%d", id), access, 0)
	} else {
		mapstructure.Decode(cacheInfo, &access)
	}
	return access
}

/**
删除数据
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func Del(ids []string) error {
	o := orm.NewOrm()
	_, errs := o.QueryTable("access").Filter("id__in", ids).Delete()
	return errs
}
