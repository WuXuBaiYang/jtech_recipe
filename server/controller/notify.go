package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
)

// 发送通知
func sendNotify(c *gin.Context, notify model.Notify, toUserIds []int64) []error {
	return nil
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
		Where("to_user_id", middleware.GetCurrUId(c))
	notifyDB.Count(&result.Total)
	notifyDB.Order("created_at DESC").Preload("FromUser", "ToUser").
		Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&result.Data)
	response.SuccessDef(c, result)
}
