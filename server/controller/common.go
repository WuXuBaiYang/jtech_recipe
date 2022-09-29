package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
	"server/tool"
	"strconv"
	"time"
)

// 用户登录注册授权信息结构体
type userAuth struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Code        string `json:"code"`
}

// GetSMS 获取短信验证码
func GetSMS(c *gin.Context) {
	// 获取请求参数
	phone := c.Param("phone")
	if !tool.VerifyPhoneNumber(phone) {
		response.FailParams(c, "手机号校验失败")
		return
	}
	// 发送短信验证码并写入redis
	rdb := common.GetRDB()
	code, err := common.SendSMSVerify(phone)
	if err != nil {
		response.FailDef(c, -1, "短信发送失败")
		return
	}
	r := rdb.Set(c, phone, *code, common.SMSExpirationTime)
	if r.Err() != nil {
		response.FailDef(c, -1, "短信发送失败")
		return
	}
	response.SuccessDef(c, true)
}

// Register 用户注册接口
func Register(c *gin.Context) {
	// 获取请求体参数
	reqAuth := &userAuth{}
	if err := c.BindJSON(reqAuth); err != nil {
		response.FailParamsDef(c)
		return
	}
	phoneNumber := reqAuth.PhoneNumber
	password := reqAuth.Password
	code := reqAuth.Code
	rdb := common.GetRDB()
	// 数据校验
	if !tool.VerifyPhoneNumber(phoneNumber) {
		response.FailParams(c, "手机号校验失败")
		return
	}
	vCode := rdb.Get(c, phoneNumber).Val()
	if len(code) == 0 || vCode != code {
		response.FailParams(c, "短信验证码校验失败")
		return
	}
	if len(password) == 0 {
		response.FailParams(c, "密码不能为空")
		return
	}
	user := model.User{
		PhoneNumber: phoneNumber,
		Password:    password,
	}
	db := common.GetDB()
	db.Where("phone_number = ?", phoneNumber).First(&user)
	if len(user.ID) != 0 {
		response.FailDef(c, -1, "用户已存在")
		return
	}
	// 在事件中执行注册
	err := db.Transaction(func(tx *gorm.DB) error {
		// 创建用户信息
		user.OrmBase = createBase()
		user.NickName = tool.GenInitUserNickName(user.ID)
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		response.FailDef(c, -1, "用户创建失败")
		return
	}
	// 构造授权信息并返回
	if auth, err := createAuthInfo(c, user); err == nil {
		response.SuccessDef(c, auth)
		// 删除使用过的短信验证码
		rdb.Del(c, phoneNumber)
		return
	}
	response.FailDef(c, -1, "授权失败")
}

// Login 用户登录接口
func Login(c *gin.Context) {
	// 获取请求体参数
	reqAuth := &userAuth{}
	if err := c.BindJSON(reqAuth); err != nil {
		response.FailParamsDef(c)
		return
	}
	phoneNumber := reqAuth.PhoneNumber
	password := reqAuth.Password
	code := reqAuth.Code
	rdb := common.GetRDB()
	// 数据校验
	if !tool.VerifyPhoneNumber(phoneNumber) {
		response.FailParams(c, "手机号校验失败")
		return
	}
	user := model.User{
		PhoneNumber: phoneNumber,
		Password:    password,
	}
	db := common.GetDB()
	db.Where("phone_number = ?", phoneNumber).First(&user)
	if len(user.ID) == 0 {
		response.FailParams(c, "用户不存在")
		return
	}
	if len(code) != 0 {
		vCode := rdb.Get(c, phoneNumber).Val()
		if vCode != code {
			response.FailParams(c, "短信验证码校验失败")
			return
		}
	} else if len(password) != 0 {
		if user.Password != password {
			response.FailParams(c, "密码错误")
			return
		}
	} else {
		response.FailParams(c, "缺少短信验证码或密码")
		return
	}
	// 构造授权信息并返回
	if auth, err := createAuthInfo(c, user); err == nil {
		response.SuccessDef(c, auth)
		// 删除使用过的短信验证码
		rdb.Del(c, phoneNumber)
		return
	}
	response.FailDef(c, -1, "授权失败")
}

// RefreshToken 刷新过期token
func RefreshToken(c *gin.Context) {
	if claims, err := middleware.GetTokenClaim(c); err == nil {
		if refreshClaims, err := middleware.GetRefreshTokenClaim(c); err == nil {
			// 校验参数
			if refreshClaims.Target != tool.MD5(middleware.GetToken(c)[7:]) {
				response.FailParams(c, "token不匹配")
				return
			}
			// 重新生成授权信息并返回
			db := common.GetDB()
			var user = model.User{}
			db.Find(&user, claims.UserId)
			newAuth, errAuth := createAuthInfo(c, user)
			if errAuth != nil {
				response.FailDef(c, -1, "授权失败")
				return
			}
			response.SuccessDef(c, newAuth)
			return
		}
	}
	// 验证失败
	response.FailParams(c, "授权信息无效")
}

// 创建授权信息
func createAuthInfo(c *gin.Context, user model.User) (*model.AuthWithUser, error) {
	platform := middleware.GetPlatform(c)
	token, err := common.ReleaseAccessToken(user, platform)
	if err != nil {
		return nil, err
	}
	target := tool.MD5(token)
	refreshToken, refreshErr := common.ReleaseRefreshToken(user, target)
	if refreshErr != nil {
		return nil, refreshErr
	}
	return &model.AuthWithUser{
		AccessToken:  token,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// 获取分页请求字段
func getPaginationParams(c *gin.Context) (model.Pagination[any], error) {
	pageIndex, _ := strconv.Atoi(c.Query("pageIndex"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pagination := model.Pagination[any]{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}
	if pageIndex <= 0 {
		return pagination, errors.New("pageIndex不合法（>0）")
	}
	if pageSize <= 0 || pageSize > 100 {
		return pagination, errors.New("pageSize不合法（1~100）")
	}
	return pagination, nil
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
