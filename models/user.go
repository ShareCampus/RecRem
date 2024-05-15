package models

import (
	"recrem/config/db"
	"recrem/utils"

	"github.com/jinzhu/gorm"
)

// User 用户
type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(30);not null;" json:"username"`  // 用户名
	Email     string `gorm:"type:varchar(30);not null;" json:"email"`     // 邮箱
	Pwd       string `gorm:"type:varchar(100);not null;" json:"pwd"`      // 密码
	Nickname  string `gorm:"type:varchar(30);" json:"nickname"`           // 昵称
	UserImg   string `gorm:"type:varchar(255);not null;" json:"user_img"` // 用户头像
	Signature string `gorm:"type:varchar(255);" json:"signature"`         // 个性签名
}

// GetByUsername 根据用户名和密码获取用户
func (user User) GetByUsername() (User, error) {
	var u User
	err := db.Db.Where("`username` = ? or `email` = ?", user.Username, user.Username).
		First(&u).Error
	if gorm.IsRecordNotFoundError(err) {
		err = nil
	}

	return u, err
}

// Create 创建用户
func (user User) Create() error {
	hashedPwd, err := utils.EncryptPwd(user.Pwd) // 加密密码
	if err != nil {
		return err
	}

	user.Pwd = hashedPwd

	return db.Db.Create(&user).Error
}

// GetByEmail 根据邮箱获取用户
func (user User) GetByEmail() (User, error) {
	u := User{}
	err := db.Db.Where("`email` = ?", user.Email).First(&u).Error
	if gorm.IsRecordNotFoundError(err) {
		err = nil
	}

	return u, err
}

// UpdatePwd 修改密码
func (user User) UpdatePwd() error {
	hashedPwd, err := utils.EncryptPwd(user.Pwd) // 加密密码
	if err != nil {
		return err
	}

	return db.Db.Model(&User{}).Where("`email` = ?", user.Email).
		Update("pwd", hashedPwd).Error
}
