package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/model"
)

// 用户勋章请求
type medalReq struct {
	Logo       string `json:"logo" binding:"required"`
	Name       string `json:"name" binding:"required,gte=2"`
	RarityCode string `json:"rarityCode"  binding:"required,dict=medal_rarity"`
}

// GetAllUserMedalList 获取全部勋章列表
func GetAllUserMedalList(c *gin.Context) {
	db := common.GetDB()
	var result []model.UserMedal
	db.Find(&result)
	response.SuccessDef(c, result)
}

// AddUserMedal 添加勋章
func AddUserMedal(c *gin.Context) {
	// 获取请求参数
	var req medalReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	// 创建并保存到数据库
	db := common.GetDB()
	result := model.UserMedal{
		OrmBase:    createBase(),
		Logo:       req.Logo,
		Name:       req.Name,
		RarityCode: req.RarityCode,
	}
	if err := db.Create(&result).Error; err != nil {
		response.FailDef(c, -1, "勋章创建失败")
		return
	}
	response.SuccessDef(c, result)
}

// UpdateUserMedal 更新勋章信息
func UpdateUserMedal(c *gin.Context) {
	// 获取请求参数
	var req medalReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	medalId := c.Param("medalId")
	if len(medalId) == 0 {
		response.FailParams(c, "勋章id不能为空")
		return
	}
	db := common.GetDB()
	var result model.UserMedal
	if hasNoRecord(&result, medalId) {
		response.FailParams(c, "勋章id不能为空")
		return
	}
	// 更新已有数据
	result.Name = req.Name
	result.Logo = req.Logo
	result.RarityCode = req.RarityCode
	if err := db.Save(&result).Error; err != nil {
		response.FailDef(c, -1, "勋章信息保存失败")
		return
	}
	response.SuccessDef(c, result)
}

// 获取用户勋章
func loadUserMedals(uId string, medals *[]model.UserMedal) error {
	db := common.GetDB()
	err := db.Model(&model.User{
		OrmBase: model.OrmBase{ID: uId},
	}).Association("Medals").Find(medals)
	return err
}
