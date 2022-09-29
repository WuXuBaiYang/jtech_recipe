package model

import "time"

// Activity 活动信息结构体
type Activity struct {
	OrmBase

	CycleTime int64 `json:"cycleTime" gorm:"comment:周期持续时间"`
	// 长期活动的话，定时器在结算上一个周期之后，会自动发起下一个周期的活动
	Always    bool     `json:"always" gorm:"comment:是否为长期活动"`
	Title     string   `json:"title" gorm:"varchar(80);comment:活动标题"`
	Contents  *[]any   `json:"contents,omitempty" gorm:"type:json;serializer:json;comment:活动内容"`
	Url       *string  `json:"url,omitempty" gorm:"varchar(200);comment:活动介绍网址"`
	TypeCodes []string `json:"typeCodes" gorm:"type:json;serializer:json;not null;comment:活动接受投稿的类型字典码集合"`
}

// ActivityRecord 活动信息记录
type ActivityRecord struct {
	OrmBase

	BeginTime  time.Time `json:"beginTime" gorm:"not null;comment:开始时间"`
	EndTime    time.Time `json:"endTime" gorm:"not null;comment:结束时间"`
	ActivityId string    `json:"activityId" gorm:"not null;comment:活动id"`
	Activity   Activity  `json:"activity" gorm:"foreignKey:ActivityId"`
}
