package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server/common"
	"server/controller/response"
	"server/model"
)

// PublishPost 发布帖子
func PublishPost(c *gin.Context) {
	db := common.GetDB()
	post, err := getAndVerifyPostInfo(c)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 插入到本地，事务处理
	user, _ := getCurrentUser(c)
	db.Preload("Profile").First(&user)
	newPost := &model.Post{
		CreatorModel: model.CreatorModel{
			CreatorID: user.ID,
		},
		Title:    post.Title,
		Contents: post.Contents,
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		tx.Create(&newPost)
		tags := post.Tags
		err := tx.Model(&newPost).Association("Tags").Append(tags)
		if err != nil {
			return err
		}
		tx.Preload("Creator.Profile").Preload("Tags").Find(&newPost)
		return nil
	})
	if err != nil {
		response.FailDef(c, -1, "帖子创建失败")
		return
	}
	// 查找关注当前用户的列表并发送通知消息
	type SubMyUser struct {
		UserId uint
	}
	var subUsers []SubMyUser
	subDb := db.Table("user_subscribe").Where("subscribe_user_id = ?", user.ID)
	subDb = subDb.Select([]string{"user_id"}).Scan(&subUsers)
	title := fmt.Sprintf("你关注的 %s 发布了一篇帖子，快去看看吧", user.Profile.NickName)
	content, _ := json.Marshal(newPost.Contents)
	uri := fmt.Sprintf("[post](postId=%v)", newPost.ID)
	var notifications []*model.Notification
	for _, item := range subUsers {
		notifications = append(notifications, &model.Notification{
			Type:         1,
			TargetUserId: item.UserId,
			Title:        title,
			Content:      string(content),
			Uri:          uri,
		})
	}
	sendNotification(c, notifications...)
	// 返回帖子信息
	response.SuccessDef(c, newPost)
}

// UpdatePost 更新帖子信息
func UpdatePost(c *gin.Context) {
	db := common.GetDB()
	// 获取帖子详情
	post, err := getAndVerifyPostInfo(c)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	postId := c.Param("postId")
	var newPost model.Post
	db.Preload("Creator.Profile").Find(&newPost, postId)
	if newPost.ID == 0 {
		response.FailDef(c, -1, "帖子不存在")
		return
	}
	user, _ := getCurrentUser(c)
	if user.ID != newPost.CreatorID {
		response.FailDef(c, -1, "没有操作该帖子的权限")
		return
	}
	newPost.Title = post.Title
	newPost.Contents = post.Contents
	err = db.Transaction(func(tx *gorm.DB) error {
		tx.Save(&newPost)
		tags := post.Tags
		tx.Find(&tags)
		err := tx.Model(&newPost).Association("Tags").Replace(tags)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		response.FailDef(c, -1, "帖子编辑失败")
		return
	}
	fillPostInfo(c, &newPost)
	response.SuccessDef(c, newPost)
}

// GetPostInfo 获取帖子详情
func GetPostInfo(c *gin.Context) {
	db := common.GetDB()
	postId := c.Param("postId")
	var newPost model.Post
	db.Preload("Creator.Profile").Preload("Tags").Find(&newPost, postId)
	if newPost.ID == 0 {
		response.FailDef(c, -1, "帖子不存在")
		return
	}
	// 加载最近点赞的5个人
	postDb := db.Preload("Profile").Model(&newPost).Limit(5)
	if err := postDb.Association("LikeUsers").Find(&newPost.LikeUserList); err != nil {
		println("最近点赞用户查询失败")
	}
	fillPostInfo(c, &newPost)
	response.SuccessDef(c, newPost)
}

// GetPostPagination 获取帖子分页列表
func GetPostPagination(c *gin.Context) {
	db := common.GetDB()
	// 获取分页参数
	pagination, err := getPaginationParams(c)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	postDb := db.Model(&model.Post{})
	// 查询过滤条件
	/// 待补充 ///
	userId := c.Query("userId")
	query := ""
	if len(userId) != 0 {
		query += "creator_id = " + userId
	}
	postDb = postDb.Where(query)
	// 分页查询
	var count int64
	postDb.Count(&count)
	var postList []*model.Post
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	postDb = postDb.Offset((pageIndex - 1) * pageSize).Limit(pageSize)
	postDb.Preload("Creator.Profile").Preload("Tags").Find(&postList)
	fillPostInfo(c, postList...)
	response.SuccessDef(c, model.Pagination{
		PageIndex:   pageIndex,
		PageSize:    pageSize,
		Total:       count,
		CurrentSize: len(postList),
		Data:        postList,
	})
}

// OperatePost 对帖子操作（浏览/点赞/取消点赞/收藏/取消收藏）
func operatePost(c *gin.Context, append bool, columnName string, errMessage string) {
	db := common.GetDB()
	user, _ := getCurrentUser(c)
	postId := c.Param("postId")
	// 校验数据
	post := &model.Post{}
	db.Where("id = ?", postId).Find(&post)
	if post.ID == 0 {
		response.FailParams(c, "帖子不存在")
		return
	}
	// 将当前用户添加到点赞列表中
	likeDb := db.Model(&post).Association(columnName)
	if append && likeDb.Append(user) != nil {
		response.FailDef(c, -1, errMessage)
		return
	} else if !append && likeDb.Delete(user) != nil {
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

// 获取并校验帖子信息
func getAndVerifyPostInfo(c *gin.Context) (model.Post, error) {
	// 获取请求参数
	var post model.Post
	err := c.BindJSON(&post)
	if err != nil {
		return post, err
	}
	// 数据校验
	title := post.Title
	contents := post.Contents
	if len(title) == 0 {
		return post, errors.New("标题不能为空")
	}
	if len(contents) == 0 {
		return post, errors.New("内容不能为空")
	}
	return post, nil
}

// 填充帖子详细信息(点赞/收藏/浏览等)
func fillPostInfo(c *gin.Context, items ...*model.Post) {
	db := common.GetDB()
	user, _ := getCurrentUser(c)
	userId := user.ID
	/// ******** 以下是垃圾代码，找到方法后立即替换 ********///
	for _, item := range items {
		// 添加浏览数据
		itemDb := db.Model(&item).Association("ViewUsers")
		item.ViewCount = itemDb.Count()
		user = &model.User{}
		if err := itemDb.Find(&user, userId); err == nil {
			item.Viewed = user.ID != 0
		}
		// 添加点赞数据
		itemDb = db.Model(&item).Association("LikeUsers")
		item.LikeCount = itemDb.Count()
		user = &model.User{}
		if err := itemDb.Find(&user, userId); err == nil {
			item.Liked = user.ID != 0
		}
		// 添加收藏数据
		itemDb = db.Model(&item).Association("CollectUsers")
		item.CollectCount = itemDb.Count()
		user = &model.User{}
		if err := itemDb.Find(&user, userId); err == nil {
			item.Collected = user.ID != 0
		}
	}
}
