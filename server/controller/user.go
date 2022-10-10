package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
	"time"
)

// 用户信息请求
type userProfileReq struct {
	userProfile

	EvaluateCode       string   `json:"evaluateCode" binding:"required"`
	RecipeCuisineCodes []string `json:"recipeCuisineCodes" binding:"required,unique,dict=recipe_cuisine"`
	RecipeTasteCodes   []string `json:"recipeTasteCodes" binding:"required,unique,dict=recipe_taste"`
}

// 用户信息
type userProfile struct {
	ID         string            `json:"id"`
	Level      int64             `json:"level"`
	NickName   string            `json:"nickName" binding:"required,gte=1"`
	Avatar     string            `json:"avatar"`
	Bio        string            `json:"bio"`
	Profession string            `json:"profession"`
	GenderCode string            `json:"genderCode" binding:"required,dict=user_gender"`
	Birth      *time.Time        `json:"birth" binding:"ltToday"`
	Medals     []model.UserMedal `json:"medals"`
}

// 用户勋章请求
type medalReq struct {
	Logo       string `json:"logo" binding:"required"`
	Name       string `json:"name" binding:"required,gte=2"`
	RarityCode string `json:"rarityCode"  binding:"required,dict=medal_rarity"`
}

// UpdateUserProfile 修改用户信息
func UpdateUserProfile(c *gin.Context) {
	// 获取请求体参数
	var req userProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	// 获取到当前用户信息并写入新的信息
	db := common.GetDB()
	user := middleware.GetCurrUser(c)
	user.NickName = req.NickName
	user.Avatar = req.Avatar
	user.Bio = req.Bio
	user.Profession = req.Profession
	user.GenderCode = req.GenderCode
	user.Birth = req.Birth
	user.EvaluateCode = req.EvaluateCode
	user.RecipeCuisineCodes = req.RecipeCuisineCodes
	user.RecipeTasteCodes = req.RecipeTasteCodes
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
	// 没有指定id则返回当前登录用户完整信息
	if len(userId) == 0 {
		user := middleware.GetCurrUser(c)
		if err := loadUserMedals(user.ID, &user.Medals); err != nil {
			response.FailDef(c, -1, "获取用户信息失败")
			return
		}
		response.SuccessDef(c, user)
		return
	}
	// 指定用户id则根据id查询目标用户的部分信息
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
	// 获取数据并校验
	subUId := c.Param("userId")
	if len(subUId) == 0 {
		response.FailParams(c, "用户id不能为空")
		return
	}
	db := common.GetDB()
	user := middleware.GetCurrUser(c)
	if subUId == user.ID {
		response.FailParams(c, "不能订阅自己")
		return
	}
	var subUser model.User
	if hasNoRecord(&subUser, subUId) {
		response.FailParams(c, "用户不存在")
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
	subUId := c.Param("userId")
	if len(subUId) == 0 {
		response.FailParams(c, "用户id不能为空")
		return
	}
	db := common.GetDB()
	var subUser model.User
	if hasNoRecord(&subUser, subUId) {
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
	var pagination model.Pagination[model.SimpleUser]
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	// 获取用户id
	uId := c.Param("userId")
	if len(uId) == 0 {
		uId = middleware.GetCurrUId(c)
	}
	// 分页查询
	db := common.GetDB()
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	user := model.User{OrmBase: model.OrmBase{ID: uId}}
	pagination.Total = db.Model(&user).Association("Subscribes").Count()
	if err := db.Model(&user).Offset((pageIndex - 1) * pageSize).Limit(pageSize).
		Association("Subscribes").Find(&pagination.Data); err != nil {
		response.FailDef(c, -1, "数据查询失败")
		return
	}
	response.SuccessDef(c, pagination)
}
