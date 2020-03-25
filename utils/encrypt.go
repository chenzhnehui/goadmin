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
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
)

/**
md5加密
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

/**
Base64Encode加密
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func Base64Encode(input interface{}) string {
	var strs []byte
	switch input.(type) {
	case []byte:
		strs = input.([]byte)
	default:
		strs = []byte(GetString(input))
	}
	return base64.StdEncoding.EncodeToString(strs)
}

/**
Base64Decode 解码
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func Base64Decode(str string) string {
	result, _ := base64.StdEncoding.DecodeString(str)
	return string(result)
}

//json编码
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func JsonEncode(value interface{}) string {
	values, _ := json.Marshal(&value)
	return string(values)
}

//json解码
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func JsonDecode(jsonStr string, args ...interface{}) interface{} {
	if jsonStr == "" {
		return ""
	}
	var obj interface{}
	json.Unmarshal([]byte(jsonStr), &obj)
	if len(args) > 0 {
		objs := obj.(map[string]interface{})
		return objs[args[0].(string)]
	}
	return obj
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize               //1.截取加密代码 段数
	padText := bytes.Repeat([]byte{byte(padding)}, padding) //2.有余数
	return append(src, padText...)                          //3.添加余数
}

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func Depadding(src []byte) []byte {
	lasteum := int(src[len(src)-1]) //1.取出最后一个元素
	return src[:len(src)-lasteum]   //2.删除和最后一个元素相等长的字节
}

/**
DES加密
keys 秘钥长度 8位
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func DesEncrypt(str, keys string) string {
	src := []byte(str)
	key := []byte(keys)
	//1.创建并返回一个使用DES算法的cipher.Block接口。
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}
	src = Padding(src, block.BlockSize()) //2.对src进行填充
	blockModel := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	blockModel.CryptBlocks(src, src) //4.crypto加密连续块
	return Base64Encode(string(src))
}

/**
DES解密
keys 秘钥长度 8位
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func DesDecrypt(str, keys string) string {
	src := []byte(Base64Decode(str))
	key := []byte(keys)
	//1.创建并返回一个使用DES算法的cipher.Block接口。
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockModel := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	blockModel.CryptBlocks(src, src) //3.解密连续块
	return string(Depadding(src))    //.删除填充数组
}
