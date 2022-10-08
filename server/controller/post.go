package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
)

// 帖子请求体
type postReq struct {
	Title            string   `json:"title" binding:"required,gte=6"`
	Contents         []any    `json:"contents" binding:"required,gte=1"`
	TagCodes         []string `json:"tagCodes" binding:"required"`
	ActivityRecordId *string  `json:"activityId"`
	RecipeId         *string  `json:"recipeId"`
}

// PublishPost 发布帖子
func PublishPost(c *gin.Context) {
	// 接收请求的消息体
	var req postReq
	if err := c.BindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	// 数据插入
	db := common.GetDB()
	result := model.Post{
		OrmBase:          createBase(),
		Creator:          createCreator(c),
		Title:            req.Title,
		Contents:         req.Contents,
		TagCodes:         req.TagCodes,
		ActivityRecordId: req.ActivityRecordId,
		RecipeId:         req.RecipeId,
	}
	if err := db.Create(&result).Error; err != nil {
		response.FailDef(c, -1, "帖子创建失败")
		return
	}
	response.SuccessDef(c, result)
}

// UpdatePost 更新帖子信息
func UpdatePost(c *gin.Context) {
	// 获取请求参数
	var req postReq
	if err := c.BindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	postId := c.Param("postId")
	if len(postId) == 0 {
		response.FailParams(c, "帖子id不能为空")
		return
	}
	db := common.GetDB()
	var result model.Post
	if err := db.First(&result, postId).Error; err != nil {
		response.FailParams(c, "帖子不存在")
		return
	}
	// 更新已有数据
	result.Title = req.Title
	result.Contents = req.Contents
	result.TagCodes = req.TagCodes
	result.ActivityRecordId = req.ActivityRecordId
	result.RecipeId = req.RecipeId
	if err := db.Save(&result).Error; err != nil {
		response.FailDef(c, -1, "帖子信息保存失败")
		return
	}
	response.SuccessDef(c, result)
}

// GetPostPagination 获取帖子分页列表
func GetPostPagination(c *gin.Context) {
	// 获取分页参数
	var pagination model.Pagination[model.Post]
	if err := c.ShouldBindJSON(&pagination); err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 分页查询
	db := common.GetDB()
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	db.Model(&model.Post{}).Count(&pagination.Total)
	if err := db.Model(&model.Post{}).Preload("Creator").
		Offset((pageIndex - 1) * pageSize).Limit(pageSize).
		Find(&pagination.Data).Error; err != nil {
		response.FailDef(c, -1, "帖子查询失败")
		return
	}
	fillPostInfo(c, &pagination.Data)
	response.SuccessDef(c, pagination)
}

// GetPostInfo 获取帖子详情
func GetPostInfo(c *gin.Context) {
	// 获取请求参数
	postId := c.Param("postId")
	if len(postId) == 0 {
		response.FailParams(c, "帖子id不能为空")
		return
	}
	db := common.GetDB()
	var result model.Post
	if err := db.Preload("Creator").
		First(&result, postId).Error; err != nil {
		response.FailParams(c, "帖子不存在")
		return
	}
	// 填充数据并返回
	fillPostDetailInfo(c, &result)
	response.SuccessDef(c, result)
}

// OperatePost 对帖子操作（浏览/点赞/取消点赞/收藏/取消收藏）
func operatePost(c *gin.Context, append bool, columnName string, errMessage string) {
	// 获取请求参数
	postId := c.Param("postId")
	if len(postId) == 0 {
		response.FailParams(c, "帖子id不能为空")
		return
	}
	db := common.GetDB()
	var result model.Post
	if err := db.First(&result, postId).Error; err != nil {
		response.FailParams(c, "帖子不存在")
		return
	}
	// 将当前用户添加到关系列表中
	user := middleware.GetCurrUser(c)
	postDB := db.Model(&result).Association(columnName)
	if append && postDB.Append(user) != nil {
		response.FailDef(c, -1, errMessage)
		return
	} else if !append && postDB.Delete(user) != nil {
		response.FailDef(c, -1, errMessage)
		return
	}
	response.SuccessDef(c, true)
}

// AddPostView 对帖子浏览
func AddPostView(c *gin.Context) {
	operatePost(c, true, "ViewUsers", "帖子浏览失败")
}

// AddPostLike 对帖子点赞
func AddPostLike(c *gin.Context) {
	operatePost(c, true, "LikeUsers", "帖子点赞失败")
}

// RemovePostLike 对帖子取消点赞
func RemovePostLike(c *gin.Context) {
	operatePost(c, false, "LikeUsers", "帖子取消点赞失败")
}

// AddPostCollect 对帖子收藏
func AddPostCollect(c *gin.Context) {
	operatePost(c, true, "CollectUsers", "帖子收藏失败")
}

// RemovePostCollect 对帖子取消收藏
func RemovePostCollect(c *gin.Context) {
	operatePost(c, false, "CollectUsers", "帖子取消收藏失败")
}

// 填充帖子信息
func fillPostInfo(c *gin.Context, items *[]model.Post) {
	/// 帖子的标签功能，有待思考实现方式
	//for i, it := range *items {
	//	(*items)[i].Title = it.Title
	//}
}

// 填充帖子详细信息
func fillPostDetailInfo(c *gin.Context, post *model.Post) {
	/// 帖子的标签功能，有待思考实现方式
	//for i, it := range *items {
	//	(*items)[i].Title = it.Title
	//}
}
