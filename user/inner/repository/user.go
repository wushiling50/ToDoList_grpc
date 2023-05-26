package repository

import (
	"errors"
	"time"

	"main/ToDoList_grpc/service"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Uid            int64
	UserName       string `gorm:"unique"`
	PasswordDigest string //存储的是密文

}

const (
	PasswordCost = 12 //加密难度
)

// 创建用户
func (*User) UserCreate(req *service.UserRequest) (user User, err error) {
	var count int64
	DB.Where("user_name=?", req.UserName).Count(&count)
	if count != 0 {
		return User{}, errors.New("该用户已存在")
	}

	user = User{
		Uid:      time.Now().Unix(),
		UserName: req.UserName,
	}
	//加密
	user.SetPassword(req.Password)
	err = DB.Create(&user).Error
	return user, err
}

// 加密
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// 验证密码
func (user *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err
}

// 检查用户信息
func (user *User) FindUserInfo(req *service.UserRequest) (err error) {
	if exist := user.CheckUserExist(req); exist {
		return nil
	}
	return errors.New("该用户不存在")
}

// 检查用户是否存在
func (user *User) CheckUserExist(req *service.UserRequest) bool {
	if err := DB.Where("user_name=?", req.UserName).First(&user).Error; err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}
