package main

import (
	"github.com/jinzhu/gorm"
)

// UserInfo 用户信息
type UserInfo struct {
	Username string
	Password string
}

func connect() (db *gorm.DB) {
	db, err := gorm.Open("mysql", "root:qwerasdf83448614@(127.0.0.1:3306)/users?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&UserInfo{})
	return db
}
