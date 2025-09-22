package repository

import (
	"koperasi-go/db"
	"koperasi-go/model"
)

func FindUserByUsername(username string) (model.User, error) {
	var user model.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

func CreateUser(user *model.User) error {
	return db.DB.Create(user).Error
}
