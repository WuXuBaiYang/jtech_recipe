package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
)

// 菜单请求结构体
type menuReq struct {
	Title            string  `json:"title" binding:"required,gte=1"`
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
	record, err := checkActivityRecord(req.ActivityRecordId,
		model.MenuActivity)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 数据插入
	db := common.GetDB()
	result := model.Menu{
		OrmBase:          createBase(),
		Creator:          createCreator(c),
		Title:            req.Title,
		Contents:         req.Contents,
		ActivityRecordId: req.ActivityRecordId,
	}
	if err := db.Create(&result).Error; err != nil {
		response.FailDef(c, -1, "菜单创建失败")
		return
	}
	result.ActivityRecord = record
	response.SuccessDef(c, result)
}

// ForkMenu 创建分支菜单（复制其他菜单并创建）
func ForkMenu(c *gin.Context) {
	// 获取请求数据
	menuId := c.Param("menuId")
	if len(menuId) == 0 {
		response.FailParams(c, "菜单id不能为空")
		return
	}
	db := common.GetDB()
	var result model.Menu
	if hasNoRecord(&result, menuId) {
		response.FailParams(c, "菜单不存在")
		return
	}
	originMenu := result
	// 数据插入
	result.OrmBase = createBase()
	result.Creator = createCreator(c)
	result.OriginId = &menuId
	result.OriginMenu = &originMenu
	result.ActivityRecordId = nil
	if err := db.Save(&result).Error; err != nil {
		response.FailDef(c, -1, "菜单创建失败")
		return
	}
	fillMenuInfo(c, &result)
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
	if err := db.Preload("ActivityRecord").
		Preload("OriginMenu").
		First(&result, menuId).Error; err != nil {
		response.FailParams(c, "菜单不存在")
		return
	}
	if result.CreatorId != middleware.GetCurrUId(c) {
		response.FailParams(c, "您不是该菜单的所有者")
		return
	}
	// 数据插入
	result.Title = req.Title
	result.Contents = req.Contents
	if err := db.Save(&result).Error; err != nil {
		response.FailDef(c, -1, "菜单保存失败")
		return
	}
	fillMenuInfo(c, &result)
	response.SuccessDef(c, result)
}

// GetMenuPagination 获取菜单分页列表
func GetMenuPagination(c *gin.Context) {
	// 获取分页参数
	var req = struct {
		model.Pagination[*model.Menu]
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
	fillMenuInfo(c, req.Data...)
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
		Preload("ActivityRecord").
		Preload("OriginMenu").
		First(&result, menuId).Error; err != nil {
		response.FailParams(c, "菜单不存在")
		return
	}
	// 填充菜单标签字典
	db.Table("sys_dict_menu_tag").
		Where("code in ?", result.TagCodes).
		Find(&result.Tags)
	fillMenuInfo(c, &result)
	response.SuccessDef(c, result)
}

// 填充菜单信息
func fillMenuInfo(c *gin.Context, items ...*model.Menu) {
	userId := middleware.GetCurrUId(c)
	db := common.GetDB()
	var ids []string
	for _, it := range items {
		ids = append(ids, it.ID)
	}
	var operates []struct {
		Liked        bool
		LikeCount    int64
		Collected    bool
		CollectCount int64
	}
	db.Raw("select (?) as 'LikeCount',(?) as 'Liked',(?) as 'CollectCount',(?) as 'Collected' from (?) as p where p.id in ?",
		db.Raw("select count(*) from sys_menu_like_users where menu_id = p.id"),
		db.Raw("select count(*) from sys_menu_like_users where menu_id = p.id and user_id = ?", userId),
		db.Raw("select count(*) from sys_menu_collect_users where menu_id = p.id"),
		db.Raw("select count(*) from sys_menu_collect_users where menu_id = p.id and user_id = ?", userId),
		db.Model(&model.Menu{}), ids).Scan(&operates)
	for i, it := range operates {
		items[i].LikeCount = it.LikeCount
		items[i].Liked = it.Liked
		items[i].CollectCount = it.CollectCount
		items[i].Collected = it.Collected
	}
}
