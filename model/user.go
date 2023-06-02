package model

import "gorm.io/gorm"

// Model struct
type User struct {
	gorm.Model
	Username string `gorm:"primaryKey" gorm:"unique_index;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Email    string `gorm:"unique_index;not null" json:"email"`
	FullName string `json:"full_name"`
	IsAdmin  bool   `json:"is_admin"`
}
