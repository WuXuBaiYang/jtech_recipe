package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/model"
)

// 帖子的评论请求体
type commentReq struct {
	PId      string `json:"pId" binding:"required,gt=0"`
	Content  string `json:"content" binding:"required,gt=0"`
	TypeCode string `json:"typeCode" binding:"required,dict=comment_type"`
}

// CreateComment 发布评论
func CreateComment(c *gin.Context) {
	// 获取参数消息体
	var req commentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	if hasNoRecord(commentType(req.TypeCode), req.PId) {
		response.FailParams(c, "评论对象不存在")
		return
	}
	// 数据插入
	db := common.GetDB()
	result := model.Comment{
		OrmBase:  createBase(),
		Creator:  createCreator(c),
		PId:      req.PId,
		TypeCode: req.TypeCode,
		Content:  req.Content,
	}
	if err := db.Create(&result).Error; err != nil {
		response.FailDef(c, -1, "评论发布失败")
		return
	}
	response.SuccessDef(c, result)
}

// GetCommentPagination 分页获取评论列表
func GetCommentPagination(c *gin.Context) {
	var req = struct {
		model.Pagination[*model.Comment]
		PId      string `form:"pId" binding:"required,gt=0"`
		TypeCode string `form:"typeCode" binding:"required,dict=comment_type"`
	}{}
	// 获取请求参数
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	db := common.GetDB()
	if hasNoRecord(commentType(req.TypeCode), req.PId) {
		response.FailParams(c, "目标不存在")
		return
	}
	// 分页查询
	pageIndex := req.PageIndex
	pageSize := req.PageSize
	commentDB := db.Model(&model.Comment{}).
		Where("p_id = ?", &req.PId)
	commentDB.Count(&req.Total)
	if err := commentDB.Preload("Creator").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).Find(&req.Data).Error; err != nil {
		response.FailDef(c, -1, "评论查询失败")
		return
	}
	fillCommentInfo(c, req.Data...)
	response.SuccessDef(c, req.Pagination)
}

// 根据传入的类型获取对应的评论父类
func commentType(v string) interface{} {
	switch v {
	case string(model.RecipeComment):
		return &model.Recipe{}
	case string(model.MenuComment):
		return &model.Menu{}
	case string(model.ActivityComment):
		return &model.Activity{}
	default:
		return &model.Post{}
	}
}

// 填充帖子评论信息
func fillCommentInfo(c *gin.Context, items ...*model.Comment) {
	// 待实现
}
