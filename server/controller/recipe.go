package controller

// 食谱请求体
type recipeReq struct {
	Title                string   `json:"title" binding:"required,gte=4"`
	Desc                 string   `json:"desc" binding:"required,gte=1"`
	Images               []string `json:"images" binding:"required,gte=1"`
	Time                 int64    `json:"time" binding:"required,gte=60000"`
	Rating               float32  `json:"rating" binding:"required,min=0,max=1"`
	Steps                []any    `json:"steps" binding:"required,gte=1"`
	CuisineCodes         []string `json:"cuisineCodes" binding:"dict=recipe_cuisine"`
	TasteCodes           []string `json:"tasteCodes"  binding:"dict=recipe_taste"`
	IngredientsMainCodes []string `json:"ingredientsMainCodes" binding:"required,gte=1,dict=recipe_ingredients_main"`
	IngredientsSubCodes  []string `json:"ingredientsSubCodes" binding:"dict=recipe_ingredients_sub"`
	TagCodes             []string `json:"tagCodes" binding:"dict=recipe_tag"`
	ActivityRecordId     *string  `json:"activityRecordId"`
}
