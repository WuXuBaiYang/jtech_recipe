package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server/common"
	"server/controller/response"
	"server/model"
	"server/tool"
	"strconv"
)

// Register 用户注册接口
func Register(c *gin.Context) {
	db := common.GetDB()
	// 获取请求体参数
	var userAuth model.UserAuth
	err := c.BindJSON(&userAuth)
	if err != nil {
		response.FailParamsDef(c)
		return
	}
	username := userAuth.UserName
	password := userAuth.Password
	// 数据校验
	if len(username) == 0 {
		response.FailParams(c, "用户名不能为空")
		return
	}
	if len(username) < 8 || len(username) > 24 {
		response.FailParams(c, "用户名长度超出限制（8-24）")
		return
	}
	if len(password) == 0 {
		response.FailParams(c, "密码不能为空")
		return
	}
	user := model.User{
		UserName: username,
		Password: password,
		UserIM: model.UserIM{
			IMUserId: tool.MD5(username),
		},
	}
	db.Where("user_name = ?", username).First(&user)
	if user.ID != 0 {
		response.FailDef(c, -1, "用户已存在")
		return
	}
	// 在事件中执行注册，im注册，im注册失败则回滚
	platform, _ := getPlatform(c)
	err = db.Transaction(func(tx *gorm.DB) error {
		// 创建用户信息
		tx.Create(&user)
		tx.Create(&model.UserProfile{
			CreatorID: user.ID,
			NickName:  username,
		})
		// 请求并创建im账号
		resBody, err := common.RegisterOnIM(*tool.Platform2Int(*platform), user.IMUserId)
		if err != nil {
			return err
		}
		user.IMToken = resBody.Data.Token
		user.IMExpired = resBody.Data.ExpiredTime
		return nil
	})
	if err != nil {
		response.FailDef(c, -1, "用户创建失败")
		return
	}
	// 构造授权信息并返回
	auth, errAuth := createAuthInfo(user, *platform)
	if errAuth != nil {
		response.FailDef(c, -1, "授权失败")
		return
	}
	response.SuccessDef(c, auth)
}

// Login 用户登录接口
func Login(c *gin.Context) {
	db := common.GetDB()
	//获取参数
	var userAuth model.UserAuth
	err := c.BindJSON(&userAuth)
	if err != nil {
		response.FailParamsDef(c)
		return
	}
	username := userAuth.UserName
	password := userAuth.Password
	//数据验证
	if len(username) == 0 {
		response.FailParams(c, "用户名不能为空")
		return
	}
	if len(password) == 0 {
		response.FailParams(c, "密码不能为空")
		return
	}
	user := model.User{
		UserName: username,
		Password: password,
	}
	db.Where("user_name = ?", username).First(&user)
	if user.ID == 0 {
		response.FailParams(c, "用户不存在")
		return
	}
	if user.Password != password {
		response.FailParams(c, "密码错误")
		return
	}
	// 请求im授权信息并赋值
	platform, _ := getPlatform(c)
	resBody, err := common.GetIMUserToken(*tool.Platform2Int(*platform), user.IMUserId)
	if err != nil {
		response.FailDef(c, -1, "IM登录失败")
		return
	}
	user.IMToken = resBody.Data.Token
	user.IMExpired = resBody.Data.ExpiredTime
	// 构造授权信息并返回
	auth, errAuth := createAuthInfo(user, *platform)
	if errAuth != nil {
		response.FailDef(c, -1, "授权失败")
		return
	}
	response.SuccessDef(c, auth)
}

// RefreshToken 刷新过期token
func RefreshToken(c *gin.Context) {
	// 获取请求参数
	accessToken := c.GetHeader("Authorization")[7:]
	refreshToken := c.GetHeader("RefreshToken")
	if len(refreshToken) == 0 {
		response.FailParams(c, "header中缺少refreshToken")
		return
	}
	// 校验参数
	_, refreshClaims, parseErr := common.ParseRefreshToken(refreshToken)
	if parseErr != nil {
		response.FailParams(c, "refreshToken非法")
		return
	}
	if refreshClaims.Target != tool.MD5(accessToken) {
		response.FailParams(c, "token不匹配")
		return
	}
	// 重新生成授权信息并返回
	user, _ := getCurrentUser(c)
	platform, _ := getPlatform(c)
	resBody, err := common.GetIMUserToken(*tool.Platform2Int(*platform), user.IMUserId)
	if err != nil {
		response.FailDef(c, -1, "IM登录失败")
		return
	}
	// 赋值结果并返回
	user.IMToken = resBody.Data.Token
	user.IMExpired = resBody.Data.ExpiredTime
	newAuth, errAuth := createAuthInfo(*user, *platform)
	if errAuth != nil {
		response.FailDef(c, -1, "授权失败")
		return
	}
	response.SuccessDef(c, newAuth)
}

// 创建授权信息
func createAuthInfo(user model.User, platform string) (model.AuthWithProfile, error) {
	token, err := common.ReleaseAccessToken(user, platform)
	if err != nil {
		return model.AuthWithProfile{}, err
	}
	target := tool.MD5(token)
	refreshToken, refreshErr := common.ReleaseRefreshToken(target)
	if refreshErr != nil {
		return model.AuthWithProfile{}, refreshErr
	}
	return model.AuthWithProfile{
		AccessToken:  token,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// 获取分页请求字段
func getPaginationParams(c *gin.Context) (model.Pagination, error) {
	pageIndex, _ := strconv.Atoi(c.Query("pageIndex"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pagination := model.Pagination{
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

// 从上下文中获取到当前访问的用户信息
func getCurrentUser(c *gin.Context) (*model.User, error) {
	user, _ := c.Get("user")
	if u, ok := user.(model.User); ok {
		return &u, nil
	}
	return nil, errors.New("用户不存在")
}

// 获取平台信息
func getPlatform(c *gin.Context) (*string, error) {
	platform, _ := c.Get("platform")
	if v, ok := platform.(string); ok {
		return &v, nil
	}
	return nil, errors.New("平台信息不存在")
}
