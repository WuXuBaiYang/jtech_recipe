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

// 获取用户勋章
func loadUserMedals(userId string, medals *[]model.UserMedal) error {
	db := common.GetDB()
	err := db.Model(&model.User{
		OrmBase: model.OrmBase{ID: userId},
	}).Association("Medals").Find(medals)
	return err
}
