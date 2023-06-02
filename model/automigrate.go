package model

import (
	"fmt"

	"gorm.io/gorm"
)

func RunAutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	fmt.Println("User Database Migrated")

	fmt.Println("Media Database Migrated")
}
