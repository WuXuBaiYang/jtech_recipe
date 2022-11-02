package model

import "time"

// Activity 活动信息结构体
type Activity struct {
	OrmBase

	CycleDays int64 `json:"cycleDays" gorm:"not null;comment:周期持续天数"`
	// 长期活动的话，定时器在结算上一个周期之后，会自动发起下一个周期的活动
	Always    bool     `json:"always" gorm:"not null;comment:是否为长期活动"`
	Title     string   `json:"title" gorm:"varchar(80);not null;comment:活动标题"`
	Url       string   `json:"url" gorm:"varchar(200);not null;comment:活动介绍网址"`
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

// ActivityType 活动类型枚举
type ActivityType string

const (
	PostActivity   ActivityType = "11"
	MenuActivity   ActivityType = "12"
	RecipeActivity ActivityType = "13"
)
