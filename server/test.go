package main

import (
	"context"
	"github.com/gin-gonic/gin"
)

var c = context.Background()

func main() {
	r := gin.Default()
	r.GET("test", testAPI)
	panic(r.Run(":9528"))
}

type pagination struct {
	PageIndex int64 `form:"pageIndex" binding:"required,gte=1"`
	PageSize  int64 `form:"pageSize" binding:"required,gte=10"`
	Total     int64 `json:"total"`
}

// 测试接口
func testAPI(c *gin.Context) {
	//var req pagination
	//if err := c.ShouldBindQuery(&req); err != nil {
	//	c.JSON(-1, err.Error())
	//	return
	//}
	c.JSON(200, nil)
}
