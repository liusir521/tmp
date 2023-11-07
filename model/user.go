package model

type User struct {
	ID    int64
	Name  string  `gorm:"column:name"`
	Phone string  `gorm:"column:phone"`
	Money float64 `gorm:"column:money"`
}

func (u User) TableName() string {
	//绑定MYSQL表名为users
	return "user"
}
