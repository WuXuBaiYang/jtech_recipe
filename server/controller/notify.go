package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
)

// 接收消息通知发送的结构体
type notifyItem struct {
	ToUsers []string `json:"toUsers"`
	model.Notify
}

// 发送通知
func sendNotify(c *gin.Context, notify notifyItem) error {
	db := common.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		// 创建要插入的数据集合并插入数据库
		var notifies []model.Notify
		fromUId := middleware.GetCurrUId(c)
		if len(notify.ToUsers) == 0 {
			notify.ToUsers = append(notify.ToUsers, "")
		}
		for _, uId := range notify.ToUsers {
			notifies = append(notifies, model.Notify{
				OrmBase:    createBase(),
				FromUserId: fromUId,
				ToUserId:   uId,
				TypeCode:   notify.TypeCode,
				Title:      notify.Title,
				Content:    notify.Content,
				Uri:        notify.Uri,
			})
		}
		if err := tx.CreateInBatches(notifies,
			100).Error; err != nil {
			return err
		}
		// 可能会需要执行消息推送行为
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// PushNotify 推送一条通知
func PushNotify(c *gin.Context) {
	// 获取请求参数体
	var notify notifyItem
	if err := c.BindJSON(&notify); err != nil {
		response.FailParamsDef(c)
		return
	}
	// 校验参数
	title := notify.Title
	content := notify.Content
	typeCode := notify.TypeCode
	if len(title) == 0 {
		response.FailParams(c, "标题不能为空")
		return
	}
	if len(content) == 0 {
		response.FailParams(c, "内容不能为空")
		return
	}
	if len(typeCode) == 0 {
		response.FailParams(c, "消息类型不能为空")
		return
	}
	// 交给方法实现消息推送和写入
	if err := sendNotify(c, notify); err != nil {
		response.FailDef(c, -1, "消息发送失败")
		return
	}
	response.SuccessDef(c, true)
}

// GetNotifyPagination 分页获取通知消息列表
func GetNotifyPagination(c *gin.Context) {
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
	result := model.Pagination[model.Notify]{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}
	notifyDB := db.Model(&model.Notify{}).
		Where("to_user_id in ?", []string{middleware.GetCurrUId(c), ""})
	notifyDB.Count(&result.Total)
	notifyDB.Order("created_at DESC").Preload("FromUser").
		Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&result.Data)
	response.SuccessDef(c, result)
}
