package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/model"
)

// 帖子的评论请求体
type commentReq struct {
	PId     string `json:"pId" binding:"required,gt=0"`
	Content string `json:"content" binding:"required,gt=0"`
}

// CreateComment 发布评论
func CreateComment[T interface{}](commentType model.CommentType) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取参数消息体
		var req commentReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.FailParamsDef(c, err)
			return
		}
		var target T
		if hasNoRecord(&target, req.PId) {
			response.FailParams(c, "评论对象不存在")
			return
		}
		// 数据插入
		db := common.GetDB()
		result := model.Comment{
			OrmBase:  createBase(),
			Creator:  createCreator(c),
			PId:      req.PId,
			Content:  req.Content,
			TypeCode: commentType,
		}
		if err := db.Create(&result).Error; err != nil {
			response.FailDef(c, -1, "评论发布失败")
			return
		}
		response.SuccessDef(c, result)
	}
}

// GetCommentPagination 分页获取评论列表
func GetCommentPagination[T interface{}](commentType model.CommentType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req = struct {
			model.Pagination[*model.Comment]
			PId string `form:"pId" binding:"required,gt=0"`
		}{}
		// 获取请求参数
		if err := c.ShouldBindQuery(&req); err != nil {
			response.FailParamsDef(c, err)
			return
		}
		db := common.GetDB()
		var target T
		if hasNoRecord(&target, req.PId) {
			response.FailParams(c, "目标不存在")
			return
		}
		// 分页查询
		pageIndex := req.PageIndex
		pageSize := req.PageSize
		commentDB := db.Model(&model.Comment{}).
			Where("p_id = ? and type_code = ?", req.PId, commentType)
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
}

// 填充帖子评论信息
func fillCommentInfo(c *gin.Context, items ...*model.Comment) {
	// 待实现
}
