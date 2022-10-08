package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
)

// 回复请求体
type replayReq struct {
	PId     string `json:"pId" binding:"required,gt=0"`
	Content string `json:"content" binding:"required,gt=0"`
}

// PublishReplay 发布回复
func PublishReplay(c *gin.Context) {
	// 获取请求参数
	var req replayReq
	if err := c.BindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	db := common.GetDB()
	if err := db.First(&model.Comment{}, req.PId).
		Error; err != nil {
		response.FailParams(c, "评论不存在")
		return
	}
	// 数据插入
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
		model.Pagination[model.Replay]
		PId string `form:"pId" binding:"required,gt=0"`
	}{}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailParams(c, err.Error())
		return
	}
	db := common.GetDB()
	if err := db.First(&model.Comment{}, req.PId).
		Error; err != nil {
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
	fillReplayInfo(c, &req.Data)
	response.SuccessDef(c, req.Pagination)
}

// OperateReplay 对帖子评论回复操作（点赞/取消点赞）
func OperateReplay(c *gin.Context, append bool, columnName string, errMessage string) {
	// 获取数据并校验
	replayId := c.Param("replayId")
	if len(replayId) == 0 {
		response.FailParams(c, "回复id不存在")
		return
	}
	db := common.GetDB()
	var replay model.Replay
	if err := db.First(&replay, replayId).Error; err != nil {
		response.FailParams(c, "回复不存在")
		return
	}
	// 将当前用户添加到点赞列表中
	user := middleware.GetCurrUser(c)
	replayDB := db.Model(&replay).Association(columnName)
	if append && replayDB.Append(user) != nil {
		response.FailDef(c, -1, errMessage)
		return
	} else if !append && replayDB.Delete(user) != nil {
		response.FailDef(c, -1, errMessage)
		return
	}
	response.SuccessDef(c, true)
}

// AddReplayLike 对回复点赞
func AddReplayLike(c *gin.Context) {
	OperateReplay(c, true, "LikeUsers", "回复点赞失败")
}

// RemoveReplayLike 对回复取消点赞
func RemoveReplayLike(c *gin.Context) {
	OperateReplay(c, false, "LikeUsers", "回复取消点赞失败")
}

// 填充回复信息
func fillReplayInfo(c *gin.Context, items *[]model.Replay) {
	//for i, it := range *items {
	//	(*items)[i].Title = it.Title
	//}
}
