package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
)

// 帖子的评论/回复请求体
type commentAndReplayReq struct {
	Content string `json:"content" binding:"required,gt=0"`
}

// PublishPostComment 发布帖子评论
func PublishPostComment(c *gin.Context) {
	// 获取参数消息体
	var req commentAndReplayReq
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
	if err := db.First(&model.Post{}, postId).Error; err != nil {
		response.FailParams(c, "帖子不存在")
		return
	}
	// 数据插入
	result := model.PostComment{
		OrmBase: createBase(),
		Creator: createCreator(c),
		PId:     postId,
		Content: req.Content,
	}
	if err := db.Create(&result).Error; err != nil {
		response.FailDef(c, -1, "帖子评论发布失败")
		return
	}
	response.SuccessDef(c, result)
}

// GetPostCommentPagination 分页获取帖子评论列表
func GetPostCommentPagination(c *gin.Context) {
	// 获取分页参数
	var pagination model.Pagination[model.PostComment]
	if err := c.ShouldBindJSON(&pagination); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	postId := c.Param("postId")
	if len(postId) == 0 {
		response.FailParams(c, "帖子id不能为空")
		return
	}
	db := common.GetDB()
	if err := db.First(&model.Post{}, postId).Error; err != nil {
		response.FailParams(c, "帖子不存在")
		return
	}
	// 分页查询
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	commentDB := db.Model(&model.PostComment{}).Where("p_id = ?", postId)
	commentDB.Count(&pagination.Total)
	if err := commentDB.Preload("Creator").Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).Find(&pagination.Data).Error; err != nil {
		response.FailDef(c, -1, "评论查询失败")
		return
	}
	fillPostCommentInfo(c, &pagination.Data)
	response.SuccessDef(c, pagination)
}

// OperatePostComment 对帖子评论操作（点赞/取消点赞）
func OperatePostComment(c *gin.Context, append bool, columnName string, errMessage string) {
	// 获取请求参数
	commentId := c.Param("commentId")
	if len(commentId) == 0 {
		response.FailParams(c, "评论id不能为空")
		return
	}
	db := common.GetDB()
	var comment model.PostComment
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

// AddPostCommentLike 对帖子评论点赞
func AddPostCommentLike(c *gin.Context) {
	OperatePostComment(c, true, "LikeUsers", "评论点赞失败")
}

// RemovePostCommentLike 对帖子评论取消点赞
func RemovePostCommentLike(c *gin.Context) {
	OperatePostComment(c, false, "LikeUsers", "评论取消点赞失败")
}

// PublishPostCommentReplay 发布帖子评论的回复
func PublishPostCommentReplay(c *gin.Context) {
	// 获取请求参数
	var req commentAndReplayReq
	if err := c.BindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	commentId := c.Param("commentId")
	if len(commentId) == 0 {
		response.FailParams(c, "评论id不能为空")
		return
	}
	db := common.GetDB()
	if err := db.First(&model.PostComment{}, commentId).
		Error; err != nil {
		response.FailParams(c, "评论不存在")
		return
	}
	// 数据插入
	result := model.PostCommentReplay{
		OrmBase: createBase(),
		Creator: createCreator(c),
		PId:     commentId,
		Content: req.Content,
	}
	if err := db.Create(&result).Error; err != nil {
		response.FailDef(c, -1, "评论回复发布失败")
		return
	}
	response.SuccessDef(c, result)
}

// GetPostCommentReplayPagination 分页获取帖子评论回复列表
func GetPostCommentReplayPagination(c *gin.Context) {
	// 获取分页参数
	var pagination model.Pagination[model.PostCommentReplay]
	if err := c.ShouldBindJSON(&pagination); err != nil {
		response.FailParams(c, err.Error())
		return
	}
	commentId := c.Param("commentId")
	if len(commentId) == 0 {
		response.FailParams(c, "评论id不能为空")
		return
	}
	db := common.GetDB()
	if err := db.First(&model.PostComment{}, commentId).
		Error; err != nil {
		response.FailParams(c, "评论不存在")
		return
	}
	// 分页查询
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	replayDB := db.Model(&model.PostCommentReplay{}).
		Where("p_id = ?", commentId)
	replayDB.Count(&pagination.Total)
	if err := replayDB.Preload("Creator").Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).Find(&pagination.Data).Error; err != nil {
		response.FailDef(c, -1, "回复查询失败")
		return
	}
	fillPostCommentReplayInfo(c, &pagination.Data)
	response.SuccessDef(c, pagination)
}

// OperatePostCommentReplay 对帖子评论回复操作（点赞/取消点赞）
func OperatePostCommentReplay(c *gin.Context, append bool, columnName string, errMessage string) {
	// 获取数据并校验
	replayId := c.Param("replayId")
	if len(replayId) == 0 {
		response.FailParams(c, "评论回复id不存在")
		return
	}
	db := common.GetDB()
	var replay model.PostCommentReplay
	if err := db.First(&replay, replayId).Error; err != nil {
		response.FailParams(c, "回复不存在")
		return
	}
	// 将当前用户添加到点赞列表中
	user := middleware.GetCurrUser(c)
	replayDB := db.Model(&replay).Association(columnName)
	if append && replayDB.Append(user) != nil {
		response.FailDef(c, -1, errMessage)
		return
	} else if !append && replayDB.Delete(user) != nil {
		response.FailDef(c, -1, errMessage)
		return
	}
	response.SuccessDef(c, true)
}

// AddPostCommentReplayLike 对帖子评论回复点赞
func AddPostCommentReplayLike(c *gin.Context) {
	OperatePostCommentReplay(c, true, "LikeUsers", "回复点赞失败")
}

// RemovePostCommentReplayLike 对帖子评论回复取消点赞
func RemovePostCommentReplayLike(c *gin.Context) {
	OperatePostCommentReplay(c, false, "LikeUsers", "回复取消点赞失败")
}

// 填充帖子评论信息
func fillPostCommentInfo(c *gin.Context, items *[]model.PostComment) {
	//for i, it := range *items {
	//	(*items)[i].Title = it.Title
	//}
}

// 填充帖子评论回复信息
func fillPostCommentReplayInfo(c *gin.Context, items *[]model.PostCommentReplay) {
	//for i, it := range *items {
	//	(*items)[i].Title = it.Title
	//}
}
