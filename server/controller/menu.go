package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/model"
)

// 菜单请求结构体
type menuReq struct {
	Contents         []any   `json:"contents" binding:"required,gte=1"`
	ActivityRecordId *string `json:"activityRecordId"`
}

// CreateMenu 创建菜单
func CreateMenu(c *gin.Context) {
	// 接收请求体
	var req menuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	if err := checkActivityRecord(req.ActivityRecordId,
		model.MenuActivity); err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 数据插入
	db := common.GetDB()
	result := model.Menu{
		OrmBase:          createBase(),
		Creator:          createCreator(c),
		Contents:         req.Contents,
		ActivityRecordId: req.ActivityRecordId,
	}
	if err := db.Create(&result).Error; err != nil {
		response.FailDef(c, -1, "菜单创建失败")
		return
	}
	response.SuccessDef(c, result)
}

// UpdateMenu 编辑菜单
func UpdateMenu(c *gin.Context) {
	// 接收请求体
	var req menuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	menuId := c.Param("menuId")
	if len(menuId) == 0 {
		response.FailParams(c, "菜单id不能为空")
		return
	}
	db := common.GetDB()
	var result model.Menu
	if err := db.First(&result, menuId).Error; err != nil {
		response.FailParams(c, "菜单不存在")
		return
	}
	// 数据插入
	result.Contents = req.Contents
	if err := db.Save(&result).Error; err != nil {
		response.FailDef(c, -1, "菜单保存失败")
		return
	}
	response.SuccessDef(c, result)
}

// GetMenuPagination 获取菜单分页列表
func GetMenuPagination(c *gin.Context) {
	// 获取分页参数
	var req = struct {
		model.Pagination[model.Menu]
		UserId string `form:"userId"`
	}{}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 分页查询
	db := common.GetDB()
	pageIndex := req.PageIndex
	pageSize := req.PageSize
	menuDB := db.Model(&model.Menu{})
	if len(req.UserId) != 0 {
		menuDB.Where("creator_id = ?", req.UserId)
	}
	menuDB.Count(&req.Total)
	if err := menuDB.Preload("Creator").
		Offset((pageIndex - 1) * pageSize).Limit(pageSize).
		Find(&req.Data).Error; err != nil {
		response.FailDef(c, -1, "菜单查询失败")
		return
	}
	response.SuccessDef(c, req.Pagination)
}

// GetMenuInfo 获取菜单详情
func GetMenuInfo(c *gin.Context) {
	// 获取请求参数
	menuId := c.Param("menuId")
	if len(menuId) == 0 {
		response.FailParams(c, "菜单id不能为空")
		return
	}
	db := common.GetDB()
	var result model.Menu
	if err := db.Preload("Creator").
		First(&result, menuId).Error; err != nil {
		response.FailParams(c, "菜单不存在")
		return
	}
	response.SuccessDef(c, result)
}
