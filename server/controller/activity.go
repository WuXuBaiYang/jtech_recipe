package controller

import (
	"github.com/gin-gonic/gin"
)

// 活动请求
type activityReq struct {
	CycleTime int64    `json:"cycleTime" validate:"required,gte=86400000"`
	Always    bool     `json:"always" validate:"required"`
	Title     string   `json:"title" validate:"required,gte=6"`
	Url       string   `json:"url" validate:"required,url"`
	TypeCodes []string `json:"typeCodes" validate:"required,unique,gte=1"`
}

// PublishActivity 发布活动信息
func PublishActivity(c *gin.Context) {

}

// UpdateActivity 编辑活动信息
func UpdateActivity(c *gin.Context) {

}

// StartActivity 开始一个活动
func StartActivity(c *gin.Context) {
	// 需要检查是否该活动已启动
}

// GetAllActivityList 获取所有活动列表
func GetAllActivityList(c *gin.Context) {

}
