package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/model"
)

// 回复请求体
type replayReq struct {
	PId     string `json:"pId" binding:"required,gt=0"`
	Content string `json:"content" binding:"required,gt=0"`
}

// CreateReplay 发布回复
func CreateReplay(c *gin.Context) {
	// 获取请求参数
	var req replayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	if hasNoRecord(&model.Comment{}, req.PId) {
		response.FailParams(c, "评论不存在")
		return
	}
	// 数据插入
	db := common.GetDB()
	result := model.Replay{
		OrmBase: createBase(),
		Creator: createCreator(c),
		PId:     req.PId,
		Content: req.Content,
	}
	if err := db.Create(&result).Error; err != nil {
		response.FailDef(c, -1, "评论回复发布失败")
		return
	}
	response.SuccessDef(c, result)
}

// GetReplayPagination 分页获取回复列表
func GetReplayPagination(c *gin.Context) {
	// 获取请求参数
	var req = struct {
		model.Pagination[*model.Replay]
		PId string `form:"pId" binding:"required,gt=0"`
	}{}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailParams(c, err.Error())
		return
	}
	db := common.GetDB()
	if hasNoRecord(&model.Comment{}, req.PId) {
		response.FailParams(c, "评论不存在")
		return
	}
	// 分页查询
	pageIndex := req.PageIndex
	pageSize := req.PageSize
	replayDB := db.Model(&model.Replay{}).
		Where("p_id = ?", req.PId)
	replayDB.Count(&req.Total)
	if err := replayDB.Preload("Creator").Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).Find(&req.Data).Error; err != nil {
		response.FailDef(c, -1, "回复查询失败")
		return
	}
	fillReplayInfo(c, req.Data...)
	response.SuccessDef(c, req.Pagination)
}

// 填充回复信息
func fillReplayInfo(c *gin.Context, items ...*model.Replay) {
	// 待实现
}
