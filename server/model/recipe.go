package model

// Recipe 食谱信息结构体
type Recipe struct {
	OrmBase
	Creator

	Title                string          `json:"title" gorm:"varchar(40);not null;comment:食谱标题"`
	Desc                 string          `json:"desc" gorm:"varchar(200);comment:食谱描述"`
	Cover                string          `json:"cover" gorm:"comment:封面图"`
	Images               []string        `json:"images" gorm:"type:json;serializer:json;not null;comment:食谱图片集合"`
	Time                 int64           `json:"time" gorm:"comment:预计耗时"`
	Rating               float32         `json:"rating" gorm:"comment:难度评分(0-1,0.5一个单位)"`
	Steps                []any           `json:"steps" gorm:"type:json;serializer:json;not null;comment:操作步骤集合"`
	CuisineCodes         []string        `json:"cuisineCodes" gorm:"type:json;serializer:json;comment:所属菜系字典码集合"`
	TasteCodes           []string        `json:"tasteCodes" gorm:"type:json;serializer:json;comment:口味字典码集合"`
	IngredientsMainCodes []string        `json:"ingredientsMainCodes" gorm:"type:json;serializer:json;not null;comment:主料字典码集合"`
	IngredientsSubCodes  []string        `json:"ingredientsSubCodes" gorm:"type:json;serializer:json;not null;comment:辅料字典码集合"`
	TagCodes             []string        `json:"tagCodes" gorm:"type:json;serializer:json;comment:标签字典码集合"`
	Tags                 []SimpleDict    `json:"tags" gorm:"-"`
	ActivityRecordId     *string         `json:"activityRecordId" gorm:"comment:活动id"`
	ActivityRecord       *ActivityRecord `json:"activityRecord,omitempty" gorm:"foreignKey:ActivityRecordId"`

	// 点赞/收藏过的用户
	LikeUsers    []User `json:"-" gorm:"many2many:recipe_like_users"`
	Liked        bool   `json:"liked" gorm:"-"`
	LikeCount    int64  `json:"likeCount" gorm:"-"`
	CollectUsers []User `json:"-" gorm:"many2many:recipe_collect_users"`
	Collected    bool   `json:"collected" gorm:"-"`
	CollectCount int64  `json:"collectCount" gorm:"-"`
}
