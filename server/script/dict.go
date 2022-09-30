package script

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"regexp"
	"server/common"
	"server/model"
	"server/tool"
	"strings"
	"time"
)

// 自定义字典表
type dictItem struct {
	Fun func() interface{}
	Res string
}

// 字典项资源
var tableList = []dictItem{
	{ // 索引总表
		Fun: func() interface{} {
			type dictIndex struct{ model.Dict }
			return dictIndex{}
		},
		Res: "./res/index_info.json",
	},
	{ // 用户性别字典
		Fun: func() interface{} {
			type dictUserGender struct{ model.Dict }
			return dictUserGender{}
		},
		Res: "./res/user_gender.json",
	},
	{ // 用户水平等级
		Fun: func() interface{} {
			type dictUserEvaluate struct{ model.Dict }
			return dictUserEvaluate{}
		},
		Res: "./res/user_evaluate.json",
	},
	{ // 用户收货地址标签
		Fun: func() interface{} {
			type dictUserAddressTag struct{ model.Dict }
			return dictUserAddressTag{}
		},
		Res: "./res/user_address_tag.json",
	},
	{ // 消息通知类型
		Fun: func() interface{} {
			type dictNoticeType struct{ model.Dict }
			return dictNoticeType{}
		},
		Res: "./res/notice_type.json",
	},
	{ // 活动类型
		Fun: func() interface{} {
			type dictActivityType struct{ model.Dict }
			return dictActivityType{}
		},
		Res: "./res/activity_type.json",
	},
	{ // 帖子标签
		Fun: func() interface{} {
			type dictPostTag struct{ model.Dict }
			return dictPostTag{}
		},
		Res: "./res/post_tag.json",
	},
	{ // 省市区三级联动
		Fun: func() interface{} {
			type dictAddress struct{ model.Dict }
			return dictAddress{}
		},
		Res: "./res/address.json",
	},
	{ // 食谱菜系
		Fun: func() interface{} {
			type dictRecipeCuisine struct{ model.Dict }
			return dictRecipeCuisine{}
		},
		Res: "./res/recipe_cuisine.json",
	},
	{ // 食谱口味
		Fun: func() interface{} {
			type dictRecipeTaste struct{ model.Dict }
			return dictRecipeTaste{}
		},
		Res: "./res/recipe_taste.json",
	},
	{ // 食谱标签
		Fun: func() interface{} {
			type dictRecipeTag struct{ model.Dict }
			return dictRecipeTag{}
		},
		Res: "./res/recipe_tag.json",
	},
	{ // 食谱主材
		Fun: func() interface{} {
			type dictRecipeIngredientsMain struct{ model.Dict }
			return dictRecipeIngredientsMain{}
		},
		Res: "./res/recipe_ingredients_main.json",
	},
	{ // 食谱辅料
		Fun: func() interface{} {
			type dictRecipeIngredientsSub struct{ model.Dict }
			return dictRecipeIngredientsSub{}
		},
		Res: "./res/recipe_ingredients_sub.json",
	},
	{ // 勋章稀有度等级
		Fun: func() interface{} {
			type dictMedalRarity struct{ model.Dict }
			return dictMedalRarity{}
		},
		Res: "./res/medal_rarity.json",
	},
}

// 大写字母匹配
var reg = regexp.MustCompile("[A-Z]")

// InitDict 执行字典表的初始化操作
func InitDict() {
	println("开始初始化字典表")
	printDivider()
	db := common.InitDB(false)
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, item := range tableList {
			dict := item.Fun()
			name := reflect.TypeOf(dict).Name()
			name = strings.ToLower(reg.
				ReplaceAllString(name, "_${0}"))
			// 如果表不存在则创建
			if !tx.Migrator().HasTable(&dict) {
				println(fmt.Sprintf("%s 表不存在，正在创建", name))
				if err := tx.Migrator().
					CreateTable(dict); err != nil {
					return err
				}
			}
			// 加载资源文件并插入数据
			println(fmt.Sprintf("正在向 %s 表插入数据", name))
			if err := insertDict(tx,
				"sys_"+name, item.Res); err != nil {
				return err
			}
			printDivider()
		}
		return nil
	})
	if err != nil {
		panic("字典表初始化失败，已回滚:" + err.Error())
		return
	}
	println("字典表初始化完成")
}

// 加载资源文件并插入数据
func insertDict(tx *gorm.DB, name string, res string) error {
	res = strings.ReplaceAll(res, "./", "./script/")
	var jsonArray []dict
	if err := tool.ReadJsonFile(res, &jsonArray); err != nil {
		return err
	}
	result := json2Dict("0", jsonArray)
	if err := tx.Table(name).
		CreateInBatches(result, 1000).
		Error; err != nil {
		return err
	}
	return nil
}

type dict struct {
	Code     string  `json:"code,omitempty"`
	Name     string  `json:"name,omitempty"`
	Children *[]dict `json:"children,omitempty"`
}

// json转对象
func json2Dict(pCode string, result []dict) []model.Dict {
	var array []model.Dict
	for _, it := range result {
		code := it.Code
		name := it.Name
		array = append(array, model.Dict{
			OrmBase: model.OrmBase{
				ID:        tool.GenID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			SimpleDict: model.SimpleDict{
				PCode: pCode,
				Code:  code,
				Tag:   name,
				Info:  "",
			},
		})
		children := it.Children
		if children != nil {
			array = append(array, json2Dict(code, *children)...)
		}
	}
	return array
}

// 打印分割线
func printDivider() {
	println("----------------------------")
}
