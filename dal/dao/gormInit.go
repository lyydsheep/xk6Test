package dao

import (
	"email/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	dbMaster *gorm.DB
)

func DB() *Query {
	return Use(dbMaster)
}

func DBMaster() *Query {
	return Use(dbMaster)
}

func InitDB() {
	dbMaster = initDB(&config.DB.Master)
}

func initDB(option *config.DBConfigOptions) *gorm.DB {
	// 默认使用 MySQL
	db, err := gorm.Open(mysql.Open(option.Dsn), &gorm.Config{
		//Logger: _GormLogger,
	})
	if err != nil {
		panic(err)
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxOpenConns(option.MaxOpen)
	sqlDb.SetMaxIdleConns(option.MaxIdle)
	sqlDb.SetConnMaxLifetime(time.Duration(option.MaxLifeTime) * time.Second)
	return db
}
