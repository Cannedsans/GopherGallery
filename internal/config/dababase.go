package config

import (
	"github.com/Cannedsans/GopherGallery/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DataBase *gorm.DB
)

func StartBd() {
	db, err := gorm.Open(sqlite.Open("./data/database.db"), &gorm.Config{})

	if err != nil {
		return
	}

	db.AutoMigrate(&models.ImageModel{})

	DataBase = db
}
