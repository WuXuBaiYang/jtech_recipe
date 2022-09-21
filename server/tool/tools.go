package tool

import (
	"crypto/md5"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"regexp"
)

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
var telephoneRegExp = "^1(3\\d|4[5-9]|5[0-35-9]|6[567]|7[0-8]|8\\d|9[0-35-9])\\d{8}$"

// VerifyTelephone 验证手机号时候正确
func VerifyTelephone(telephone string) bool {
	ok, _ := regexp.MatchString(telephoneRegExp, telephone)
	return ok
}

// If if判断
func If[T any](isTrue bool, a, b T) T {
	if isTrue {
		return a
	}
	return b
}
