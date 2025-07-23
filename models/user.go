package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // Embeds ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `gorm:"size:255;not null"`
	Email      string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
}
