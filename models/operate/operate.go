package operate

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

type Operate struct {
	Id         int
	UserId     int
	Read       int
	Url        string
	Param      string
	CreateTime string
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func init() {
	orm.RegisterModel(new(Operate))
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
	qs := o.QueryTable("operate")
	qs = qs.SetCond(cond)
	count, _ := qs.Count()
	qs.Limit(offset, start).OrderBy("-id").Values(&lists)
	return lists, count
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Add(user_id int, read int, Url string, Param string) (int64, error) {
	o := orm.NewOrm()
	var operate Operate
	operate.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	operate.UserId = user_id
	operate.Url = Url
	operate.Read = read
	operate.Param = Param
	valid := validation.Validation{}
	valid.Required(operate.UserId, "UserId").Message("UserId不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return 0, errors.New(err.Message)
		}
	}
	id, err := o.Insert(&operate)
	return id, err
}

/**
获取信息
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func GetInfo(id int) Operate {
	o := orm.NewOrm()
	var operate Operate
	o.QueryTable("operate").Filter("id", id).One(&operate)
	return operate
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
	qs := o.QueryTable("operate")
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
	_, errs := o.QueryTable("operate").Filter("id__in", ids).Delete()
	return errs
}
