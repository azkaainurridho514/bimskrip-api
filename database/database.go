package database

import (
	"log"

	"github.com/azkaainurridho514/bimskrip/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
    return db
}

func InitDB() {
	var err error
	dsn := "root:@tcp(localhost:3306)/bimskrip-v1?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.AutoMigrate(&model.Role{}, &model.StatusName{}, &model.ProgressName{}, &model.User{}, &model.Progress{}, &model.Schedule{})
	seedData()
}

func seedData() {
	var roleCount int64
	db.Model(&model.Role{}).Count(&roleCount)
	if roleCount == 0 {
		roles := []model.Role{
			{Name: "Mahasiswa"},
			{Name: "Dosen"},
			{Name: "Staff"},
		}
		db.Create(&roles)
	}
	var statusCount int64
	db.Model(&model.StatusName{}).Count(&statusCount)
	if statusCount == 0 {
		statuses := []model.StatusName{
			{Name: "Progress"},
			{Name: "Diterima"},
			{Name: "Ditolak"},
		}
		db.Create(&statuses)
	}

	// Seed progress names
	var progressCount int64
	db.Model(&model.ProgressName{}).Count(&progressCount)
	if progressCount == 0 {
		progressNames := []model.ProgressName{
			{Name: "Proposal"},
			{Name: "BAB 1"},
			{Name: "BAB 2"},
			{Name: "BAB 3"},
			{Name: "BAB 4"},
			{Name: "BAB 5"},
		}
		db.Create(&progressNames)
	}
}