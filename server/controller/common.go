package controller

import (
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
	Password    string `json:"password" binding:"gte=8"`
	Code        string `json:"code" binding:"len=4"`
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
	code := tool.GenSMSCode(4)
	if err := common.SendSMSVerify(phone, code); err != nil {
		response.FailDef(c, -1, "短信发送失败")
		return
	}
	// 写入redis
	rdb := common.GetRDB()
	if err := rdb.Set(c, phone, code,
		common.SMSExpirationTime).Err(); err != nil {
		response.FailDef(c, -1, "短信发送失败")
		return
	}
	// 响应发送成功的结果,调试模式则返回code
	result := tool.If[any](gin.IsDebugging(), code, true)
	response.SuccessDef(c, result)
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
	rdb := common.GetRDB()
	vCode := rdb.Get(c, req.PhoneNumber).Val()
	if vCode != req.Code {
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
		OrmBase:  createBase(),
		NickName: tool.GenInitUserNickName(),
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
	// 存在密码则验证密码，否则验证校验码
	rdb := common.GetRDB()
	if len(req.Password) != 0 {
		if result.Password != req.Password {
			response.FailParams(c, "密码错误")
			return
		}
	} else {
		vCode := rdb.Get(c, req.PhoneNumber).Val()
		if vCode != req.Code {
			response.FailParams(c, "短信验证码校验失败")
			return
		}
	}
	// 构造授权信息并返回
	auth, authErr := createAuthInfo(c, result)
	if authErr != nil {
		response.FailDef(c, -1, "登录失败")
		return
	}
	// 删除使用过的短信验证码
	rdb.Del(c, req.PhoneNumber)
	response.SuccessDef(c, auth)
}

// RefreshToken 刷新过期token
func RefreshToken(c *gin.Context) {
	claims, err := middleware.GetAccessTokenClaim(c)
	if err != nil {
		response.FailParams(c, "授权信息无效")
		return
	}
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
	db.Find(&user, claims.UserId)
	auth, authErr := createAuthInfo(c, user)
	if authErr != nil {
		response.FailDef(c, -1, "授权失败")
		return
	}
	response.SuccessDef(c, auth)
}

// 创建授权信息
func createAuthInfo(c *gin.Context, user model.User) (*authRes, error) {
	token, err := common.
		ReleaseAccessToken(user, middleware.GetPlatform(c))
	if err != nil {
		return nil, err
	}
	refreshToken, rErr := common.
		ReleaseRefreshToken(user, tool.MD5(token))
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
