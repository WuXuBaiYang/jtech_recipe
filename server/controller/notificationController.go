package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/model"
)

// 发送通知
func sendNotification(c *gin.Context, notifications ...*model.Notification) []error {
	db := common.GetDB()
	var errList []error
	for _, item := range notifications {
		// 数据校验
		if item.ID != 0 {
			errList = append(errList, errors.New("禁止发送已存在通知"))
		}
		item.OrmModel = model.OrmModel{}
		user, _ := getCurrentUser(c)
		if user == nil {
			errList = append(errList, errors.New("请登陆后再发送通知"))
		}
		item.CreatorModel = model.CreatorModel{
			CreatorID: user.ID,
		}
		if item.TargetUserId == 0 {
			errList = append(errList, errors.New("目标用户不能为空"))
		}
		if item.Type < 0 {
			errList = append(errList, errors.New("通知类型错误"))
		}
		if len(item.Title) == 0 {
			errList = append(errList, errors.New("通知标题不能为空"))
		}
		db.Create(&item)
		//** 待补充：异步发送消息推送到用户 **//
	}
	return nil
}

// GetNotificationPagination 获取通知消息列表
func GetNotificationPagination(c *gin.Context) {
	db := common.GetDB()
	// 获取分页参数
	pagination, err := getPaginationParams(c)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	user, _ := getCurrentUser(c)
	notificationDb := db.Model(&model.Notification{}).Where("target_user_id = ?", user.ID)
	// 分页查询
	var count int64
	notificationDb.Count(&count)
	var notificationList []*model.Notification
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	notificationDb = notificationDb.Order("created_at DESC").Offset((pageIndex - 1) * pageSize).Limit(pageSize)
	notificationDb.Preload("Creator.Profile").Preload("Tags").Find(&notificationList)
	response.SuccessDef(c, model.Pagination{
		PageIndex:   pageIndex,
		PageSize:    pageSize,
		Total:       count,
		CurrentSize: len(notificationList),
		Data:        notificationList,
	})
}
