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
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

/**
 http请求
params map[string]string 请求参数

args[0]  请求方式 get,post,json,xml
args[1]  请求头 map[string]string
args[2]  超时时间
*/

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Http(url string, params interface{}, args ...interface{}) string {
	var method string //请求方式 get,post,json
	var paramsList []string
	var paramStr string //请求参数
	timeout := 30       //超时时间
	if len(args) > 2 {
		timeout = args[2].(int)
	}
	timeouts := time.Duration(timeout) * time.Second

	switch params.(type) {
	case string:
		paramStr = params.(string)
	case map[string]string:
		for k, v := range params.(map[string]string) {
			paramsList = append(paramsList, k+"="+v)
		}
		paramStr = strings.Join(paramsList, "&")
	}
	if len(args) > 0 {
		method = args[0].(string)
	}
	client := &http.Client{Timeout: timeouts}

	var req *http.Request
	switch method {
	case "post":
		req, _ = http.NewRequest("POST", url, strings.NewReader(paramStr))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case "json":
		req, _ = http.NewRequest("POST", url, bytes.NewBuffer([]byte(paramStr)))
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
	case "xml":
		req, _ = http.NewRequest("POST", url, bytes.NewBuffer([]byte(paramStr)))
		req.Header.Set("Content-Type", "application/xml;charset=utf-8")
	default:
		if paramStr != "" && len(strings.Split(paramStr, "=")) > 0 {
			if string(url[len(url)-1]) == "?" {
				url += "&" + paramStr
			} else {
				url += "?" + paramStr
			}
		}
		req, _ = http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if len(args) > 1 {
		for k, v := range args[1].(map[string]string) {
			req.Header.Set(k, v)
		}
	}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	//cookies := resp.Cookies() //遍历cookies
	//	fmt.Println("status", resp.Status)
	//	fmt.Println("response:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}
