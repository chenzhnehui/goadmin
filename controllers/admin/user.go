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
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils/pagination"
	"goadmin/models/access"
	"goadmin/models/login"
	"goadmin/models/user"
	"goadmin/utils"
	"strconv"
	"strings"
)

type UserController struct {
	BaseController
}

/**
登录验证
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func (this *UserController) Login() {
	if this.Ctx.Input.IsPost() {
		adminInfo, err := user.CheckUserPwd(this.GetString("username"), this.GetString("password"))
		if err != nil {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		} else {
			this.SetSession("adminInfo", adminInfo)
			this.SetSession("adminId", adminInfo.Id)
			go login.Add(adminInfo.Id, this.Ctx.Input.IP(), 0)
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "登录成功", "url": "/" + this.MODULE_NAME + "/"}
		}
		this.ServeJSON()
		return
	} else {
		this.TplName = this.TplNames
	}
}

/**
登录退出
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func (this *UserController) Logout() {
	utils.DelCache("Admin.GetInfoById.id." + fmt.Sprintf("%d", this.Ctx.Input.Session("adminId")))
	utils.DelCache("Menu.CheckPrivilege.accessList.id." + fmt.Sprintf("%d", this.Ctx.Input.Session("adminId")))
	this.DelSession("adminInfo")
	this.DelSession("adminId")
	if this.Ctx.Input.IsPost() {
		this.Data["json"] = map[string]interface{}{"code": 1, "msg": "退出成功", "url": this.Urls("login")}
		this.ServeJSON()
		return
	} else {
		this.Ctx.Redirect(302, this.Urls("login"))
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *UserController) Lists() {
	condArr := make(map[string]string)
	condArr["username"] = this.GetString("username")
	page, _ := this.GetInt("p")
	offset, err := beego.AppConfig.Int("pageoffset")
	if err != nil {
		offset = 15
	}

	list, count := user.Lists(condArr, page, offset)
	this.Data["paginator"] = pagination.SetPaginator(this.Ctx, offset, count)
	this.Data["list"] = list
	this.Data["count"] = count
	this.Data["params"] = condArr
	this.TplName = this.TplNames
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *UserController) Add() {
	if this.Ctx.Input.IsPost() {
		status, _ := this.GetInt("status", 0)
		supper, _ := this.GetInt("supper", 0)
		_, err := user.Add(this.GetString("username"), this.GetString("password"), supper, status, "0")
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "添加成功", "url": this.Urls("lists")}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		this.TplName = this.TplNames
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *UserController) Edit() {
	if this.Ctx.Input.IsPost() {
		id, _ := this.GetInt("id", 0)
		status, _ := this.GetInt("status", 0)
		supper, _ := this.GetInt("supper", 0)
		err := user.Edit(id, "", supper, status, "")
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "编辑成功", "url": this.Urls("lists")}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		id, _ := this.GetInt("id")
		userInfo, _ := user.GetInfoById(id)
		this.Data["info"] = userInfo
		this.TplName = this.TplNames
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *UserController) Password() {
	if this.Ctx.Input.IsPost() {
		id := this.GetSession("adminId").(int)
		password := this.GetString("password")
		repassword := this.GetString("repassword")
		if password != repassword {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": "重复密码不一致"}
			this.ServeJSON()
			return
		}
		err := user.Edit(id, password, -1, -1, "")
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "修改成功", "url": "/" + this.MODULE_NAME + "/"}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
		this.ServeJSON()
		return
	} else {
		this.TplName = this.TplNames
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *UserController) Privilege() {
	id, _ := this.GetInt("id")
	if this.Ctx.Input.IsPost() {
		err := user.Edit(id, "", -1, -1, this.GetString("access_id"))
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "操作成功"}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
	} else {
		userInfo, _ := user.GetInfoById(id)
		accessList, _ := access.Lists(map[string]string{}, 1, 1000)
		for k, v := range accessList {
			accessList[k]["Select"] = false
			for _, v1 := range strings.Split(userInfo.AccessId, ",") {
				if v1 == strconv.Itoa(int(v["Id"].(int64))) {
					accessList[k]["Select"] = true
					break
				}
			}
		}
		this.Data["json"] = accessList
	}
	this.ServeJSON()
	return
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *UserController) Del() {
	ids := strings.Split(this.GetString("id"), ",")
	if this.GetString("id") == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "msg": "您没有选择任何数据"}
	} else {
		err := user.Del(ids)
		if err == nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "msg": "删除成功"}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": err.Error()}
		}
	}
	this.ServeJSON()
	return
}
