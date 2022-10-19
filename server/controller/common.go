package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
	"server/tool"
	"time"
)

// 授权请求
type authReq struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,phone,gte=11"`
	Password    string `json:"password" binding:"omitempty,gte=8"`
	Code        string `json:"code" binding:"omitempty,len=4"`
}

// 授权响应
type authRes struct {
	AccessToken  string     `json:"accessToken"`
	RefreshToken string     `json:"refreshToken"`
	User         model.User `json:"user"`
}

// GetSMS 获取短信验证码
func GetSMS(c *gin.Context) {
	// 获取请求参数并校验手机号
	phone := c.Param("phone")
	if !tool.VerifyPhoneNumber(phone) {
		response.FailParams(c, "手机号校验失败")
		return
	}
	// 发送短信验证码
	code := tool.If(gin.IsDebugging(), phone[len(phone)-4:], tool.GenSMSCode(4))
	if err := common.SendSMSVerify(phone, code); err != nil {
		response.FailDef(c, -1, "短信发送失败")
		return
	}
	// 写入redis
	rdb := common.GetBaseRDB()
	if err := rdb.Set(c, phone, code,
		common.SMSExpirationTime).Err(); err != nil {
		response.FailDef(c, -1, "短信发送失败")
		return
	}
	// 响应发送成功的结果,调试模式则返回code
	response.SuccessDef(c, true)
}

// Auth 用户请求授权
func Auth(c *gin.Context) {
	// 获取请求体参数
	var req authReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	// 校验短信验证码
	rdb := common.GetBaseRDB()
	vCode := rdb.Get(c, req.PhoneNumber)
	if vCode.Err() != nil || vCode.Val() != req.Code {
		response.FailParams(c, "短信验证码校验失败")
		return
	}
	// 判断用户（手机号）是否已存在
	var result model.User
	db := common.GetDB()
	if err := db.Where("phone_number = ?", req.PhoneNumber).
		First(&result).Error; err != nil {
		// 用户不存在则创建用户
		result.OrmBase = createBase()
		result.PhoneNumber = req.PhoneNumber
		result.NickName = tool.GenInitUserNickName()
		if err := db.Create(&result).Error; err != nil {
			response.FailDef(c, -1, "用户创建失败")
			return
		}
	}
	// 构造授权信息并返回
	auth, authErr := createAuthInfo(c, result)
	if authErr != nil {
		response.FailDef(c, -1, "授权失败")
		return
	}
	// 删除使用过的短信验证码
	rdb.Del(c, req.PhoneNumber)
	// 写入用户授权信息
	response.SuccessDef(c, auth)
}

// Register 用户注册接口
func Register(c *gin.Context) {
	// 获取请求体参数
	var req authReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	// 校验短信验证码
	rdb := common.GetBaseRDB()
	vCode := rdb.Get(c, req.PhoneNumber)
	if vCode.Err() != nil || vCode.Val() != req.Code {
		response.FailParams(c, "短信验证码校验失败")
		return
	}
	// 判断用户（手机号）是否已存在
	db := common.GetDB()
	if err := db.Where("phone_number = ?", req.PhoneNumber).
		First(&model.User{}).Error; err == nil {
		response.FailDef(c, -1, "用户已存在")
		return
	}
	// 插入用户信息
	result := model.User{
		OrmBase:     createBase(),
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
		NickName:    tool.GenInitUserNickName(),
	}
	if err := db.Create(&result).Error; err != nil {
		response.FailDef(c, -1, "用户创建失败")
		return
	}
	// 构造授权信息并返回
	auth, authErr := createAuthInfo(c, result)
	if authErr != nil {
		response.FailDef(c, -1, "注册失败")
		return
	}
	// 删除使用过的短信验证码
	rdb.Del(c, req.PhoneNumber)
	// 写入用户授权信息
	response.SuccessDef(c, auth)
}

// Login 用户登录接口
func Login(c *gin.Context) {
	// 获取请求体参数
	var req authReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	// 判断用户是否存在
	db := common.GetDB()
	var result model.User
	if err := db.Where("phone_number = ?", req.PhoneNumber).
		First(&result).Error; err != nil {
		response.FailParams(c, "用户不存在")
		return
	}
	// 检查该账户是否已被封锁
	if result.Blocked {
		response.FailAuth(c, "该账号已被封锁")
		return
	}
	// 存在密码则验证密码，否则验证校验码
	smsRDB := common.GetBaseRDB()
	if len(req.Code) != 0 {
		vCode := smsRDB.Get(c, req.PhoneNumber)
		if vCode.Err() != nil || vCode.Val() != req.Code {
			response.FailParams(c, "短信验证码校验失败")
			return
		}
	} else if len(req.Password) != 0 {
		if result.Password != req.Password {
			response.FailParams(c, "密码错误")
			return
		}
	} else {
		response.FailParams(c, "至少使用密码/验证码进行登录")
		return
	}
	// 构造授权信息并返回
	auth, authErr := createAuthInfo(c, result)
	if authErr != nil {
		response.FailDef(c, -1, "登录失败")
		return
	}
	// 删除使用过的短信验证码
	smsRDB.Del(c, req.PhoneNumber)
	response.SuccessDef(c, auth)
}

