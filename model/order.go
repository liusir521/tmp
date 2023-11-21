package model

import "time"

type Order struct {
	ID         int64
	UserId     string    `gorm:"column:userid" json:"userId"`
	DriverId   string    `gorm:"column:driverid" json:"driverId"`
	Carid      string    `gorm:"column:carid" json:"carId"`
	Money      float64   `gorm:"column:money" json:"money"`
	CreateTime time.Time `gorm:"column:createtime" json:"createTime"`
	StartPlace string    `gorm:"column:stratplace" json:"startPlace"`
	EndPlace   string    `gorm:"column:endplace" json:"endPlace"`
	Status     string    `gorm:"column:status" json:"status"`
}

func (u Order) TableName() string {
	//绑定MYSQL表名为users
	return "order"
}
