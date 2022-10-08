package tool

import (
	"reflect"
	"server/model"
)

// 评论所属类型对照表
var commentBelongMap = map[string]reflect.Type{
	"post":   reflect.TypeOf(model.Post{}),
	"recipe": reflect.TypeOf(model.Recipe{}),
	"menu":   reflect.TypeOf(model.RecipeMenu{}),
}

// CommentBelongVerify 验证评论类型是否存在是否存在
func CommentBelongVerify(v any) bool {
	switch v.(type) {
	case string:
		return Platform2Int(v.(string)) != nil
	case int:
		return Platform2Tag(v.(int)) != nil
	}
	return false
}