// RefreshToken 刷新过期token
func RefreshToken(c *gin.Context) {
	rClaims, rErr := middleware.GetRefreshTokenClaim(c)
	if rErr != nil {
		response.FailParams(c, "授权信息无效")
		return
	}
	// 校验参数
	if rClaims.Target != tool.MD5(middleware.GetAccessToken(c)[7:]) {
		response.FailParams(c, "token不匹配")
		return
	}
	// 重新生成授权信息并返回
	db := common.GetDB()
	var user model.User
	db.Find(&user, rClaims.UserId)
	auth, authErr := createAuthInfo(c, user)
	if authErr != nil {
		response.FailDef(c, -1, "授权失败")
		return
	}
	response.SuccessDef(c, auth)
}

// ForcedOffline 用户强制下线
func ForcedOffline(c *gin.Context) {
	// 获取请求体
	var req = struct {
		UserList []string `json:"userList" binding:"required,gte=1"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 删除授权凭据
	if cmd := common.ClearRDBToken(c,
		req.UserList...); cmd.Err() != nil {
		response.FailDef(c, -1, "强制下线失败")
		return
	}
	response.SuccessDef(c, true)
}

// BlockOut 用户封锁
func BlockOut(c *gin.Context) {
	// 获取请求体
	var req = struct {
		UserList []string `json:"userList" binding:"required,gte=1"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 写入封锁记录
	db := common.GetDB()
	if err := db.Model(&model.User{}).
		Where("id in ?", req.UserList).
		UpdateColumn("blocked", true).
		Error; err != nil {
		response.FailDef(c, -1, "状态更新失败")
		return
	}
	// 清除被封锁的token
	if cmd := common.ClearRDBToken(c,
		req.UserList...); cmd.Err() != nil {
		response.FailDef(c, -1, "授权清除失败")
		return
	}
	response.SuccessDef(c, true)
}

// UnBlockOut 解除用户封锁
func UnBlockOut(c *gin.Context) {
	// 获取请求体
	var req = struct {
		UserList []string `json:"userList" binding:"required,gte=1"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 写入封锁记录
	db := common.GetDB()
	if err := db.Model(&model.User{}).
		Where("id in ?", req.UserList).
		UpdateColumn("blocked", false).
		Error; err != nil {
		response.FailDef(c, -1, "状态更新失败")
		return
	}
	response.SuccessDef(c, true)
}

// 创建授权信息
func createAuthInfo(c *gin.Context, user model.User) (*authRes, error) {
	token, err := common.ReleaseAccessToken(c, user)
	if err != nil {
		return nil, err
	}
	refreshToken, rErr := common.
		ReleaseRefreshToken(c, user, tool.MD5(token))
	if rErr != nil {
		return nil, rErr
	}
	return &authRes{
		AccessToken:  token,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// 创建基础结构体
func createBase() model.OrmBase {
	return model.OrmBase{
		ID:        tool.GenID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// 创建创建者信息
func createCreator(c *gin.Context) model.Creator {
	return model.Creator{
		CreatorId: middleware.GetCurrUId(c),
	}
}

// 检查活动是否符合条件
func checkActivityRecord(activityRecordId *string, activityType model.ActivityType) (*model.ActivityRecord, error) {
	db := common.GetDB()
	var record model.ActivityRecord
	if activityRecordId != nil && len(*activityRecordId) != 0 {
		if err := db.Where("end_time > ?", time.Now()).
			Preload("Activity").
			First(&record, activityRecordId).
			Error; err != nil {
			return nil, errors.New("活动不存在/已结束")
		}
		if !tool.IsContain(record.Activity.TypeCodes, string(activityType)) {
			return nil, errors.New("活动类型不允许")
		}
	}
	return &record, nil
}

// 对对象操作
func operateTarget[T interface{}](c *gin.Context, append bool, columnName string, errMessage string) {
	// 获取请求参数
	targetId := c.Param("targetId")
	if len(targetId) == 0 {
		response.FailParams(c, "id不能为空")
		return
	}
	db := common.GetDB()
	var target T
	if hasNoRecord(&target, targetId) {
		response.FailParams(c, "目标不存在")
		return
	}
	// 操作用户列表
	user := model.User{OrmBase: model.OrmBase{
		ID: middleware.GetCurrUId(c)}}
	subDB := db.Model(&target).Association(columnName)
	if append && subDB.Append(&user) != nil {
		response.FailDef(c, -1, errMessage)
		return
	} else if !append && subDB.Delete(&user) != nil {
		response.FailDef(c, -1, errMessage)
		return
	}
	response.SuccessDef(c, true)
}

// OperateLike 点赞操作
func OperateLike[T interface{}](append bool) gin.HandlerFunc {
	columnName := "LikeUsers"
	errMessage := "点赞操作失败"
	return func(c *gin.Context) {
		operateTarget[T](c, append, columnName, errMessage)
	}
}

// OperateCollect 收藏操作
func OperateCollect[T interface{}](append bool) gin.HandlerFunc {
	columnName := "CollectUsers"
	errMessage := "收藏操作失败"
	return func(c *gin.Context) {
		operateTarget[T](c, append, columnName, errMessage)
	}
}

// 检查记录是否不存在
func hasNoRecord(target interface{}, id string) bool {
	var count int64
	db := common.GetDB()
	db.Model(&target).
		Where("id = ?", id).
		Count(&count).First(&target)
	return count == 0
}
