package model

type Driver struct {
	ID        int64
	Name      string  `gorm:"column:name"`
	Phone     string  `gorm:"column:phone"`
	Money     float64 `gorm:"column:money"`
	IsWorking bool    `gorm:"column:isworking"`
	CarId     string  `gorm:"column:carid"`
	RunCount  int64   `gorm:"column:runcount"`
	Password  string  `gorm:"column:password"`
}

func (u Driver) TableName() string {
	//绑定MYSQL表名为users
	return "driver"
}
