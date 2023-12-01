package model

import (
	"time"
)

type TgOrder struct {
	ID         int64
	UserCouunt int64     `gorm:"column:usercount" json:"userCount"`
	DriverId   int64     `gorm:"column:driverid" json:"driverId"`
	Carid      string    `gorm:"column:carid" json:"carId"`
	Money      float64   `gorm:"column:money" json:"money"`
	CreateTime time.Time `gorm:"column:createtime" json:"createTime"`
	StartPlace string    `gorm:"column:startplace" json:"startPlace"`
	EndPlace   string    `gorm:"column:endplace" json:"endPlace"`
	Status     string    `gorm:"column:status" json:"status"`
}

func (u TgOrder) TableName() string {
	//绑定MYSQL表名为users
	return "tgorder"
}
