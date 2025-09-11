package service

import (
	"errors"
	"justn0w-bot/internal/model"
	"justn0w-bot/pkg/utils"
	"time"
)

type UserService struct {
}

func (u UserService) Register(name string, password string) error {
	//1 校验数据库中是否存在该用户
	userDb, err := model.FindUserByName(name)
	if err != nil {
		return err
	}
	if userDb.Id > 0 {
		return errors.New("用户已存在")
	}

	//2 对密码进行加密
	hashPassword, err := utils.GenerateHash(password)
	if err != nil {
		return err
	}

	//3 插入数据库
	user := model.User{
		Name:      name,
		Password:  hashPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return model.InsertUser(user)
}

func (u UserService) Login(name string, password string) error {
	// 1 判断用户是否存在
	userDb, err := model.FindUserByName(name)
	if err != nil {
		return err
	}
	if userDb.Id <= 0 {
		return errors.New("用户不存在")
	}

	// 2 比较密码是否一致
	isSame := utils.CompareHashAndPassword(userDb.Password, password)
	if !isSame {
		return errors.New("密码错误")
	} else {
		return nil
	}
}
