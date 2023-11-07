package model

import "time"

type Order struct {
	ID         int64
	UserId     string    `gorm:"column:userid"`
	DriverId   string    `gorm:"column:driverid"`
	Money      float64   `gorm:"column:money"`
	CreateTime time.Time `gorm:"column:createtime"`
	StartPlace string    `gorm:"column:stratplace"`
	EndPlace   string    `gorm:"column:endplace"`
	Status     string    `gorm:"column:status"`
}

func (u Order) TableName() string {
	//绑定MYSQL表名为users
	return "order"
}
