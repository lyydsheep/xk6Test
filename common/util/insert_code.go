package util

import (
	"bufio"
	"email/dal/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func InsertCode(filePath string, amount int32) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// 数据库连接字符串
	dsn := "user02:+r1zvGQg%kB~IFoICu(*@tcp(47.86.177.131:3306)/ccpay_02?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 读取文件并插入数据
	tx := db.Begin()
	defer tx.Rollback()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		code := scanner.Text()
		entry := model.EmlRedemptionCode{
			Cid:    2,
			Code:   code,
			Amount: amount,
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
