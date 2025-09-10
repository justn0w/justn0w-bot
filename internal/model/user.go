package model

import (
	"justn0w-bot/config"
	"time"
)

type User struct {
	Id        int64     `json:"id"`         //用户ID
	Name      string    `json:"name"`       //用户名
	Password  string    `json:"password"`   //用户密码
	CreatedAt time.Time `json:"created_at"` //创建时间
	UpdatedAt time.Time `json:"updated_at"` //更新时间
}

func (User) TableName() string {
	return "user"
}

func FindUserByName(name string) (User, error) {
	var user User
	tx := config.Db.Where("name = ?", name).Find(&user)
	return user, tx.Error
}

func InsertUser(user User) error {
	tx := config.Db.Create(&user)
	return tx.Error
}
