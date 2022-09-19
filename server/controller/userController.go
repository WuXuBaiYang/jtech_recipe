package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server/common"
	"server/controller/response"
	"server/model"
	"server/tool"
	"time"
)

// UpdateUserProfile 修改用户信息
func UpdateUserProfile(c *gin.Context) {
	db := common.GetDB()
	// 获取请求体参数
	var profile model.UserProfile
	err := c.BindJSON(&profile)
	if err != nil {
		response.FailParamsDef(c)
		return
	}
	// 校验参数
	nickName := profile.NickName
	gender := profile.Gender
	birthday := profile.Birthday
	telephone := profile.Telephone
	if len(nickName) == 0 {
		response.FailParams(c, "昵称不能为空")
		return
	}
	if gender < 0 || gender > 2 {
		response.FailParams(c, "性别填写错误")
		return
	}
	if time.Now().Before(birthday) {
		response.FailParams(c, "生日超出范围")
		return
	}
	if !tool.VerifyTelephone(telephone) {
		response.FailParams(c, "手机号格式错误")
		return
	}
	// 获取到当前用户信息并写入新的信息
	user, _ := getCurrentUser(c)
	db.Preload("Profile").First(&user)
	user.Profile = &model.UserProfile{
		OrmModel:   user.Profile.OrmModel,
		CreatorID:  user.Profile.CreatorID,
		NickName:   profile.NickName,
		Avatar:     profile.Avatar,
		Telephone:  profile.Telephone,
		Bio:        profile.Bio,
		Address:    profile.Address,
		Location:   profile.Location,
		Profession: profile.Profession,
		Email:      profile.Email,
		Gender:     profile.Gender,
		Birthday:   profile.Birthday,
	}
	// 启用事务处理，如果im更新失败，则回滚状态
	platform, _ := getPlatform(c)
	imToken := c.GetHeader("IMToken")
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&user.Profile).Error; err != nil {
			return err
		}
		user.IMToken = imToken
		if err := common.UpdateIMUserInfo(*tool.Platform2Int(*platform), *user); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		response.FailDef(c, -1, "用户信息修改失败")
		return
	}
	response.SuccessDef(c, user.Profile)
}

// GetUserProfile 获取用户信息
func GetUserProfile(c *gin.Context) {
	db := common.GetDB()
	userId := c.Param("userId")
	if len(userId) != 0 {
		var user model.User
		db.Preload("Profile").First(&user, userId)
		if user.ID == 0 {
			response.FailDef(c, -1, "用户不存在")
			return
		}
		/// 处理黑名单，敏感数据等操作
		response.SuccessDef(c, user)
		return
	}
	user, _ := getCurrentUser(c)
	db.Preload("Profile").First(&user)
	response.SuccessDef(c, user)
}
