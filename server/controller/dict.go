package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
	"server/tool"
)

// 字典项请求体
type dictReq struct {
	Tag  string `json:"tag" binding:"required,gte=1,lte=18"`
	Info string `json:"info" binding:"lte=30"`
}

// AddDict 批量添加字典项
func AddDict(dictType model.DictType) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求体
		var req = struct {
			DictList []dictReq `json:"dictList" binding:"required,gte=1,lte=20"`
		}{}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.FailParamsDef(c, err)
			return
		}
		// 插入字典项
		db := common.GetDB()
		var results []model.Dict
		userId := middleware.GetCurrUId(c)
		for _, it := range req.DictList {
			code := tool.MD5(tool.JoinV(userId, tool.GenID()))
			results = append(results, model.Dict{
				OrmBase: createBase(),
				Creator: model.Creator{
					CreatorId: userId,
				},
				SimpleDict: model.SimpleDict{
					Code: code,
					Tag:  it.Tag,
					Info: it.Info,
				},
			})
		}
		if err := db.Table(string(dictType)).
			Save(&results).Error; err != nil {
			response.FailDef(c, -1, "标签插入失败")
			return
		}
		response.SuccessDef(c, results)
	}
}

// GetDictPagination 分页获取字典项
func GetDictPagination(dictType model.DictType) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求参数
		var req = struct {
			model.Pagination[model.SimpleDict]
			UserId string `form:"userId"`
		}{}
		if err := c.ShouldBindQuery(&req); err != nil {
			response.FailParamsDef(c, err)
			return
		}
		// 分页查询
		db := common.GetDB()
		pageIndex := req.PageIndex
		pageSize := req.PageSize
		subDB := db.Table(string(dictType))
		if len(req.UserId) != 0 {
			subDB = subDB.Where("creator_id in ?",
				[]string{"", req.UserId})
		}
		subDB.Count(&req.Total)
		if err := subDB.Offset((pageIndex - 1) * pageSize).
			Limit(pageSize).Find(&req.Data).Error; err != nil {
			response.FailDef(c, -1, "标签查询失败")
			return
		}
		response.SuccessDef(c, req.Pagination)
	}
}

// GetUserAddressDictPagination 分页获取用户地址字典项
func GetUserAddressDictPagination(c *gin.Context) {
	// 获取请求参数
	var req = struct {
		model.Pagination[model.SimpleDict]
	}{}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	// 分页查询
	db := common.GetDB()
	pageIndex := req.PageIndex
	pageSize := req.PageSize
	userId := middleware.GetCurrUId(c)
	subDB := db.Table(string(model.UserAddressTagDict)).
		Where("creator_id in ?", []string{"", userId})
	subDB.Count(&req.Total)
	if err := subDB.Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).Find(&req.Data).Error; err != nil {
		response.FailDef(c, -1, "标签查询失败")
		return
	}
	response.SuccessDef(c, req.Pagination)
}
