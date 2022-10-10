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
	TagCodes         []string `json:"tagCodes" binding:"required,unique,dict=post_tag"`
	ActivityRecordId *string  `json:"activityRecordId"`
	RecipeId         *string  `json:"recipeId"`
}

// CreatePost 发布帖子
func CreatePost(c *gin.Context) {
	// 接收请求的消息体
	var req postReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	record, err := checkActivityRecord(req.ActivityRecordId,
		model.PostActivity)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	db := common.GetDB()
	var recipe *model.Recipe
	if req.RecipeId != nil && len(*req.RecipeId) != 0 {
		if hasNoRecord(&recipe, *req.RecipeId) {
			response.FailParams(c, "菜谱不存在")
			return
		}
	}
	// 数据插入
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
	fillPostInfo(c, &result)
	result.ActivityRecord = record
	result.Recipe = recipe
	response.SuccessDef(c, result)
}

// UpdatePost 更新帖子信息
func UpdatePost(c *gin.Context) {
	// 获取请求参数
	var req postReq
	if err := c.ShouldBindJSON(&req); err != nil {
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
	if err := db.Preload("ActivityRecord").
		Preload("Recipe").First(&result, postId).
		Error; err != nil {
		response.FailParams(c, "帖子不存在")
		return
	}
	if result.CreatorId != middleware.GetCurrUId(c) {
		response.FailParams(c, "您不是该帖子的所有者")
		return
	}
	// 更新已有数据
	result.Title = req.Title
	result.Contents = req.Contents
	result.TagCodes = req.TagCodes
	if err := db.Save(&result).Error; err != nil {
		response.FailDef(c, -1, "帖子保存失败")
		return
	}
	fillPostInfo(c, &result)
	response.SuccessDef(c, result)
}

// GetPostPagination 获取帖子分页列表
func GetPostPagination(c *gin.Context) {
	// 获取分页参数
	var req = struct {
		model.Pagination[*model.Post]
		UserId string `form:"userId"`
	}{}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 分页查询
	db := common.GetDB()
	pageIndex := req.PageIndex
	pageSize := req.PageSize
	postDB := db.Model(&model.Post{})
	if len(req.UserId) != 0 {
		postDB.Where("creator_id = ?", req.UserId)
	}
	postDB.Count(&req.Total)
	if err := postDB.Preload("Creator").
		Offset((pageIndex - 1) * pageSize).Limit(pageSize).
		Find(&req.Data).Error; err != nil {
		response.FailDef(c, -1, "帖子查询失败")
		return
	}
	fillPostInfo(c, req.Data...)
	response.SuccessDef(c, req.Pagination)
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
		Preload("ActivityRecord").
		Preload("Recipe").
		First(&result, postId).Error; err != nil {
		response.FailParams(c, "帖子不存在")
		return
	}
	// 填充数据并返回
	fillPostInfo(c, &result)
	response.SuccessDef(c, result)
}

// 填充帖子信息
func fillPostInfo(c *gin.Context, items ...*model.Post) {
	/// 待实现
}
