package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
)

// 分页获取用户操作列表(帖子/菜单/食谱)
func getUserOperatePagination[T interface{}](c *gin.Context, columnName string, errMessage string) {
	// 获取分页参数
	var pagination model.Pagination[T]
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 分页查询
	db := common.GetDB()
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	user := model.User{OrmBase: model.OrmBase{ID: middleware.GetCurrUId(c)}}
	pagination.Total = db.Model(&user).Association(columnName).Count()
	if err := db.Model(&user).Offset((pageIndex - 1) * pageSize).Limit(pageSize).
		Association(columnName).Find(&pagination.Data); err != nil {
		response.FailDef(c, -1, errMessage)
		return
	}
	response.SuccessDef(c, pagination)
}

// GetUserLikePostPagination 分页获取帖子点赞列表
func GetUserLikePostPagination(c *gin.Context) {
	getUserOperatePagination[model.Post](c, "LikePosts", "帖子点赞列表获取失败")
}

// GetUserLikeMenuPagination 分页获取菜单点赞列表
func GetUserLikeMenuPagination(c *gin.Context) {
	getUserOperatePagination[model.Menu](c, "LikeMenus", "菜单点赞列表获取失败")
}

// GetUserLikeRecipePagination 分页获取食谱点赞列表
func GetUserLikeRecipePagination(c *gin.Context) {
	getUserOperatePagination[model.Recipe](c, "LikeRecipes", "食谱点赞列表获取失败")
}

// GetUserCollectPostPagination 分页获取帖子收藏列表
func GetUserCollectPostPagination(c *gin.Context) {
	getUserOperatePagination[model.Post](c, "CollectPosts", "帖子收藏列表获取失败")
}

// GetUserCollectMenuPagination 分页获取菜单收藏列表
func GetUserCollectMenuPagination(c *gin.Context) {
	getUserOperatePagination[model.Menu](c, "CollectMenus", "菜单收藏列表获取失败")
}

// GetUserCollectRecipePagination 分页获取食谱收藏列表
func GetUserCollectRecipePagination(c *gin.Context) {
	getUserOperatePagination[model.Recipe](c, "CollectRecipes", "食谱收藏列表获取失败")
}
