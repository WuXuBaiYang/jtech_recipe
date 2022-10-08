package main

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/model"
)

func main() {
	//r := gin.Default()
	//r.GET("test", testAPI)
	//panic(r.Run(":9528"))
	db := common.InitDB(true)
	err := db.Where("phone_number = ?", "18600574971").First(&model.User{}).Error
	println(err.Error())
}

type pagination struct {
	PageIndex int64 `form:"pageIndex" binding:"required,gte=1"`
	PageSize  int64 `form:"pageSize" binding:"required,gte=10"`
	Total     int64 `json:"total"`
}

// 测试接口
func testAPI(c *gin.Context) {
	var req pagination
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(-1, err.Error())
		return
	}
	c.JSON(200, req)
}
