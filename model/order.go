package model

import (
	"encoding/json"
	"time"
)

type Order struct {
	ID         int64
	UserId     int64     `gorm:"column:userid" json:"userId"`
	DriverId   int64     `gorm:"column:driverid" json:"driverId"`
	Carid      string    `gorm:"column:carid" json:"carId"`
	Money      float64   `gorm:"column:money" json:"money"`
	CreateTime time.Time `gorm:"column:createtime" json:"createTime"`
	StartPlace string    `gorm:"column:startplace" json:"startPlace"`
	EndPlace   string    `gorm:"column:endplace" json:"endPlace"`
	Status     string    `gorm:"column:status" json:"status"`
}

func (u Order) TableName() string {
	//绑定MYSQL表名为users
	return "order"
}

func (m Order) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m Order) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
