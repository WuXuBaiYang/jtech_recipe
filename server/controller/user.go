package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
	"time"
)

// 用户详细信息结构体
type userProfile struct {
	ID         string            `json:"id"`
	Level      int64             `json:"level"`
	NickName   string            `json:"nickName"`
	Avatar     string            `json:"avatar"`
	Bio        string            `json:"bio"`
	Profession string            `json:"profession"`
	GenderCode string            `json:"genderCode"`
	Birth      *time.Time        `json:"birth"`
	Medals     []model.UserMedal `json:"medals,omitempty"`
}

// 全部可修改的用户信息
type fullUserProfile struct {
	userProfile

	EvaluateCode       string   `json:"evaluateCode"`
	RecipeCuisineCodes []string `json:"recipeCuisineCodes"`
	RecipeTasteCodes   []string `json:"recipeTasteCodes"`
}

// UpdateUserProfile 修改用户信息
func UpdateUserProfile(c *gin.Context) {
	// 获取请求体参数
	var profile fullUserProfile
	if err := c.BindJSON(&profile); err != nil {
		response.FailParamsDef(c)
		return
	}
	// 校验参数
	nickName := profile.NickName
	gender := profile.GenderCode
	birth := profile.Birth
	evaluateCode := profile.EvaluateCode
	if len(nickName) == 0 {
		response.FailParams(c, "昵称不能为空")
		return
	}
	if len(gender) == 0 {
		response.FailParams(c, "性别不能为空")
		return
	}
	if birth != nil && time.Now().Before(*birth) {
		response.FailParams(c, "生日填写错误")
		return
	}
	if len(evaluateCode) == 0 {
		response.FailParams(c, "个人评价信息有误")
		return
	}
	// 获取到当前用户信息并写入新的信息
	db := common.GetDB()
	user := middleware.GetCurrUser(c)
	user.NickName = profile.NickName
	user.Avatar = profile.Avatar
	user.Bio = profile.Bio
	user.Profession = profile.Profession
	user.GenderCode = profile.GenderCode
	user.Birth = profile.Birth
	user.EvaluateCode = profile.EvaluateCode
	user.RecipeCuisineCodes = profile.RecipeCuisineCodes
	user.RecipeTasteCodes = profile.RecipeTasteCodes
	if err := db.Save(&user).Error; err != nil {
		response.FailDef(c, -1, "用户信息修改失败")
		return
	}
	response.SuccessDef(c, user)
}

// GetUserProfile 获取用户信息
func GetUserProfile(c *gin.Context) {
	db := common.GetDB()
	userId := c.Param("userId")
	if len(userId) == 0 { // 获取个人信息（完整）
		user := middleware.GetCurrUser(c)
		if err := loadUserMedals(user.ID, &user.Medals); err != nil {
			response.FailDef(c, -1, "获取用户信息失败")
			return
		}
		response.SuccessDef(c, user)
		return
	}
	var profile userProfile
	if err := db.Model(&model.User{}).First(&profile, userId).Error; err != nil {
		response.FailDef(c, -1, "用户不存在")
		return
	}
	if err := loadUserMedals(userId, &profile.Medals); err != nil {
		response.FailDef(c, -1, "获取用户信息失败")
		return
	}
	response.SuccessDef(c, profile)
}

// SubscribeUser 订阅用户
func SubscribeUser(c *gin.Context) {
	// 校验数据
	subUserId := c.Param("userId")
	if len(subUserId) == 0 {
		response.FailParams(c, "用户id不能为空")
		return
	}
	var subUser model.User
	db := common.GetDB()
	db.Where("id = ?", subUserId).Find(&subUser)
	if len(subUser.ID) == 0 {
		response.FailParams(c, "用户不存在")
		return
	}
	user := middleware.GetCurrUser(c)
	if subUser.ID == user.ID {
		response.FailParams(c, "不能订阅自己")
		return
	}
	// 添加订阅关系
	if err := db.Model(&user).
		Association("Subscribes").Append(&subUser); err != nil {
		response.FailDef(c, -1, "订阅失败")
		return
	}
	response.SuccessDef(c, true)
}

// UnSubscribeUser 取消订阅用户
func UnSubscribeUser(c *gin.Context) {
	// 校验数据
	subUserId := c.Param("userId")
	if len(subUserId) == 0 {
		response.FailParams(c, "用户id不能为空")
		return
	}
	var subUser model.User
	db := common.GetDB()
	db.Where("id = ?", subUserId).Find(&subUser)
	if len(subUser.ID) == 0 {
		response.FailParams(c, "用户不存在")
		return
	}
	// 移除订阅关系
	if err := db.Model(middleware.GetCurrUser(c)).
		Association("Subscribes").Delete(&subUser); err != nil {
		response.FailDef(c, -1, "取消订阅失败")
		return
	}
	response.SuccessDef(c, true)
}

