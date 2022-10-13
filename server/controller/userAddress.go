package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server/common"
	"server/controller/response"
	"server/middleware"
	"server/model"
	"server/tool"
)

// 用户收货地址请求体
type userAddressReq struct {
	Receiver      string   `json:"receiver" binding:"required,gte=2"`
	Contact       string   `json:"contact" binding:"required,phone"`
	AddressCodes  []string `json:"addressCodes" binding:"required,unique,gte=3"`
	AddressDetail string   `json:"addressDetail" binding:"required,gte=6"`
	TagCode       string   `json:"tagCode" binding:"dict=user_address_tag"`
	Default       bool     `json:"default" binding:"required"`
	Order         int64    `json:"order"`
}

// AddUserAddress 添加用户收货地址
func AddUserAddress(c *gin.Context) {
	// 获取请求体
	var req userAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	// 添加数据
	db := common.GetDB()
	result := model.UserAddress{
		OrmBase:       createBase(),
		Creator:       createCreator(c),
		Receiver:      req.Receiver,
		Contact:       req.Contact,
		AddressCodes:  req.AddressCodes,
		AddressDetail: req.AddressDetail,
		TagCode:       req.TagCode,
		Default:       req.Default,
		Order:         req.Order,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&result).Error; err != nil {
			return err
		}
		// 如果当前收货地址被标记为默认，则更新当前用户的其他收货地址为非默认
		if result.Default {
			if err := setAddressDefault(
				c, tx, result.ID); err != nil {
				return nil
			}
		}
		return nil
	})
	if err != nil {
		response.FailDef(c, -1, "收货地址创建失败")
		return
	}
	response.SuccessDef(c, result)
}

// UpdateUserAddress 更新用户收货地址
func UpdateUserAddress(c *gin.Context) {
	// 获取请求参数
	var req userAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	addressId := c.Param("addressId")
	if len(addressId) == 0 {
		response.FailParams(c, "收货地址id不能为空")
		return
	}
	var result model.UserAddress
	if hasNoRecord(&result, addressId) {
		response.FailParams(c, "收货地址不存在")
		return
	}
	if result.CreatorId != middleware.GetCurrUId(c) {
		response.FailParams(c, "您不是该收货地址的所有者")
		return
	}
	// 数据插入
	db := common.GetDB()
	result.Receiver = req.Receiver
	result.Contact = req.Contact
	result.AddressCodes = req.AddressCodes
	result.AddressDetail = req.AddressDetail
	result.TagCode = req.TagCode
	result.Default = req.Default
	result.Order = req.Order
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&result).Error; err != nil {
			return err
		}
		// 如果当前收货地址被标记为默认，则更新当前用户的其他收货地址为非默认
		if result.Default {
			if err := setAddressDefault(
				c, tx, addressId); err != nil {
				return nil
			}
		}
		return nil
	})
	if err != nil {
		response.FailDef(c, -1, "收货地址创建失败")
		return
	}
	response.SuccessDef(c, result)
}

// UpdateUserAddressDefault 修改用户收货地址为默认地址
func UpdateUserAddressDefault(c *gin.Context) {
	// 获取请求参数
	addressId := c.Param("addressId")
	if len(addressId) == 0 {
		response.FailParams(c, "收货地址id不能为空")
		return
	}
	var result model.UserAddress
	if hasNoRecord(&result, addressId) {
		response.FailParams(c, "收货地址不存在")
		return
	}
	if result.CreatorId != middleware.GetCurrUId(c) {
		response.FailParams(c, "您不是该收货地址的所有者")
		return
	}
	// 更新默认状态
	db := common.GetDB()
	if err := setAddressDefault(c, db, addressId); err != nil {
		response.FailDef(c, -1, "默认状态更新失败")
		return
	}
	response.SuccessDef(c, true)
}

// UpdateUserAddressOrder 修改用户收货地址排序
func UpdateUserAddressOrder(c *gin.Context) {
	// 获取请求参数
	var req struct {
		Order int64 `json:"order" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailParamsDef(c, err)
		return
	}
	addressId := c.Param("addressId")
	if len(addressId) == 0 {
		response.FailParams(c, "收货地址id不能为空")
		return
	}
	var result model.UserAddress
	if hasNoRecord(&result, addressId) {
		response.FailParams(c, "收货地址不存在")
		return
	}
	if result.CreatorId != middleware.GetCurrUId(c) {
		response.FailParams(c, "您不是该收货地址的所有者")
		return
	}
	// 更新默认状态
	db := common.GetDB()
	if err := db.Model(&result).
		UpdateColumn("order", req.Order).Error; err != nil {
		response.FailDef(c, -1, "排序更新失败")
		return
	}
	response.SuccessDef(c, true)
}

// GetAllUserAddress 获取全部收货地址信息
func GetAllUserAddress(c *gin.Context) {
	// 获取当前用户的收货地址
	db := common.GetDB()
	userId := middleware.GetCurrUId(c)
	var result []*model.UserAddress
	if err := db.Model(&model.UserAddress{}).
		Where("creator_id=?", userId).
		Preload("Creator").
		Order("`order`").
		Scan(&result).Error; err != nil {
		response.FailDef(c, -1, "获取收货地址失败")
		return
	}
	fillUserAddressInfo(result...)
	response.SuccessDef(c, tool.If(
		result != nil, result, []*model.UserAddress{}))
}

// GetUserAddressInfo 获取用户收货地址详情
func GetUserAddressInfo(c *gin.Context) {
	// 获取请求参数
	addressId := c.Param("addressId")
	if len(addressId) == 0 {
		response.FailParams(c, "收货地址id不能为空")
		return
	}
	var result model.UserAddress
	if hasNoRecord(&result, addressId) {
		response.FailParams(c, "收货地址不存在")
		return
	}
	if result.CreatorId != middleware.GetCurrUId(c) {
		response.FailParams(c, "您不是该收货地址的所有者")
		return
	}
	// 获取当前用户的收货地址
	fillUserAddressInfo(&result)
	response.SuccessDef(c, result)
}

// 更新收货地址为默认
func setAddressDefault(c *gin.Context, db *gorm.DB, id string) error {
	userId := middleware.GetCurrUId(c)
	return db.Exec("update sys_user_address set `default`=if(id=?,true,false) where creator_id = ?",
		id, userId).Error
}

// 填充收货地址信息
func fillUserAddressInfo(items ...*model.UserAddress) {
	db := common.GetDB()
	var codes []string
	for _, it := range items {
		codes = append(codes, it.TagCode)
	}
	var tags []*model.SimpleDict
	db.Table("sys_dict_user_address_tag").
		Where("code in ?", codes).
		Find(&tags)
	for i, it := range tags {
		items[i].Tag = it
	}
}
