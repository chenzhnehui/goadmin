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
	"fmt"
	"strings"
	"time"
)

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

/**
获取日期
 */
func Date(format ...interface{})  string{
	var date string
	if(len(format) > 1){
		date = format[1].(string)
	}else{
		date = Now()
	}
   if(len(format) > 0){
	    dateStr := format[0].(string)
		dateTimeArr := strings.Split(date," ")
		timeArr := strings.Split(dateTimeArr[1],":")
		dateArr := strings.Split(dateTimeArr[0],"-")
	    dateStr = strings.Replace(dateStr,"Y",dateArr[0],-1)
	    dateStr = strings.Replace(dateStr,"M",dateArr[1],-1)
	    dateStr = strings.Replace(dateStr,"d",dateArr[2],-1)
	    dateStr = strings.Replace(dateStr,"H",timeArr[0],-1)
	    dateStr = strings.Replace(dateStr,"i",timeArr[1],-1)
	    dateStr = strings.Replace(dateStr,"s",timeArr[2],-1)
		return dateStr
   }
	return date
}

/**
获取当前时间秒
 */
func Time(millisecond ...interface{}) int64 {
	if (len(millisecond) > 0){
		return time.Now().UnixNano() / 1e6
	}else{
		return time.Now().Unix()
	}
}

/**
秒转时间
 */
func TimeToDate(second int64,format ...interface{}) string {
	t := time.Unix(second,0)
	nows := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	if(len(format) > 0){
		return Date(format[0].(string),nows)
	}else{
		return nows
	}
}

/**
时间转秒
 */
func DateToTime(date string) int64  {
	loc,_:= time.LoadLocation("Local")
	//使用模板在对应时区转化为time.time类型
	theTime,_:= time.ParseInLocation("2006-01-02 15:04:05",date,loc)
	return theTime.Unix()
}