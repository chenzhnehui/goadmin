package user

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

type Admin struct {
	Id         int
	Username   string
	Password   string
	Status     int
	Supper     int
	AccessId   string
	CreateTime string
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func init() {
	orm.RegisterModel(new(Admin))
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func CheckUserPwd(username, password string) (Admin, error) {
	var admin Admin
	if username == "" || password == "" {
		return admin, errors.New("用户名密码不能为空")
	}
	admin, err := GetInfoByUserName(username)
	if err != nil {
		return admin, err
	}
	if utils.Md5(password) != admin.Password {
		return admin, errors.New("密码不正确")
	}
	if admin.Status != 1 {
		return admin, errors.New("用户已被禁用")
	}
	return admin, nil
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
	if condArr["username"] != "" {
		cond = cond.And("username__icontains", condArr["username"])
	}

	start := (page - 1) * offset
	var lists []orm.Params
	o := orm.NewOrm()
	qs := o.QueryTable("admin")
	qs = qs.SetCond(cond)
	count, _ := qs.Count()
	qs.Limit(offset, start).OrderBy("-id").Values(&lists)
	return lists, count
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Add(username string, password string, supper int, status int, access_id string) (int64, error) {
	o := orm.NewOrm()
	var admin Admin
	admin.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	admin.Username = username
	admin.Password = utils.Md5(password)
	admin.Supper = supper
	admin.AccessId = access_id
	admin.Status = status
	valid := validation.Validation{}
	valid.Required(admin.Username, "Username").Message("用户名不能为空")
	valid.MinSize(password, 6, "Password").Message("密码字符不能低于6位")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return 0, errors.New(err.Message)
		}
	}
	_, err := GetInfoByUserName(username)
	if err == nil {
		return 0, errors.New("该用户名已存在")
	}
	id, err := o.Insert(&admin)
	return id, err
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Edit(id int, password string, supper int, status int, access_id string) error {
	o := orm.NewOrm()
	params := orm.Params{}
	valid := validation.Validation{}
	if password != "" {
		valid.MinSize(password, 6, "Password")
		if valid.HasErrors() {
			return errors.New("密码字符不能低于6位")
		}
		params["password"] = utils.Md5(password)
	}
	if supper >= 0 {
		params["supper"] = supper
	}
	if status >= 0 {
		params["status"] = status
	}
	if access_id != "" {
		params["access_id"] = access_id
	}
	_, err := o.QueryTable("admin").Filter("id", id).Update(params)
	utils.DelCache("Admin.GetInfoById.id." + fmt.Sprintf("%d", id))
	return err
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetInfoById(id int) (Admin, error) {
	var admin Admin
	var cacheInfo interface{}
	utils.GetCache("Admin.GetInfoById.id."+fmt.Sprintf("%d", id), &cacheInfo)
	if cacheInfo == nil {
		o := orm.NewOrm()
		err := o.QueryTable("admin").Filter("id", id).One(&admin)
		if err == orm.ErrNoRows {
			return admin, errors.New("用户不存在")
		}
		utils.SetCache("Admin.GetInfoById.id."+fmt.Sprintf("%d", id), admin, 0)
	} else {
		mapstructure.Decode(cacheInfo, &admin)
	}
	return admin, nil
}

/**
获取信息
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func GetInfoByUserName(username string) (Admin, error) {
	o := orm.NewOrm()
	var admin Admin
	err := o.QueryTable("admin").Filter("username", username).One(&admin)
	if err == orm.ErrNoRows {
		return admin, errors.New("用户不存在")
	}
	return admin, nil
}

/**
删除数据
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func Del(ids []string) error {
	o := orm.NewOrm()
	_, errs := o.QueryTable("admin").Filter("id__in", ids).Delete()
	return errs
}
