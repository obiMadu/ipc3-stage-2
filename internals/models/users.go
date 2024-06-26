package models

import "gorm.io/gorm"

type Users struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"unique;not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Fullname string `json:"fullname,omitempty"`
}

func CreateUser(db *gorm.DB, user Users) error {
	return db.Create(&user).Error
}

func GetAll(db *gorm.DB) ([]Users, error) {
	var users []Users
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByID(db *gorm.DB, id uint) (*Users, error) {
	var user Users
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(db *gorm.DB, username string) (*Users, error) {
	var user Users
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUserByID(db *gorm.DB, id uint, user Users) error {
	return db.Model(&Users{}).Where("id = ?", id).Updates(user).Error
}

func UpdateUserByUsername(db *gorm.DB, username string, user Users) error {
	return db.Model(&Users{}).Where("username = ?", username).Updates(user).Error
}

func DeleteUserByID(db *gorm.DB, id uint) error {
	user, err := GetUserByID(db, id)
	if err != nil {
		return err
	}
	return db.Delete(&Users{}, user).Error
}

func DeleteUserByUsername(db *gorm.DB, username string) error {
	user, err := GetUserByUsername(db, username)
	if err != nil {
		return err
	}

	return db.Delete(&user).Error
}
