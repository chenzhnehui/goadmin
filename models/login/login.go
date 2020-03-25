package login

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
	"time"
)

type Login struct {
	Id         int
	UserId     int
	Read       int
	Ip         string
	CreateTime string
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func init() {
	orm.RegisterModel(new(Login))
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
	if condArr["read"] != "" {
		cond = cond.And("read", condArr["read"])
	}
	if condArr["user_id"] != "" {
		cond = cond.And("user_id", condArr["user_id"])
	}
	start := (page - 1) * offset
	var lists []orm.Params
	o := orm.NewOrm()
	qs := o.QueryTable("login")
	qs = qs.SetCond(cond)
	count, _ := qs.Count()
	qs.Limit(offset, start).OrderBy("-id").Values(&lists)
	return lists, count
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Add(user_id int, ip string, read int) (int64, error) {
	o := orm.NewOrm()
	var login Login
	login.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	login.UserId = user_id
	login.Ip = ip
	login.Read = read
	valid := validation.Validation{}
	valid.Required(login.UserId, "UserId").Message("UserId不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return 0, errors.New(err.Message)
		}
	}
	id, err := o.Insert(&login)
	return id, err
}

/**
获取信息
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func GetInfo(id int) Login {
	o := orm.NewOrm()
	var login Login
	o.QueryTable("login").Filter("id", id).One(&login)
	return login
}

//更新已读
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Reads(user_id int, id int) error {
	o := orm.NewOrm()
	params := orm.Params{}
	params["read"] = 1
	cond := orm.NewCondition()
	cond = cond.And("user_id", user_id)
	if id > 0 {
		cond = cond.And("id", id)
	}
	qs := o.QueryTable("login")
	qs = qs.SetCond(cond)
	_, err := qs.Update(params)
	return err
}

/**
删除数据
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func Del(ids []string) error {
	o := orm.NewOrm()
	_, errs := o.QueryTable("login").Filter("id__in", ids).Delete()
	return errs
}
