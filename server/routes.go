package main

import (
	"github.com/gin-gonic/gin"
	"server/controller"
	"server/middleware"
	"server/model"
)

// CollectRoutes 统一注册路由方法
func CollectRoutes(r *gin.Engine) *gin.Engine {
	//** 根节点，使用api版本区分 **//
	group := r.Group("/api", middleware.Common)
	//** 授权校验组 **//
	authGroup := group.Group("", middleware.AuthCheck)
	//** 授权相关 **//
	authRoutes(group)
	//** 用户相关 **//
	userRoutes(authGroup.Group("/user"))
	//** 帖子相关 **//
	postRoutes(authGroup.Group("/post"))
	//** 菜单相关 **//
	menuRoutes(authGroup.Group("/menu"))
	//** 食谱相关 **//
	recipeRoutes(authGroup.Group("/recipe"))
	//** 评论相关 **//
	commentRoutes(authGroup.Group("/comment"))
	//** 回复相关 **//
	replayRoutes(authGroup.Group("/replay"))
	//** 活动相关 **//
	activityRoutes(authGroup.Group("/activity"))
	//** 通知相关 **//
	notifyRoutes(authGroup.Group("/notification"))
	return r
}

// 授权相关路由
func authRoutes(group *gin.RouterGroup) {
	// 发送短信验证码
	group.POST("/sms/:phone", controller.GetSMS)
	// 用户注册
	group.POST("/register", controller.Register)
	// 用户登录
	group.POST("/login", controller.Login)
	// token刷新
	group.POST("/refreshToken", controller.RefreshToken)
	// 用户强制下线
	group.POST("/forceOffline", middleware.PermissionCheck, controller.ForcedOffline)
	// 封锁用户
	group.POST("/blockOut", middleware.PermissionCheck, controller.BlockOut)
	// 解除用户封锁
	group.POST("/unBlockOut", middleware.PermissionCheck, controller.UnBlockOut)
}

// 用户相关路由
func userRoutes(group *gin.RouterGroup) {
	// 订阅用户
	group.POST("/subscribe/:userId", controller.SubscribeUser)
	// 取消订阅用户
	group.DELETE("/subscribe/:userId", controller.UnSubscribeUser)
	// 分页获取订阅列表
	group.GET("/subscribe", controller.GetSubscribePagination)
	// 分页获取目标用户的订阅列表
	group.GET("/subscribe/:userId", controller.GetSubscribePagination)
	// 获取用户信息
	group.GET("/info/:userId", controller.GetUserProfile)
	// 获取当前用户信息
	group.GET("/info", controller.GetUserProfile)
	// 编辑当前用户信息
	group.PUT("/info", controller.UpdateUserProfile)
	// 分页获取帖子点赞列表
	group.GET("/post/like", controller.GetUserLikePostPagination)
	// 分页获取菜单点赞列表
	group.GET("/menu/like", controller.GetUserLikeMenuPagination)
	// 分页获取食谱点赞列表
	group.GET("/recipe/like", controller.GetUserLikeRecipePagination)
	// 分页获取帖子收藏列表
	group.GET("/post/collect", controller.GetUserCollectPostPagination)
	// 分页获取菜单收藏列表
	group.GET("/menu/collect", controller.GetUserCollectMenuPagination)
	// 分页获取食谱收藏列表
	group.GET("/recipe/collect", controller.GetUserCollectRecipePagination)
	// 获取全部勋章列表
	group.GET("/medal", controller.GetAllUserMedalList)
	// 添加勋章[权限]
	group.POST("/medal", middleware.PermissionCheck, controller.AddUserMedal)
	// 更新勋章信息[权限]
	group.PUT("/medal/:medalId", middleware.PermissionCheck, controller.UpdateUserMedal)
}

