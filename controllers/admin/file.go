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
	"goadmin/utils"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type FileController struct {
	BaseController
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func (this *FileController) Upload() {
	jsons := make(map[string]interface{})
	jsons["code"] = 0
	name := this.GetString("filename", "")
	if name == "" {
		name = "imgFile"
	}
	f, h, err := this.GetFile(name)
	if err == nil {
		defer f.Close()
		//生成上传路径
		now := time.Now()
		dir := "./static/uploadfile/" + strconv.Itoa(now.Year()) + "-" + strconv.Itoa(int(now.Month())) + "/" + strconv.Itoa(now.Day())
		err1 := os.MkdirAll(dir, 0755)
		if err1 != nil {
			jsons["msg"] = "目录权限不够"
		} else {
			filename := utils.Md5(time.RFC3339Nano+h.Filename+strconv.Itoa(rand.Intn(int(now.Unix())))) + path.Ext(h.Filename)
			if err != nil {
				jsons["msg"] = err.Error()
			} else {
				this.SaveToFile(name, dir+"/"+filename)
				jsons["code"] = 1
				jsons["url"] = strings.Replace(dir, ".", "", 1) + "/" + filename
				jsons["name"] = h.Filename
				jsons["size"] = h.Size
				jsons["msg"] = "上传成功"
			}
		}
	} else {
		jsons["msg"] = err.Error()
	}

	this.Data["json"] = jsons
	this.ServeJSON()
	return
}
