package routers

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
	"github.com/astaxie/beego"
	"goadmin/controllers"
	"goadmin/controllers/admin"
	"goadmin/controllers/home"
)

func init() {
	beego.Router("/", &home.IndexController{}, "get,post:Index")

	beego.AddNamespace(beego.NewNamespace("/admin",
		beego.NSAutoRouter(&admin.IndexController{}),
		beego.NSAutoRouter(&admin.UserController{}),
		beego.NSAutoRouter(&admin.MenuController{}),
		beego.NSAutoRouter(&admin.GroupsController{}),
		beego.NSAutoRouter(&admin.AccessController{}),
		beego.NSAutoRouter(&admin.LoginController{}),
		beego.NSAutoRouter(&admin.OperateController{}),
		beego.NSAutoRouter(&admin.ConfigController{}),
		beego.NSAutoRouter(&admin.FileController{}),
		beego.NSAutoRouter(&admin.ModelController{}),
		beego.NSAutoRouter(&admin.AttributeController{}),
		beego.NSAutoRouter(&admin.CategoryController{}),
		beego.NSAutoRouter(&admin.ArticleController{}),
		beego.NSAutoRouter(&admin.PublicController{}), //无权限控制
		beego.NSRouter("*", &admin.IndexController{}, "get,post:Index"),
	)) //后台管理 路由

	beego.AddNamespace(beego.NewNamespace("/home", //如果url不要 home,可以把如下路由放在最外层，配置文件 urlmode 设置为 2
		beego.NSAutoRouter(&home.IndexController{}),
		beego.NSRouter("*", &home.IndexController{}, "get,post:Index"),
	)) //前台 路由

	beego.Router("/ws", &controllers.WsController{}) //websocket路由
}