// 帖子相关路由
func postRoutes(group *gin.RouterGroup) {
	// 发布帖子
	group.POST("", controller.CreatePost)
	// 编辑帖子
	group.PUT("/:postId", controller.UpdatePost)
	// 获取帖子分页列表
	group.GET("", controller.GetPostPagination)
	// 获取帖子详情信息
	group.GET("/:postId", controller.GetPostInfo)
	// 对帖子点赞
	group.POST("/like/:targetId", controller.OperateLike[model.Post](true))
	// 对帖子取消点赞
	group.DELETE("/like/:targetId", controller.OperateLike[model.Post](false))
	// 对帖子收藏
	group.POST("/collect/:targetId", controller.OperateCollect[model.Post](true))
	// 对帖子取消收藏
	group.DELETE("/collect/:targetId", controller.OperateCollect[model.Post](false))
}

// 菜单相关路由
func menuRoutes(group *gin.RouterGroup) {
	// 创建一个菜单
	group.POST("", controller.CreateMenu)
	// 创建分支菜单
	group.POST("/fork/:menuId", controller.ForkMenu)
	// 编辑一个菜单
	group.PUT("/:menuId", controller.UpdateMenu)
	// 分页获取菜单列表
	group.GET("", controller.GetMenuPagination)
	// 获取菜单详情
	group.GET("/:menuId", controller.GetMenuInfo)
	// 对菜单点赞
	group.POST("/like/:targetId", controller.OperateLike[model.Menu](true))
	// 对菜单取消点赞
	group.DELETE("/like/:targetId", controller.OperateLike[model.Menu](false))
	// 对菜单收藏
	group.POST("/collect/:targetId", controller.OperateCollect[model.Menu](true))
	// 对菜单取消收藏
	group.DELETE("/collect/:targetId", controller.OperateCollect[model.Menu](false))
}

// 食谱相关路由
func recipeRoutes(group *gin.RouterGroup) {
	// 创建一个食谱
	group.POST("", controller.CreateRecipe)
	// 编辑一个食谱
	group.PUT("/:recipeId", controller.UpdateRecipe)
	// 分页获取食谱列表
	group.GET("", controller.GetRecipePagination)
	// 获取食谱详情
	group.GET("/:recipeId", controller.GetRecipeInfo)
	// 对食谱点赞
	group.POST("/like/:targetId", controller.OperateLike[model.Recipe](true))
	// 对食谱取消点赞
	group.DELETE("/like/:targetId", controller.OperateLike[model.Recipe](false))
	// 对食谱收藏
	group.POST("/collect/:targetId", controller.OperateCollect[model.Recipe](true))
	// 对食谱取消收藏
	group.DELETE("/collect/:targetId", controller.OperateCollect[model.Recipe](false))
}

// 评论相关
func commentRoutes(group *gin.RouterGroup) {
	// 发布评论
	group.POST("", controller.CreateComment)
	// 分页查询评论
	group.GET("", controller.GetCommentPagination)
	// 对评论点赞
	group.POST("/like/:targetId", controller.OperateLike[model.Comment](true))
	// 对评论取消点赞
	group.DELETE("/like/:targetId", controller.OperateLike[model.Comment](false))
}

// 回复相关
func replayRoutes(group *gin.RouterGroup) {
	// 发布评论回复
	group.POST("", controller.CreateReplay)
	// 分页评论回复
	group.GET("", controller.GetReplayPagination)
	// 对评论回复点赞
	group.POST("/like/:targetId", controller.OperateLike[model.Replay](true))
	// 对评论回复取消点赞
	group.DELETE("/like/:targetId", controller.OperateLike[model.Replay](false))
}

// 活动相关路由
func activityRoutes(group *gin.RouterGroup) {
	// 发布一个活动
	group.POST("", middleware.PermissionCheck, controller.CreateActivity)
	// 编辑一个活动
	group.PUT("/:activityId", middleware.PermissionCheck, controller.UpdateActivity)
	// 开始一个活动
	group.POST("/start/:activityId", middleware.PermissionCheck, controller.StartActivity)
	// 获取全部活动列表
	group.GET("", controller.GetAllActivityList)
	// 获取全部进行中的活动列表
	group.GET("/process", controller.GetAllActivityProcessList)
}

// 通知相关路由
func notifyRoutes(group *gin.RouterGroup) {
	// 分页获取通知列表
	group.GET("", controller.GetNotifyPagination)
	// 发送消息通知[权限]
	group.POST("", middleware.PermissionCheck, controller.PushNotify)
}
