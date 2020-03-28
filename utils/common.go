package utils

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
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"net/smtp"
	"strconv"
	"strings"
)

/**
 *  to: example@example.com;example1@163.com;example2@sina.com.cn;...
 *  subject:The subject of mail
 *  body: The content of mail，可以是 html
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func SendMail(to string, subject string, body string) error {
	user := beego.AppConfig.String("mailuser")
	password := beego.AppConfig.String("mailpassword")
	host := beego.AppConfig.String("mailhost")

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	content_type = "Content-type:text/html;charset=utf-8"

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

// 判断某一个值是否含在切片之中
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func InArray(need interface{}, haystack interface{}) bool {
	switch key := need.(type) {
	case int:
		for _, item := range haystack.([]int) {
			if item == key {
				return true
			}
		}
	case string:
		for _, item := range haystack.([]string) {
			if item == key {
				return true
			}
		}
	case int64:
		for _, item := range haystack.([]int64) {
			if item == key {
				return true
			}
		}
	case float64:
		for _, item := range haystack.([]float64) {
			if item == key {
				return true
			}
		}
	default:
		return false
	}
	return false
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func InString(need interface{}, haystack interface{}) bool {
	var needStr string
	var needStrArr string
	switch need.(type) {
	case int:
		needStr = strconv.Itoa(need.(int))
	case string:
		needStr = need.(string)
	case int64:
		needStr = strconv.FormatInt(need.(int64), 10)
	case float64:
		needStr = strconv.Itoa(int(need.(float64)))
	}
	switch haystack.(type) {
	case int:
		needStrArr = strconv.Itoa(haystack.(int))
	case string:
		needStrArr = haystack.(string)
	case int64:
		needStrArr = strconv.FormatInt(haystack.(int64), 10)
	case float64:
		needStrArr = strconv.Itoa(int(haystack.(float64)))
	}
	return InArray(needStr, strings.Split(needStrArr, ","))
}

//字符串到数组
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func StrToArray(str string, sep string) interface{} {
	if str == "" {
		return ""
	}
	return strings.Split(str, sep)
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Equal(val1 interface{}, val2 interface{}) bool {
	switch val1.(type) {
	case int:
		switch val2.(type) {
		case int:
			return val1.(int) == val2.(int)
		case string:
			val, _ := strconv.Atoi(val2.(string))
			return val1.(int) == val
		case int64:
			return val1.(int) == int(val2.(int64))
		case float64:
			return val1.(int) == int(val2.(float64))
		}
	case string:
		switch val2.(type) {
		case int:
			return val1.(string) == strconv.Itoa(val2.(int))
		case string:
			return val1.(string) == val2.(string)
		case int64:
			return val1.(string) == strconv.FormatInt(val2.(int64), 10)
		case float64:
			return val1.(string) == strconv.Itoa(int(val2.(float64)))
		}
	case int64:
		switch val2.(type) {
		case int:
			return val1.(int64) == int64(val2.(int))
		case string:
			val, _ := strconv.Atoi(val2.(string))
			return val1.(int64) == int64(val)
		case int64:
			return val1.(int64) == val2.(int64)
		case float64:
			return val1.(int64) == int64(val2.(float64))
		}
	case float64:
		switch val2.(type) {
		case int:
			return val1.(float64) == float64(val2.(int))
		case string:
			val, _ := strconv.ParseFloat(val2.(string), 64)
			return val1.(float64) == val
		case int64:
			return val1.(float64) == float64(val2.(int64))
		case float64:
			return val1.(float64) == val2.(float64)
		}
	}
	return val1 == val2
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetInt(id interface{}) int {
	switch id.(type) {
	case int:
		return id.(int)
	case string:
		ids, _ := strconv.Atoi(id.(string))
		return ids
	case int64:
		return int(id.(int64))
	case float64:
		return int(id.(float64))
	default:
		return id.(int)
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetString(id interface{}) string {
	switch id.(type) {
	case int:
		return strconv.Itoa(id.(int))
	case string:
		return id.(string)
	case int64:
		return strconv.FormatInt(id.(int64), 10)
	case float64:
		return strconv.Itoa(int(id.(float64)))
	default:
		return id.(string)
	}
}

//基本运算
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Operation(val1 interface{}, types string, val2 interface{}) interface{} {
	val1Int := GetInt(val1)
	val2Int := GetInt(val2)
	switch types {
	case "+":
		return val1Int + val2Int
	case "-":
		return val1Int - val2Int
	case "*":
		return val1Int * val2Int
	case "/":
		return val1Int / val2Int
	}
	return val1Int + val2Int
}

//解析字符串
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func AnalysisStr(str string, sep string, args ...interface{}) interface{} {
	if str == "" {
		return ""
	}
	strArr := strings.Split(str, sep)
	strLen := len(strArr)
	var arrList []map[string]interface{}
	for _, v1 := range strArr {
		arrInfo := make(map[string]interface{})
		strs := strings.Split(v1, ":")
		arrInfo["key"] = strs[0]
		arrInfo["val"] = strs[1]
		arrInfo["len"] = strLen
		arrList = append(arrList, arrInfo)
	}
	if len(args) > 0 {
		for _, v := range arrList {
			if v["key"] == args[0] {
				return v["val"]
			}
		}
		return ""
	} else {
		return arrList
	}
}

// 1字符 2数字 3文本 4单选 5复选 6select 7富文本 8图片 9文件 10数组
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetConfigType(type_id interface{}) interface{} {
	typesId := GetInt(type_id)
	types := map[int]string{1: "字符", 2: "数字", 3: "文本", 4: "单选", 5: "复选", 6: "枚举", 7: "富文本", 8: "图片", 9: "文件", 10: "数组", 11: "多图", 12: "日期", 13: "时间", 14: "小时", 15: "颜色"}
	if typesId > 0 {
		return types[typesId]
	}
	return types
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetConfigTypeField(type_id interface{}) interface{} {
	typesId := GetInt(type_id)
	types := map[int]string{1: "varchar(255) NOT NULL", 2: "varchar(10) NOT NULL", 3: "text NOT NULL", 4: "varchar(10) NOT NULL", 5: "varchar(200) NOT NULL", 6: "varchar(10) NOT NULL", 7: "text NOT NULL", 8: "varchar(200) NOT NULL", 9: "varchar(255) NOT NULL", 10: "text NOT NULL", 11: "text NOT NULL", 12: "varchar(50) NOT NULL", 13: "varchar(50) NOT NULL", 14: "varchar(50) NOT NULL", 15: "varchar(50) NOT NULL"}
	if typesId > 0 {
		return types[typesId]
	}
	return types
}

//格式化 post参数
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Inputs(params interface{}) map[string]string {
	var obj map[string]interface{}
	configArr := make(map[string][]string)
	param := make(map[string]string)
	json.Unmarshal([]byte(JsonEncode(params)), &obj)
	for k, v := range obj {
		for _, v1 := range v.([]interface{}) {
			configArr[k] = append(configArr[k], v1.(string))
		}
	}
	for k, v := range configArr {
		param[k] = strings.Join(v, ",")
	}
	return param
}

//返回字段信息
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetFields(infos interface{}, field interface{}) interface{} {
	switch infos.(type) {
	case orm.Params:
		infoVal := infos.(orm.Params)
		return infoVal[GetString(field)]
	case map[string]interface{}:
		infoVal := infos.(map[string]interface{})
		return infoVal[GetString(field)]
	case map[string]string:
		infoVal := infos.(map[string]string)
		return infoVal[GetString(field)]
	}
	return ""
}

//数据库存储转义
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func SqlEscape(str string) string {
	str = strings.Replace(str, "\\", "\\\\", -1)
	str = strings.Replace(str, "/", "\\/", -1)
	str = strings.Replace(str, "<", "\\<", -1)
	str = strings.Replace(str, "`", "\\`", -1)
	str = strings.Replace(str, "=", "\\=", -1)
	str = strings.Replace(str, ".", "\\.", -1)
	str = strings.Replace(str, "?", "\\?", -1)
	str = strings.Replace(str, ">", "\\>", -1)
	str = strings.Replace(str, "'", "\\'", -1)
	str = strings.Replace(str, "\"", "\\\"", -1)
	str = strings.Replace(str, "&", "\\&", -1)
	return str
}

//返回url
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Urls(baseUrl map[string]string, url string, args ...interface{}) string {
	if url == "" {
		url = "/" + baseUrl["MODULE_NAME"] + "/" + baseUrl["CONTROLLER_NAME"] + "/" + baseUrl["ACTION_NAME"]
	}
	if url[0:1] == "/" {
		url = url[1:len(url)]
	}
	urls := strings.Split(url, "/")
	switch len(urls) {
	case 1:
		url = "/" + baseUrl["MODULE_NAME"] + "/" + baseUrl["CONTROLLER_NAME"] + "/" + urls[0]
	case 2:
		url = "/" + baseUrl["MODULE_NAME"] + "/" + urls[0] + "/" + urls[1]
	case 3:
		url = "/" + urls[0] + "/" + urls[1] + "/" + urls[2]
	}
	if len(args) > 0 {
		var params []string
		switch (args[0]).(type) {
		case map[string]string:
			for k, v := range (args[0]).(map[string]string) {
				params = append(params, k+"="+v)
			}
		case string:
			params = append(params, (args[0]).(string))
		case []interface{}:
			for _, v := range (args[0]).([]interface{}) {
				for k1, v1 := range v.(map[string]string) {
					params = append(params, k1+"="+v1)
				}
			}
		}
		url = url + "?" + strings.Join(params, "&")
	}
	if Equal(beego.AppConfig.String("urlmode"), 2) && len(url) > 5 && url[0:5] == "/home" {
		url = url[5:len(url)]
	}
	return url
}

//返回树形结构，子分类是Son ,pidField 是父 ID字段 ,Id是主键
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetTree(lists []orm.Params, pidField string, Id string) []orm.Params {
	SonList := make(map[int][]interface{})
	var pidList []orm.Params
	for _, v := range lists {
		pid := int(v[pidField].(int64))
		if pid > 0 {
			SonList[pid] = append(SonList[pid], v)
		} else {
			pidList = append(pidList, v)
		}
	}

	for k, v := range pidList {
		pidList[k]["Son"] = SonList[int(v[Id].(int64))]
	}

	return pidList
}
