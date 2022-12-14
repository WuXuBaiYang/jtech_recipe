package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
)

// 食谱请求体
type recipeReq struct {
	Title                string   `json:"title" binding:"required,gte=4"`
	Desc                 string   `json:"desc" binding:"required,gte=1"`
	Images               []string `json:"images" binding:"required,gte=1"`
	Time                 int64    `json:"time" binding:"required,gte=60000"`
	Rating               float32  `json:"rating" binding:"required,min=0,max=1"`
	Steps                []any    `json:"steps" binding:"required,gte=1"`
	CuisineCodes         []string `json:"cuisineCodes" binding:"unique,dict=recipe_cuisine"`
	TasteCodes           []string `json:"tasteCodes"  binding:"unique,dict=recipe_taste"`
	IngredientsMainCodes []string `json:"ingredientsMainCodes" binding:"required,unique,gte=1,dict=recipe_ingredients_main"`
	IngredientsSubCodes  []string `json:"ingredientsSubCodes" binding:"unique,dict=recipe_ingredients_sub"`
	TagCodes             []string `json:"tagCodes" binding:"unique,dict=recipe_tag"`
	ActivityRecordId     *string  `json:"activityRecordId"`
}

// CreateRecipe 创建菜谱
func CreateRecipe(c *gin.Context) {
	// 接收请求体
	var req recipeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	record, err := checkActivityRecord(req.ActivityRecordId,
		model.RecipeActivity)
	if err != nil {
		response.FailParams(c, err.Error())
		return
	}
	// 数据插入
	db := common.GetDB()
	result := model.Recipe{
		OrmBase:              createBase(),
		Creator:              createCreator(c),
		Title:                req.Title,
		Desc:                 req.Desc,
		Images:               req.Images,
		Time:                 req.Time,
		Rating:               req.Rating,
		Steps:                req.Steps,
		CuisineCodes:         req.CuisineCodes,
		TasteCodes:           req.TasteCodes,
		IngredientsMainCodes: req.IngredientsMainCodes,
		IngredientsSubCodes:  req.IngredientsSubCodes,
		TagCodes:             req.TagCodes,
		ActivityRecordId:     req.ActivityRecordId,
	}
	if err := db.Create(&result).Error; err != nil {
		response.FailDef(c, -1, "食谱创建失败")
		return
	}
	result.ActivityRecord = record
	response.SuccessDef(c, result)
}

// UpdateRecipe 编辑食谱
func UpdateRecipe(c *gin.Context) {
	// 接收请求体
	var req recipeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	recipeId := c.Param("recipeId")
	if len(recipeId) == 0 {
		response.FailParams(c, "食谱id不能为空")
		return
	}
	db := common.GetDB()
	var result model.Recipe
	if err := db.Preload("ActivityRecord").
		First(&result, recipeId).Error; err != nil {
		response.FailParams(c, "食谱不存在")
		return
	}
	if result.CreatorId != middleware.GetCurrUId(c) {
		response.FailParams(c, "您不是该食谱的所有者")
		return
	}
	// 数据插入
	result.Title = req.Title
	result.Title = req.Title
	result.Desc = req.Desc
	result.Images = req.Images
	result.Time = req.Time
	result.Rating = req.Rating
	result.Steps = req.Steps
	result.CuisineCodes = req.CuisineCodes
	result.TasteCodes = req.TasteCodes
	result.IngredientsMainCodes = req.IngredientsMainCodes
	result.IngredientsSubCodes = req.IngredientsSubCodes
	result.TagCodes = req.TagCodes
	if err := db.Save(&result).Error; err != nil {
		response.FailDef(c, -1, "食谱保存失败")
		return
	}
	fillRecipeInfo(c, &result)
	response.SuccessDef(c, result)
}

// GetRecipePagination 获取食谱分页列表
func GetRecipePagination(c *gin.Context) {
	// 获取分页参数
	var req = struct {
		model.Pagination[*model.Recipe]
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
	recipeDB := db.Model(&model.Recipe{})
	if len(req.UserId) != 0 {
		recipeDB.Where("creator_id = ?", req.UserId)
	}
	recipeDB.Count(&req.Total)
	if err := recipeDB.Preload("Creator").
		Offset((pageIndex - 1) * pageSize).Limit(pageSize).
		Find(&req.Data).Error; err != nil {
		response.FailDef(c, -1, "食谱查询失败")
		return
	}
	fillRecipeInfo(c, req.Data...)
	response.SuccessDef(c, req.Pagination)
}

// GetRecipeInfo 获取食谱详情
func GetRecipeInfo(c *gin.Context) {
	// 获取请求参数
	recipeId := c.Param("recipeId")
	if len(recipeId) == 0 {
		response.FailParams(c, "食谱id不能为空")
		return
	}
	db := common.GetDB()
	var result model.Recipe
	if err := db.Preload("Creator").
		Preload("ActivityRecord").
		First(&result, recipeId).Error; err != nil {
		response.FailParams(c, "菜单不存在")
		return
	}
	// 补充食谱的标签
	db.Table("sys_dict_recipe_tag").
		Where("code in ?", result.TagCodes).
		Find(&result.Tags)
	fillRecipeInfo(c, &result)
	response.SuccessDef(c, result)
}

// 填充食谱信息
func fillRecipeInfo(c *gin.Context, items ...*model.Recipe) {
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
		db.Raw("select count(*) from sys_recipe_like_users where recipe_id = p.id"),
		db.Raw("select count(*) from sys_recipe_like_users where recipe_id = p.id and user_id = ?", userId),
		db.Raw("select count(*) from sys_recipe_collect_users where recipe_id = p.id"),
		db.Raw("select count(*) from sys_recipe_collect_users where recipe_id = p.id and user_id = ?", userId),
		db.Model(&model.Recipe{}), ids).Scan(&operates)
	for i, it := range operates {
		items[i].LikeCount = it.LikeCount
		items[i].Liked = it.Liked
		items[i].CollectCount = it.CollectCount
		items[i].Collected = it.Collected
	}
}
