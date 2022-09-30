package tool

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"os"
	"regexp"
	"time"
)

// GenInitUserNickName 根据格式生成初始化用户昵称
func GenInitUserNickName() string {
	return fmt.Sprintf("用户_%d", time.Now().UnixMilli())
}

// 缓存雪花算法节点
var genNode *snowflake.Node

// GenID 雪花算法生成id
func GenID() string {
	if genNode == nil {
		Node, err := snowflake.NewNode(899)
		if err != nil {
			return ""
		}
		genNode = Node
	}
	return fmt.Sprintf("%d", genNode.Generate().Int64())
}

// GenSMSCode 生成短信验证码
func GenSMSCode(count int) string {
	code := GenID()
	return code[len(code)-count:]
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

// ReadJsonFile 读取Json文件
func ReadJsonFile(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			println("文件关闭失败")
		}
	}(f)
	if err := json.NewDecoder(f).
		Decode(&v); err != nil {
		return err
	}
	return nil
}