// GetSubscribePagination 分页获取订阅列表
func GetSubscribePagination(c *gin.Context) {
	// 获取分页参数
	pagination, err := getPaginationParams(c)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 分页查询
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	userId := c.Param("userId")
	if len(userId) == 0 {
		userId = middleware.GetCurrUId(c)
	}
	db := common.GetDB()
	result := model.Pagination[model.SimpleUser]{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}
	target := model.User{OrmBase: model.OrmBase{ID: userId}}
	result.Total = db.Model(&target).Association("Subscribes").Count()
	if err := db.Model(&target).Offset((pageIndex - 1) * pageSize).Limit(pageSize).
		Association("Subscribes").Find(&result.Data); err != nil {
		response.FailDef(c, -1, "数据查询失败")
		return
	}
	response.SuccessDef(c, result)
}

// 分页获取用户帖子操作列表
func getUserPostPagination(c *gin.Context, columnName string, errMessage string) {
	// 获取分页参数
	pagination, err := getPaginationParams(c)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 分页查询
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	userId := c.Param("userId")
	if len(userId) == 0 {
		userId = middleware.GetCurrUId(c)
	}
	db := common.GetDB()
	result := model.Pagination[model.Post]{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}
	target := model.User{OrmBase: model.OrmBase{ID: userId}}
	result.Total = db.Model(&target).Association(columnName).Count()
	if err := db.Model(&target).Offset((pageIndex - 1) * pageSize).Limit(pageSize).
		Association(columnName).Find(&result.Data); err != nil {
		response.FailDef(c, -1, errMessage)
		return
	}
	fillPostInfo(c, &result.Data)
	response.SuccessDef(c, result)
}

// GetUserViewPostPagination 分页获取用户浏览过的帖子列表
func GetUserViewPostPagination(c *gin.Context) {
	getUserPostPagination(c, "ViewPosts", "用户浏览帖子列表获取失败")
}

// GetUserLikePostPagination 分页获取用户点赞过的帖子列表
func GetUserLikePostPagination(c *gin.Context) {
	getUserPostPagination(c, "LikePosts", "用户点赞帖子列表获取失败")
}

// GetUserCollectPagination 分页获取用户收藏的帖子列表
func GetUserCollectPagination(c *gin.Context) {
	getUserPostPagination(c, "CollectPosts", "用户收藏帖子列表获取失败")
}

// GetAllUserMedalList 获取全部勋章列表
func GetAllUserMedalList(c *gin.Context) {
	db := common.GetDB()
	var result []model.UserMedal
	db.Find(&result)
	response.SuccessDef(c, result)
}

// AddUserMedal 添加勋章
func AddUserMedal(c *gin.Context) {
	// 获取请求参数
	var medal model.UserMedal
	if err := c.BindJSON(&medal); err != nil {
		response.FailParamsDef(c)
		return
	}
	// 数据校验
	logo := medal.Logo
	name := medal.Name
	rarityCode := medal.RarityCode
	if len(logo) == 0 {
		response.FailParams(c, "图标不能为空")
		return
	}
	if len(name) == 0 {
		response.FailParams(c, "名称不能为空")
		return
	}
	if len(rarityCode) == 0 {
		response.FailParams(c, "稀有度不能为空")
		return
	}
	// 创建并保存到数据库
	db := common.GetDB()
	result := model.UserMedal{
		OrmBase:    createBase(),
		Logo:       logo,
		Name:       name,
		RarityCode: rarityCode,
	}
	if err := db.Create(&result).Error; err != nil {
		response.FailDef(c, -1, "勋章创建失败")
		return
	}
	response.SuccessDef(c, result)
}

// UpdateUserMedal 更新勋章信息
func UpdateUserMedal(c *gin.Context) {
	// 获取请求参数
	var medal model.UserMedal
	if err := c.BindJSON(&medal); err != nil {
		response.FailParamsDef(c)
		return
	}
	// 数据校验
	logo := medal.Logo
	name := medal.Name
	rarityCode := medal.RarityCode
	medalId := c.Param("medalId")
	db := common.GetDB()
	var result model.UserMedal
	db.Find(&result, medalId)
	if len(result.ID) == 0 {
		response.FailParams(c, "勋章信息不存在")
		return
	}
	if len(logo) == 0 {
		response.FailParams(c, "图标不能为空")
		return
	}
	if len(name) == 0 {
		response.FailParams(c, "名称不能为空")
		return
	}
	if len(rarityCode) == 0 {
		response.FailParams(c, "稀有度不能为空")
		return
	}
	if len(medalId) == 0 {
		response.FailParams(c, "勋章id不能为空")
		return
	}
	// 更新已有数据
	result.Name = name
	result.Logo = logo
	result.RarityCode = rarityCode
	if err := db.Save(&result).Error; err != nil {
		response.FailDef(c, -1, "勋章信息保存失败")
		return
	}
	response.SuccessDef(c, result)
}

// 获取用户勋章
func loadUserMedals(userId string, medals *[]model.UserMedal) error {
	db := common.GetDB()
	err := db.Model(&model.User{
		OrmBase: model.OrmBase{ID: userId},
	}).Association("Medals").Find(medals)
	return err
}
