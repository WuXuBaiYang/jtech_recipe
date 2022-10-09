package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/model"
	"time"
)

// 活动请求
type activityReq struct {
	CycleTime int64    `json:"cycleTime" binding:"required,gte=86400000"`
	Always    bool     `json:"always" binding:"required"`
	Title     string   `json:"title" binding:"required,gte=6"`
	Url       string   `json:"url" binding:"required,url"`
	TypeCodes []string `json:"typeCodes" binding:"required,unique,gte=1,dict=activity_type"`
}

// 活动记录请求
type activityRecordReq struct {
	BeginTime time.Time `json:"beginTime" binding:"required,gtToday"`
}

// PublishActivity 发布活动信息
func PublishActivity(c *gin.Context) {
	// 获取请求参数
	var req activityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	// 数据插入
	db := common.GetDB()
	result := model.Activity{
		OrmBase:   createBase(),
		CycleTime: req.CycleTime,
		Always:    req.Always,
		Title:     req.Title,
		Url:       req.Url,
		TypeCodes: req.TypeCodes,
	}
	if err := db.Create(&result).Error; err != nil {
		response.FailDef(c, -1, "活动创建失败")
		return
	}
	response.SuccessDef(c, result)
}

// UpdateActivity 编辑活动信息
func UpdateActivity(c *gin.Context) {
	// 获取请求参数
	var req activityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	activityId := c.Param("activityId")
	if len(activityId) == 0 {
		response.FailParams(c, "活动id不能为空")
		return
	}
	db := common.GetDB()
	var result model.Activity
	if err := db.First(&result, activityId).Error; err != nil {
		response.FailParams(c, "活动不存在")
		return
	}
	// 更新已有数据
	result.CycleTime = req.CycleTime
	result.Always = req.Always
	result.Title = req.Title
	result.Url = req.Url
	result.TypeCodes = req.TypeCodes
	if err := db.Save(&result).Error; err != nil {
		response.FailDef(c, -1, "活动信息保存失败")
		return
	}
	response.SuccessDef(c, result)
}

// StartActivity 开始一个活动
func StartActivity(c *gin.Context) {
	// 获取请求参数
	var req activityRecordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	activityId := c.Param("activityId")
	if len(activityId) == 0 {
		response.FailParams(c, "活动id不存在")
		return
	}
	db := common.GetDB()
	var activity model.Activity
	if err := db.First(&activity, activityId).Error; err != nil {
		response.FailParams(c, "活动不存在")
		return
	}
	// 判断活动是否已启动
	if err := db.Where("activity_id = ? and end_time > ?", activityId, time.Now()).
		First(&model.ActivityRecord{}).Error; err == nil {
		response.FailParams(c, "活动已开始，无法重复启动")
		return
	}
	// 整理并写入活动记录
	result := model.ActivityRecord{
		OrmBase:   createBase(),
		BeginTime: req.BeginTime,
		EndTime: req.BeginTime.
			Add(time.Duration(activity.CycleTime)),
		ActivityId: activityId,
	}
	if err := db.Create(&result).Error; err != nil {
		response.FailDef(c, -1, "活动启动失败")
		return
	}
	response.SuccessDef(c, result)
}

// GetAllActivityList 获取所有活动列表
func GetAllActivityList(c *gin.Context) {
	db := common.GetDB()
	var result []model.Activity
	db.Find(&result)
	response.SuccessDef(c, result)
}

// GetAllActivityProcessList 获取所有进行中的活动列表
func GetAllActivityProcessList(c *gin.Context) {
	db := common.GetDB()
	var result []model.ActivityRecord
	db.Where("end_time > ?", time.Now()).Find(&result)
	response.SuccessDef(c, result)
}
