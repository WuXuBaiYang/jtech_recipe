package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
)

// PublishPost 发布帖子
func PublishPost(c *gin.Context) {
	// 接收请求的消息体
	var post model.Post
	if err := c.BindJSON(&post); err != nil {
		response.FailParamsDef(c)
		return
	}
	// 校验数据完整性
	title := post.Title
	contents := post.Contents
	if len(title) == 0 {
		response.FailParams(c, "标题不能为空")
		return
	}
	if len(contents) == 0 {
		response.FailParams(c, "内容不能为空")
		return
	}
	// 数据插入
	result := model.Post{
		OrmBase:          createBase(),
		Creator:          createCreator(c),
		Title:            title,
		Contents:         contents,
		TagCodes:         post.TagCodes,
		ActivityRecordId: post.ActivityRecordId,
		RecipeId:         post.RecipeId,
	}
	db := common.GetDB()
	if err := db.Create(&result).Error; err != nil {
		response.FailDef(c, -1, "帖子创建失败")
		return
	}
	response.SuccessDef(c, result)
}

// UpdatePost 更新帖子信息
func UpdatePost(c *gin.Context) {
	// 接收请求的消息体
	var post model.Post
	if err := c.BindJSON(&post); err != nil {
		response.FailParamsDef(c)
		return
	}
	// 校验数据完整性
	title := post.Title
	contents := post.Contents
	postId := c.Param("postId")
	db := common.GetDB()
	var result model.Post
	db.Find(&result, postId)
	if len(result.ID) == 0 {
		response.FailParams(c, "帖子不存在")
		return
	}
	if len(title) == 0 {
		response.FailParams(c, "标题不能为空")
		return
	}
	if len(contents) == 0 {
		response.FailParams(c, "内容不能为空")
		return
	}
	// 更新已有数据
	result.Title = title
	result.Contents = contents
	result.TagCodes = post.TagCodes
	result.ActivityRecordId = post.ActivityRecordId
	result.RecipeId = post.RecipeId
	if err := db.Save(&result).Error; err != nil {
		response.FailDef(c, -1, "帖子信息保存失败")
		return
	}
	response.SuccessDef(c, result)
}

// GetPostPagination 获取帖子分页列表
func GetPostPagination(c *gin.Context) {
	// 获取分页参数
	pagination, err := getPaginationParams(c)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 分页查询
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	db := common.GetDB()
	result := model.Pagination[model.Post]{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}
	postDB := db.Model(&model.Post{})
	postDB.Count(&result.Total)
	if err := postDB.Preload("Creator").Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).Find(&result.Data).Error; err != nil {
		response.FailDef(c, -1, "帖子查询失败")
		return
	}
	fillPostInfo(c, &result.Data)
	response.SuccessDef(c, result)
}

// GetPostInfo 获取帖子详情
func GetPostInfo(c *gin.Context) {
	// 获取帖子id并校验是否存在
	postId := c.Param("postId")
	db := common.GetDB()
	var result model.Post
	db.Preload("Creator").Find(&result, postId)
	if len(result.ID) == 0 {
		response.FailParams(c, "帖子不存在")
		return
	}
	// 填充数据并返回
	results := []model.Post{result}
	fillPostInfo(c, &results)
	response.SuccessDef(c, results[0])
}

// OperatePost 对帖子操作（浏览/点赞/取消点赞/收藏/取消收藏）
func operatePost(c *gin.Context, append bool, columnName string, errMessage string) {
	// 获取并校验数据有效性
	postId := c.Param("postId")
	db := common.GetDB()
	var post model.Post
	db.Find(&post, postId)
	if len(post.ID) == 0 {
		response.FailParams(c, "帖子不存在")
		return
	}
	// 将当前用户添加到点赞列表中
	user := middleware.GetCurrUser(c)
	postDB := db.Model(&post).Association(columnName)
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

// PublishPostComment 发布帖子评论
func PublishPostComment(c *gin.Context) {
	// 获取参数消息体
	var comment model.PostComment
	if err := c.BindJSON(&comment); err != nil {
		response.FailParamsDef(c)
		return
	}
	// 校验数据完整性
	content := comment.Content
	postId := c.Param("postId")
	if len(postId) == 0 {
		response.FailParams(c, "帖子id不能为空")
		return
	}
	if len(content) == 0 {
		response.FailParams(c, "评论内容不能为空")
		return
	}
	// 数据插入
	db := common.GetDB()
	result := model.PostComment{
		OrmBase: createBase(),
		Creator: createCreator(c),
		PId:     postId,
		Content: content,
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
	pagination, err := getPaginationParams(c)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 分页查询
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	postId := c.Param("postId")
	db := common.GetDB()
	result := model.Pagination[model.PostComment]{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}
	commentDB := db.Model(&model.PostComment{}).
		Where("p_id = ?", postId)
	commentDB.Count(&result.Total)
	if err := commentDB.Preload("Creator").Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).Find(&result.Data).Error; err != nil {
		response.FailDef(c, -1, "评论查询失败")
		return
	}
	fillPostCommentInfo(c, &result.Data)
	response.SuccessDef(c, result)
}

// OperatePostComment 对帖子评论操作（点赞/取消点赞）
func OperatePostComment(c *gin.Context, append bool, columnName string, errMessage string) {
	// 获取数据并校验有效性
	commentId := c.Param("commentId")
	db := common.GetDB()
	var comment model.PostComment
	db.First(&comment, commentId)
	if len(comment.ID) == 0 {
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
	// 获取参数消息体
	var replay model.PostCommentReplay
	if err := c.BindJSON(&replay); err != nil {
		response.FailParamsDef(c)
		return
	}
	// 校验数据完整性
	content := replay.Content
	commentId := c.Param("commentId")
	if len(commentId) == 0 {
		response.FailParams(c, "评论id不能为空")
		return
	}
	if len(content) == 0 {
		response.FailParams(c, "回复内容不能为空")
		return
	}
	// 数据插入
	db := common.GetDB()
	result := model.PostCommentReplay{
		OrmBase: createBase(),
		Creator: createCreator(c),
		PId:     commentId,
		Content: content,
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
	pagination, err := getPaginationParams(c)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 分页查询
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	commentId := c.Param("commentId")
	db := common.GetDB()
	result := model.Pagination[model.PostCommentReplay]{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}
	replayDB := db.Model(&model.PostCommentReplay{}).
		Where("p_id = ?", commentId)
	replayDB.Count(&result.Total)
	if err := replayDB.Preload("Creator").Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).Find(&result.Data).Error; err != nil {
		response.FailDef(c, -1, "回复查询失败")
		return
	}
	fillPostCommentReplayInfo(c, &result.Data)
	response.SuccessDef(c, result)
}

// OperatePostCommentReplay 对帖子评论回复操作（点赞/取消点赞）
func OperatePostCommentReplay(c *gin.Context, append bool, columnName string, errMessage string) {
	// 获取数据并校验
	replayId := c.Param("replayId")
	db := common.GetDB()
	var replay model.PostCommentReplay
	db.Find(&replay, replayId)
	if len(replay.ID) == 0 {
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

// 填充帖子信息
func fillPostInfo(c *gin.Context, items *[]model.Post) {
	/// 帖子的标签功能，有待思考实现方式
	//for i, it := range *items {
	//	(*items)[i].Title = it.Title
	//}
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
