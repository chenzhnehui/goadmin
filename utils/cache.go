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
	"github.com/astaxie/beego/cache"
	"strconv"
	"time"
)

var cc cache.Cache

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func init() {
	InitCache()
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func InitCache() {
	var err error
	cc, err = cache.NewCache("file", `{"CachePath":"./tmp/cache","FileSuffix":".cache","DirectoryLevel":"2","EmbedExpiry":"17200"}`)
	if err != nil {
		beego.Info(err)
	}
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func SetCache(key string, value interface{}, timeout int) error {
	if timeout == 0 {
		timeout, _ = strconv.Atoi(beego.AppConfig.String("cachetimeout"))
	}
	if timeout < 0 {
		timeout = 86400 * 365 * 10 // 10年
	}
	timeouts := time.Duration(timeout) * time.Second
	values, err := json.Marshal(&value)
	err = cc.Put(key, string(values), timeouts)
	return err
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func GetCache(key string, obj *interface{}) {
	str := cc.Get(key)
	json.Unmarshal([]byte(str.(string)), &obj)
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func DelCache(key string) {
	cc.Delete(key)
}
