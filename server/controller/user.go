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
	RecipeCuisineCodes []string `json:"recipeCuisineCodes" binding:"required"`
	RecipeTasteCodes   []string `json:"recipeTasteCodes" binding:"required"`
}

// 用户信息
type userProfile struct {
	ID         string            `json:"id"`
	Level      int64             `json:"level"`
	NickName   string            `json:"nickName" binding:"required,gte=1"`
	Avatar     string            `json:"avatar"`
	Bio        string            `json:"bio"`
	Profession string            `json:"profession"`
	GenderCode string            `json:"genderCode" binding:"required"`
	Birth      *time.Time        `json:"birth" binding:"ltToday"`
	Medals     []model.UserMedal `json:"medals"`
}

// 用户勋章请求
type medalReq struct {
	Logo       string `json:"logo" binding:"required"`
	Name       string `json:"name" binding:"required,gte=2"`
	RarityCode string `json:"rarityCode"  binding:"required"`
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
	if err := db.First(&subUser, subUId).Error; err != nil {
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
	if err := db.First(&subUser, subUId).Error; err != nil {
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
	target := model.User{OrmBase: model.OrmBase{ID: uId}}
	pagination.Total = db.Model(&target).Association("Subscribes").Count()
	if err := db.Model(&target).Offset((pageIndex - 1) * pageSize).Limit(pageSize).
		Association("Subscribes").Find(&pagination.Data); err != nil {
		response.FailDef(c, -1, "数据查询失败")
		return
	}
	response.SuccessDef(c, pagination)
}

// 分页获取用户帖子操作列表
func getUserPostPagination(c *gin.Context, columnName string, errMessage string) {
	// 获取分页参数
	var pagination model.Pagination[model.Post]
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.FailParams(c, err.Error())
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
	target := model.User{OrmBase: model.OrmBase{ID: uId}}
	pagination.Total = db.Model(&target).Association(columnName).Count()
	if err := db.Model(&target).Offset((pageIndex - 1) * pageSize).Limit(pageSize).
		Association(columnName).Find(&pagination.Data); err != nil {
		response.FailDef(c, -1, errMessage)
		return
	}
	fillPostInfo(c, &pagination.Data)
	response.SuccessDef(c, pagination)
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
	var req medalReq
	if err := c.BindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	// 创建并保存到数据库
	db := common.GetDB()
	result := model.UserMedal{
		OrmBase:    createBase(),
		Logo:       req.Logo,
		Name:       req.Name,
		RarityCode: req.RarityCode,
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
	var req medalReq
	if err := c.BindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	medalId := c.Param("medalId")
	if len(medalId) == 0 {
		response.FailParams(c, "勋章id不能为空")
		return
	}
	db := common.GetDB()
	var result model.UserMedal
	if err := db.First(&result, medalId).Error; err != nil {
		response.FailParams(c, "勋章id不能为空")
		return
	}
	// 更新已有数据
	result.Name = req.Name
	result.Logo = req.Logo
	result.RarityCode = req.RarityCode
	if err := db.Save(&result).Error; err != nil {
		response.FailDef(c, -1, "勋章信息保存失败")
		return
	}
	response.SuccessDef(c, result)
}

// 获取用户勋章
func loadUserMedals(uId string, medals *[]model.UserMedal) error {
	db := common.GetDB()
	err := db.Model(&model.User{
		OrmBase: model.OrmBase{ID: uId},
	}).Association("Medals").Find(medals)
	return err
}
