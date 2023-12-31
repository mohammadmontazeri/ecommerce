package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string  `gorm:"column:username;unique;" json:"username"`
	Email    string  `gorm:"column:email;unique;not null;index:em_index" json:"email"`
	Password string  `gorm:"column:password;not null" json:"password"`
	Orders   []Order `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
