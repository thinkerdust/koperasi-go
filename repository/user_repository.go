package repository

import (
	"koperasi-go/db"
	"koperasi-go/model"
)

func FindUserByNIK(nik string) (model.User, error) {
	var user model.User
	err := db.DB.Where("nik = ?", nik).First(&user).Error
	return user, err
}

func CreateUser(user *model.User) error {
	return db.DB.Create(user).Error
}

func CheckUserToken(userId uint, token string) (bool, error) {
	var dbToken string
	err := db.DB.Raw("SELECT token FROM users WHERE id = ?", userId).Scan(&dbToken).Error
	if err != nil {
		return false, err
	}

	return dbToken == token, nil
}

func ClearUserToken(userId uint) error {
	return db.DB.Exec("UPDATE users SET token = NULL WHERE id = ?", userId).Error
}

func UpdateUserToken(userId uint, token string) error {
	return db.DB.Model(&model.User{}).
		Where("id = ?", userId).
		Update("token", token).Error
}