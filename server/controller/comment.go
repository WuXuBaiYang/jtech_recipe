package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
	"server/tool"
)

// 帖子的评论请求体
type commentReq struct {
	PId     string `json:"pId" binding:"required,gt=0"`
	Content string `json:"content" binding:"required,gt=0"`
	Type    string `json:"type" binding:"required,type=comment"`
}

// PublishComment 发布评论
func PublishComment(c *gin.Context) {
	// 获取参数消息体
	var req commentReq
	if err := c.BindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	db := common.GetDB()
	if err := db.First(tool.CommentType(req.Type), req.PId).
		Error; err != nil {
		response.FailParams(c, "评论目标不存在")
		return
	}
	// 数据插入
	result := model.Comment{
		OrmBase: createBase(),
		Creator: createCreator(c),
		PId:     req.PId,
		Type:    req.Type,
		Content: req.Content,
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
		model.Pagination[model.Comment]
		PId  string `form:"pId" binding:"required,gt=0"`
		Type string `form:"type" binding:"required,type=comment"`
	}{}
	// 获取请求参数
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	db := common.GetDB()
	if err := db.First(tool.CommentType(req.Type), req.PId).
		Error; err != nil {
		response.FailParams(c, "目标不存在")
		return
	}
	// 分页查询
	pageIndex := req.PageIndex
	pageSize := req.PageSize
	commentDB := db.Model(&model.Comment{}).
		Where("p_id = ?", &req.PId)
	commentDB.Count(&req.Total)
	if err := commentDB.Preload("Creator").Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).Find(&req.Data).Error; err != nil {
		response.FailDef(c, -1, "评论查询失败")
		return
	}
	fillCommentInfo(c, &req.Data)
	response.SuccessDef(c, req.Pagination)
}

// OperateComment 对评论操作（点赞/取消点赞）
func OperateComment(c *gin.Context, append bool, columnName string, errMessage string) {
	// 获取请求参数
	commentId := c.Param("commentId")
	if len(commentId) == 0 {
		response.FailParams(c, "评论id不能为空")
		return
	}
	db := common.GetDB()
	var comment model.Comment
	if err := db.First(&comment, commentId).Error; err != nil {
		response.FailParams(c, "评论不存在")
		return
	}
	// 将当前用户添加到点赞列表中
	user := middleware.GetCurrUser(c)
	commentDB := db.Model(&comment).Association(columnName)
	if append && commentDB.Append(user) != nil {
		response.FailDef(c, -1, errMessage)
		return
	} else if !append && commentDB.Delete(user) != nil {
		response.FailDef(c, -1, errMessage)
		return
	}
	response.SuccessDef(c, true)
}

// AddCommentLike 对评论点赞
func AddCommentLike(c *gin.Context) {
	OperateComment(c, true, "LikeUsers", "评论点赞失败")
}

// RemoveCommentLike 对评论取消点赞
func RemoveCommentLike(c *gin.Context) {
	OperateComment(c, false, "LikeUsers", "评论取消点赞失败")
}

// 填充帖子评论信息
func fillCommentInfo(c *gin.Context, items *[]model.Comment) {
	//for i, it := range *items {
	//	(*items)[i].Title = it.Title
	//}
}
