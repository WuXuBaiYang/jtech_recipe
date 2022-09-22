package main

import (
	"github.com/gin-gonic/gin"
	"server/controller"
	"server/middleware"
)

// CollectRoutes 统一注册路由方法
func CollectRoutes(r *gin.Engine) *gin.Engine {
	//** 根节点，使用api版本区分 **//
	group := r.Group("/api", middleware.CommonMiddleware())
	//** 用户授权相关 **//
	// 发送短信验证码
	group.POST("/sms/:phone", controller.GetSMS)
	// 用户注册
	group.POST("/register", controller.Register)
	// 用户登录
	group.POST("/login", controller.Login)
	// token刷新
	group.POST("/refreshToken", middleware.AuthMiddleware(false), controller.RefreshToken)

	//// 创建用户相关请求组
	//userGroup := group.Group("/user", middleware.AuthMiddleware(true))
	//// 获取用户信息
	//userGroup.GET("/info/:userId", controller.GetUserProfile)
	//// 获取当前用户信息
	//userGroup.GET("/info", controller.GetUserProfile)
	//// 编辑当前用户信息
	//userGroup.PUT("/info", controller.UpdateUserProfile)
	//// 订阅用户
	//userGroup.POST("/subscribe/:userId", controller.SubscribeUser)
	//// 取消订阅用户
	//userGroup.DELETE("/subscribe/:userId", controller.UnSubscribeUser)
	//// 分页获取订阅列表
	//userGroup.GET("/subscribe", controller.GetSubscribePagination)
	//// 分页获取目标用户的订阅列表
	//userGroup.GET("/subscribe/:userId", controller.GetSubscribePagination)
	//// 分页获取用户点赞帖子列表
	//userGroup.GET("/common/like", controller.GetUserLikePostPagination)
	//// 分页获取目标用户点赞帖子列表
	//userGroup.GET("/common/like/:userId", controller.GetUserLikePostPagination)
	//// 分页获取用户浏览帖子列表
	//userGroup.GET("/common/view", controller.GetUserViewPostPagination)
	//// 分页获取目标用户浏览帖子列表
	//userGroup.GET("/common/view/:userId", controller.GetUserViewPostPagination)
	//// 分页获取用户收藏帖子列表
	//userGroup.GET("/common/collect", controller.GetUserCollectPagination)
	//// 分页获取目标用户收藏帖子列表
	//userGroup.GET("/common/collect/:userId", controller.GetUserCollectPagination)
	//
	////** 帖子相关 **//
	//postGroup := group.Group("/post", middleware.AuthMiddleware(true))
	//// 获取帖子分页列表
	//postGroup.GET("", controller.GetPostPagination)
	//// 获取帖子详情信息
	//postGroup.GET("/:postId", controller.GetPostInfo)
	//// 发布帖子
	//postGroup.POST("", controller.PublishPost)
	//// 编辑帖子
	//postGroup.PUT("/:postId", controller.UpdatePost)
	//// 对帖子浏览
	//postGroup.POST("/view/:postId", controller.AddPostView)
	//// 对帖子点赞
	//postGroup.POST("/like/:postId", controller.AddPostLike)
	//// 对帖子取消点赞
	//postGroup.DELETE("/like/:postId", controller.RemovePostLike)
	//// 对帖子收藏
	//postGroup.POST("/collect/:postId", controller.AddPostCollect)
	//// 对帖子取消收藏
	//postGroup.DELETE("/collect/:postId", controller.RemovePostCollect)
	//
	////** 帖子标签相关 **//
	//postTagGroup := postGroup.Group("/tag")
	//// 新增标签
	//postTagGroup.POST("", controller.AddNewPostTag)
	//// 查询所有标签
	//postTagGroup.GET("", controller.GetPostTagPagination)
	//// 查询特定用户的标签
	//postTagGroup.GET("/:userId", controller.GetPostTagPagination)
	//
	////** 帖子评论/回复相关 **//
	//postCommentGroup := postGroup.Group("/comment")
	//// 发布帖子评论
	//postCommentGroup.POST("/:postId", controller.PublishPostComment)
	//// 分页查询帖子评论和简略（3条）评论回复
	//postCommentGroup.GET("/:postId", controller.GetPostCommentPagination)
	//// 发布帖子评论回复
	//postCommentGroup.POST("/replay/:commentId", controller.PublishPostCommentReplay)
	//// 分页查询帖子评论回复
	//postCommentGroup.GET("/replay/:commentId", controller.GetPostCommentReplayPagination)
	//// 对帖子评论点赞
	//postCommentGroup.POST("/like/:commentId", controller.AddPostCommentLike)
	//// 对帖子评论取消点赞
	//postCommentGroup.DELETE("/like/:commentId", controller.RemovePostCommentLike)
	//// 对帖子评论回复点赞
	//postCommentGroup.POST("/replay/like/:replayId", controller.AddPostCommentReplayLike)
	//// 对帖子评论回复取消点赞
	//postCommentGroup.DELETE("/replay/like/:replayId", controller.RemovePostCommentReplayLike)
	//
	////** 通知相关 **//
	//notificationGroup := group.Group("/notification", middleware.AuthMiddleware(true))
	//// 分页获取通知列表
	//notificationGroup.GET("", controller.GetNotificationPagination)
	return r
}
