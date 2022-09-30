package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
)

// 消息请求
type notifyReq struct {
	ToUsers  []string `json:"toUsers" validate:"required,unique"`
	TypeCode string   `json:"typeCode" validate:"required"`
	Title    string   `json:"title" validate:"required,gte=6"`
	Content  string   `json:"content"`
	Uri      string   `json:"uri"`
}

// 发送通知
func sendNotify(c *gin.Context, notify notifyReq) error {
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
	var req notifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c)
		return
	}
	// 交给方法实现消息推送和写入
	if err := sendNotify(c, req); err != nil {
		response.FailDef(c, -1, "消息发送失败")
		return
	}
	response.SuccessDef(c, true)
}

// GetNotifyPagination 分页获取通知消息列表
func GetNotifyPagination(c *gin.Context) {
	// 获取分页参数
	var pagination model.Pagination[model.Notify]
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.FailParamsDef(c)
		return
	}
	// 分页查询
	db := common.GetDB()
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	notifyDB := db.Model(&model.Notify{}).
		Where("to_user_id in ?", []string{middleware.GetCurrUId(c), ""})
	notifyDB.Count(&pagination.Total)
	notifyDB.Order("created_at DESC").Preload("FromUser").
		Offset((pageIndex - 1) * pageSize).Limit(pageSize).
		Find(&pagination.Data)
	response.SuccessDef(c, pagination)
}
