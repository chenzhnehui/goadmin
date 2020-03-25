package config

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
	"strings"
	"time"
)

type Config struct {
	Id         int
	Name       string
	Need       int
	Status     int
	Type       int
	Group      int
	Extra      string
	Title      string
	Remark     string
	Value      string
	Sort       int
	CreateTime string
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func init() {
	orm.RegisterModel(new(Config))
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
	if condArr["title"] != "" {
		cond = cond.And("title__icontains", condArr["title"])
	}
	if condArr["name"] != "" {
		cond = cond.And("name__icontains", condArr["name"])
	}
	if condArr["type"] > "0" {
		cond = cond.And("type", condArr["type"])
	}
	var orderBy string
	if condArr["order"] != "" {
		orderBy = condArr["order"]
	} else {
		orderBy = "-id"
	}
	start := (page - 1) * offset
	var lists []orm.Params
	o := orm.NewOrm()
	qs := o.QueryTable("config")
	qs = qs.SetCond(cond)
	count, _ := qs.Count()
	qs.Limit(offset, start).OrderBy(orderBy).Values(&lists)
	return lists, count
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Add(Name string, Title string, Need int, Status int, Type int, Group int, Sort int, Extra string, Remark string, Value string) (int64, error) {
	o := orm.NewOrm()
	var config Config
	config.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	config.Name = strings.ToUpper(Name)
	config.Need = Need
	config.Status = Status
	config.Type = Type
	config.Group = Group
	config.Sort = Sort
	config.Title = Title
	config.Extra = strings.Replace(Extra, "\r", "", -1)
	config.Remark = Remark
	config.Value = strings.Replace(Value, "\r", "", -1)
	_, err := GetInfoByName(config.Name)
	if err == nil {
		return 0, errors.New("该标识已存在")
	}
	valid := validation.Validation{}
	valid.Required(config.Title, "Title").Message("标题不能为空")
	valid.Required(config.Name, "Name").Message("标识不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return 0, errors.New(err.Message)
		}
	}
	id, err := o.Insert(&config)
	utils.DelCache("Config.GetConfig.aconfigListCacheInfo")
	return id, err
}

/**
获取信息
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func GetInfoByName(name string) (Config, error) {
	o := orm.NewOrm()
	var config Config
	err := o.QueryTable("config").Filter("name", name).One(&config)
	if err == orm.ErrNoRows {
		return config, errors.New("数据不存在")
	}
	return config, nil
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Edit(id int, Name string, Title string, Need int, Status int, Type int, Group int, Sort int, Extra string, Remark string, Value string) error {
	o := orm.NewOrm()
	config := Config{Id: id}
	config.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	config.Name = strings.ToUpper(Name)
	config.Title = Title
	config.Need = Need
	config.Status = Status
	config.Type = Type
	config.Group = Group
	config.Sort = Sort
	config.Extra = strings.Replace(Extra, "\r", "", -1)
	config.Remark = Remark
	config.Value = strings.Replace(Value, "\r", "", -1)
	info, err := GetInfoByName(config.Name)
	if err == nil {
		if info.Id != id {
			return errors.New("该标识已存在")
		}
	}
	if id == 1 {
		if config.Name != "CONFIG_GROUP_LIST" {
			return errors.New("特殊标识，不能修改")
		}
		if config.Status != 1 {
			return errors.New("特殊标识，不能禁用")
		}
		if config.Type != 10 {
			return errors.New("特殊标识，不能修改字段类型")
		}
		valArr := strings.Split(config.Value, "\n")
		if valArr[0] != "0:不分组" {
			return errors.New("特殊标识，值0:不分组需要存在")
		}
	}
	valid := validation.Validation{}
	valid.Required(config.Title, "Title").Message("标题不能为空")
	valid.Required(config.Name, "Name").Message("标识不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
	_, err = o.Update(&config)
	utils.DelCache("Config.GetConfig.aconfigListCacheInfo")
	return err
}

//按name更新信息
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func UpdateConfigVal(name string, val string) error {
	val = strings.Replace(val, "\r", "", -1)
	if name == "CONFIG_GROUP_LIST" {
		valArr := strings.Split(val, "\n")
		if valArr[0] != "0:不分组" {
			return errors.New("特殊标识，值0:不分组  需要存在")
		}
	}
	_, err := orm.NewOrm().QueryTable("config").Filter("name", name).Update(orm.Params{"value": val})
	utils.DelCache("Config.GetConfig.aconfigListCacheInfo")
	return err
}

/**
获取信息
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func GetInfo(id int) Config {
	o := orm.NewOrm()
	var config Config
	o.QueryTable("config").Filter("id", id).One(&config)
	return config
}

//获取配置信息
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetConfig(key string, args ...interface{}) interface{} {
	var configListCacheInfo interface{}
	utils.GetCache("Config.GetConfig.aconfigListCacheInfo", &configListCacheInfo)
	if configListCacheInfo == nil {
		var configList []Config
		orm.NewOrm().QueryTable("config").Filter("status", 1).All(&configList)
		if configList == nil {
			return ""
		}
		utils.SetCache("Config.GetConfig.aconfigListCacheInfo", configList, -1)
		utils.GetCache("Config.GetConfig.aconfigListCacheInfo", &configListCacheInfo)
	}
	var infos map[string]interface{}
	for _, v := range configListCacheInfo.([]interface{}) {
		val := v.(map[string]interface{})
		if val["Name"] == key {
			infos = val
			break
		}
	}
	arrList := make(map[string]interface{})
	if len(args) > 0 {
		for _, v1 := range strings.Split(infos["Extra"].(string), "\n") {
			strs := strings.Split(v1, ":")
			arrList[strs[0]] = strs[1]
		}
		return arrList
	} else {
		if utils.Equal(10, infos["Type"]) {
			for _, v1 := range strings.Split(infos["Value"].(string), "\n") {
				strs := strings.Split(v1, ":")
				arrList[strs[0]] = strs[1]
			}
			return arrList
		}
		return infos["Value"]
	}
}

/**
删除数据
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func Del(ids []string) error {
	if utils.InArray("1", ids) {
		return errors.New("CONFIG_GROUP_LIST标识不能删除")
	}
	o := orm.NewOrm()
	_, errs := o.QueryTable("config").Filter("id__in", ids).Delete()
	utils.DelCache("Config.GetConfig.aconfigListCacheInfo")
	return errs
}
