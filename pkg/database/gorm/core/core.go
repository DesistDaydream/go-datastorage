package core

import (
	"github.com/DesistDaydream/go-datastorage/pkg/database/gorm/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("连接数据库失败: %v", err)
	}

	db.AutoMigrate(&models.User{})
}
