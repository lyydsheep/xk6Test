package util

import (
	"bufio"
	"email/common/enum"
	"email/dal/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
)

func InsertUser(filePath string) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// 数据库连接字符串
	dsn := "crm_prod:1qaz#EDC@pv12@tcp(rm-t4nvxa5dq73ua02khuo.mysql.singapore.rds.aliyuncs.com:3306)/crm_prod?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 读取文件并插入数据
	tx := db.Begin()
	defer tx.Rollback()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		email := scanner.Text()
		email = strings.ReplaceAll(email, "|", "")
		email = strings.TrimSpace(email)
		entry := model.UsrUser{
			UserID:   email,
			FullName: email,
			Email:    email,
			Language: enum.EmailTemplateEN,
			Cid:      2,
			Tags:     "odps",
		}

		// 插入数据
		if err := tx.Create(&entry).Error; err != nil {
			log.Printf("Failed to insert data: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	tx.Commit()
	fmt.Println("Data insertion complete.")
}
