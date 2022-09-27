package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/model"
)

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
	if subUser.ID == 0 {
		response.FailParams(c, "用户不存在")
		return
	}
	user := getCurrUser(c)
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
	if subUser.ID == 0 {
		response.FailParams(c, "用户不存在")
		return
	}
	// 移除订阅关系
	if err := db.Model(getCurrUser(c)).
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
	userId := parseId(c.Param("userId"))
	if userId == 0 {
		userId = *getCurrUId(c)
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
	userId := parseId(c.Param("userId"))
	if userId == 0 {
		userId = *getCurrUId(c)
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
