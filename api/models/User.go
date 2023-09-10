package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null;type:varchar(255)"`
	Password string `gorm:"uniqueIndex;not null;type:varchar(255)"`
	Email    string `gorm:"uniqueIndex;not null;type:varchar(255)"`
}
