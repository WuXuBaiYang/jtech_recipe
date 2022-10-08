package tool

import (
	"server/model"
)

// 评论所属类型对照表
var commentMap = map[string]interface{}{
	"post":     &model.Post{},
	"recipe":   &model.Recipe{},
	"menu":     &model.RecipeMenu{},
	"activity": &model.Activity{},
}

// CommentType 获取评论类型
func CommentType(v string) interface{} {
	return commentMap[v]
}

// CommentTypeVerify 验证评论类型是否存在是否存在
func CommentTypeVerify(v string) bool {
	return CommentType(v) != nil
}
