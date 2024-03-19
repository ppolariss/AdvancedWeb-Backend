package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
)

var DB *gorm.DB

var Config struct {
	DbURL string `env:"DB_URL" envDefault:"root:root@tcp(localhost:3306)/advanced_web?charset=utf8mb4&parseTime=True&loc=Local"`
}

func InitDB() error {
	dsn := "root:root@tcp(localhost:3306)/advanced_web?charset=utf8mb4&parseTime=True&loc=Local"
	// 根据你的 MySQL 配置进行修改
	var err error
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到终端
		logger.Config{
			LogLevel: logger.Info, // 设置日志级别为 Info，以打印 SQL 语句
		},
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `Login` would be `user` with this option enabled
		}})
	if err != nil {
		return err
	}

	// 迁移数据库，确保 Login 表存在
	err = DB.AutoMigrate(&User{}, &UserJwtSecret{})
	if err != nil {
		return err
	}
	err = DB.AutoMigrate(&Message{})

	return err
}

//func setupDatabase() (*gorm.DB, error) {
//	dsn := "root:root@tcp(localhost:3306)/price_comparator"
//	// 根据你的 MySQL 配置进行修改
//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true, // use singular table name, table for `Login` would be `user` with this option enabled
//		}})
//	if err != nil {
//		return nil, err
//	}
//
//	// 迁移数据库，确保 Login 表存在
//	err = db.AutoMigrate(&Login{})
//	if err != nil {
//		return nil, err
//	}
//
//	return db, nil
//}
