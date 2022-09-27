package tool

import (
	"crypto/md5"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"regexp"
	"time"
)

// GenInitUserNickName 根据格式生成初始化用户昵称
func GenInitUserNickName(uId int64) string {
	return fmt.Sprintf("用户_%d_%v", uId, time.Now().UnixMilli())
}

// 缓存雪花算法节点
var genNode *snowflake.Node

// GenID 雪花算法生成id
func GenID() int64 {
	if genNode == nil {
		Node, err := snowflake.NewNode(899)
		if err != nil {
			return 0
		}
		genNode = Node
	}
	return genNode.Generate().Int64()
}

// MD5 计算md5
func MD5(str string) string {
	data := []byte(str)
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash)
}

// 手机号校验正则
var phoneNumberRegExp = "^1(3\\d|4[5-9]|5[0-35-9]|6[567]|7[0-8]|8\\d|9[0-35-9])\\d{8}$"

// VerifyPhoneNumber 验证手机号时候正确
func VerifyPhoneNumber(phone string) bool {
	ok, _ := regexp.MatchString(phoneNumberRegExp, phone)
	return ok
}

// If if判断
func If[T any](isTrue bool, a, b T) T {
	if isTrue {
		return a
	}
	return b
}
