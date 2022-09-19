package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/model"
)

// AddNewPostTag 新增帖子标签
func AddNewPostTag(c *gin.Context) {
	db := common.GetDB()
	var tag model.PostTag
	err := c.BindJSON(&tag)
	if err != nil {
		response.FailParamsDef(c)
		return
	}
	// 校验数据
	name := tag.Name
	if len(name) <= 0 || len(name) > 20 {
		response.FailParams(c, "标签名称超出限制（1~20）")
		return
	}
	db.Where("name = ?", name).First(&tag)
	if tag.ID != 0 {
		response.FailParams(c, "标签已存在，请勿重复创建")
		return
	}
	user, _ := getCurrentUser(c)
	newTag := model.PostTag{
		CreatorModel: model.CreatorModel{
			CreatorID: user.ID,
		},
		Name: name,
	}
	db.Create(&newTag)
	db.Preload("Creator.Profile").Find(&newTag)
	response.SuccessDef(c, newTag)
}

// GetPostTagPagination 查询所有标签
func GetPostTagPagination(c *gin.Context) {
	db := common.GetDB()
	// 获取分页参数
	pagination, err := getPaginationParams(c)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 分页查询
	pageIndex := pagination.PageIndex
	pageSize := pagination.PageSize
	tagDb := db.Model(&model.PostTag{})
	// 如果查询某个用户的标签，则添加用户条件
	userId := c.Param("userId")
	if len(userId) != 0 {
		tagDb.Where("creator_id = ?", userId)
	}
	var count int64
	tagDb.Count(&count)
	var tagList []model.PostTag
	tagDb.Preload("Creator.Profile").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&tagList)
	response.SuccessDef(c, model.Pagination{
		PageIndex:   pageIndex,
		PageSize:    pageSize,
		Total:       count,
		CurrentSize: len(tagList),
		Data:        tagList,
	})
}
