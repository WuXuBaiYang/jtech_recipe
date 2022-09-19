package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/model"
)

// PublishPostComment 发布帖子评论
func PublishPostComment(c *gin.Context) {
	db := common.GetDB()
	// 获取参数
	var comment model.PostComment
	err := c.BindJSON(&comment)
	if err != nil {
		response.FailParamsDef(c)
		return
	}
	// 校验参数
	content := comment.Content
	if len(content) == 0 {
		response.FailParams(c, "评论内容不能为空")
		return
	}
	postId := c.Param("postId")
	var post model.Post
	db.Where("id = ?", postId).Find(&post)
	if post.ID == 0 {
		response.FailParams(c, "帖子不存在")
		return
	}
	// 评论写入
	user, _ := getCurrentUser(c)
	newComment := &model.PostComment{
		CreatorModel: model.CreatorModel{
			CreatorID: user.ID,
		},
		PostID:  post.ID,
		Content: content,
	}
	db.Create(&newComment)
	db.Preload("Creator.Profile").Find(&newComment)
	response.SuccessDef(c, newComment)
}

// GetPostCommentPagination 分页查询帖子评论以及简略（3条）回复
func GetPostCommentPagination(c *gin.Context) {
	db := common.GetDB()
	// 获取分页参数
	pagination, err := getPaginationParams(c)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 分页查询
	postId := c.Param("postId")
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	commentDb := db.Where("post_id = ?", postId).Model(&model.PostComment{})
	var count int64
	commentDb.Count(&count)
	var commentList []*model.PostComment
	commentDb = commentDb.Offset((pageIndex - 1) * pageSize).Limit(pageSize)
	commentDb = commentDb.Preload("Creator.Profile").Preload("Replays.Creator.Profile")
	commentDb.Find(&commentList)
	fillPostCommentInfo(c, commentList...)
	response.SuccessDef(c, model.Pagination{
		PageIndex:   pageIndex,
		PageSize:    pageSize,
		Total:       count,
		CurrentSize: len(commentList),
		Data:        commentList,
	})
}

// PublishPostCommentReplay 发布帖子评论回复
func PublishPostCommentReplay(c *gin.Context) {
	db := common.GetDB()
	// 获取参数
	var replay model.PostCommentReplay
	err := c.BindJSON(&replay)
	if err != nil {
		response.FailParamsDef(c)
		return
	}
	// 校验参数
	content := replay.Content
	if len(content) == 0 {
		response.FailParams(c, "回复内容不能为空")
		return
	}
	commentId := c.Param("commentId")
	var comment model.PostComment
	db.Where("id = ?", commentId).Find(&comment)
	if comment.ID == 0 {
		response.FailParams(c, "评论不存在")
		return
	}
	// 回复写入
	user, _ := getCurrentUser(c)
	newReplay := &model.PostCommentReplay{
		CreatorModel: model.CreatorModel{
			CreatorID: user.ID,
		},
		CommentID: comment.ID,
		Content:   content,
	}
	db.Create(&newReplay)
	db.Preload("Creator.Profile").Find(&newReplay)
	response.SuccessDef(c, newReplay)
}

// GetPostCommentReplayPagination 分页查询帖子评论回复列表
func GetPostCommentReplayPagination(c *gin.Context) {
	db := common.GetDB()
	// 获取分页参数
	pagination, err := getPaginationParams(c)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 分页查询
	commentId := c.Param("commentId")
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	replayDb := db.Where("comment_id = ?", commentId).Model(&model.PostCommentReplay{})
	var count int64
	replayDb.Count(&count)
	var replayList []*model.PostCommentReplay
	replayDb = replayDb.Offset((pageIndex - 1) * pageSize).Limit(pageSize)
	replayDb = replayDb.Preload("Creator.Profile").Find(&replayList)
	fillPostCommentReplayInfo(c, replayList...)
	response.SuccessDef(c, model.Pagination{
		PageIndex:   pageIndex,
		PageSize:    pageSize,
		Total:       count,
		CurrentSize: len(replayList),
		Data:        replayList,
	})
}

// OperatePostComment 对帖子评论操作（点赞/取消点赞）
func OperatePostComment(c *gin.Context, append bool, columnName string, errMessage string) {
	db := common.GetDB()
	user, _ := getCurrentUser(c)
	commentId := c.Param("commentId")
	// 校验数据
	comment := &model.PostComment{}
	db.Where("id = ?", commentId).First(&comment)
	if comment.ID == 0 {
		response.FailParams(c, "评论不存在")
		return
	}
	// 将当前用户添加到点赞列表中
	likeDb := db.Model(&comment).Association(columnName)
	if append && likeDb.Append(user) != nil {
		response.FailDef(c, -1, errMessage)
		return
	} else if !append && likeDb.Delete(user) != nil {
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

// OperatePostCommentReplay 对帖子评论回复操作（点赞/取消点赞）
func OperatePostCommentReplay(c *gin.Context, append bool, columnName string, errMessage string) {
	db := common.GetDB()
	user, _ := getCurrentUser(c)
	replayId := c.Param("replayId")
	// 校验数据
	replay := &model.PostCommentReplay{}
	db.Where("id = ?", replayId).Find(&replay)
	if replay.ID == 0 {
		response.FailParams(c, "回复不存在")
		return
	}
	// 将当前用户添加到点赞列表中
	likeDb := db.Model(&replay).Association(columnName)
	if append && likeDb.Append(user) != nil {
		response.FailDef(c, -1, errMessage)
		return
	} else if !append && likeDb.Delete(user) != nil {
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

// 填充帖子评论详细信息(点赞)
func fillPostCommentInfo(c *gin.Context, items ...*model.PostComment) {
	db := common.GetDB()
	/// ******** 以下是垃圾代码，找到方法后立即替换 ********///
	user, _ := getCurrentUser(c)
	userId := user.ID
	for _, item := range items {
		// 添加点赞数据
		itemDb := db.Model(&item).Association("LikeUsers")
		item.LikeCount = itemDb.Count()
		user = &model.User{}
		if err := itemDb.Find(&user, userId); err == nil {
			item.Liked = user.ID != 0
		}
		// 添加评论总评论数
		itemDb = db.Model(&item).Association("Replays")
		item.ReplayCount = itemDb.Count()
	}
}

// 填充帖子评论回复详细信息（点赞）
func fillPostCommentReplayInfo(c *gin.Context, items ...*model.PostCommentReplay) {
	db := common.GetDB()
	/// ******** 以下是垃圾代码，找到方法后立即替换 ********///
	user, _ := getCurrentUser(c)
	userId := user.ID
	for _, item := range items {
		// 添加点赞数据
		itemDb := db.Model(&item).Association("LikeUsers")
		item.LikeCount = itemDb.Count()
		user = &model.User{}
		if err := itemDb.Find(&user, userId); err == nil {
			item.Liked = user.ID != 0
		}
	}
}
