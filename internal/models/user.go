package models

import "time"

type User struct {
	ID           int64  `json:"id" gorm:"primaryKey"`
	FirstName    string `json:"first_name" gorm:"not null"`
	LastName     string `json:"last_name" gorm:"not null"`
	Username     string `json:"username" gorm:"not null;uniqueIndex"`
	PasswordHash string `json:"-" gorm:"not null"`

	Opinions []Opinion `json:"opinions" gorm:"foreignKey:UserID"`

	Following []User    `gorm:"many2many:user_follows;foreignKey:ID;joinForeignKey:FollowerID;joinReferences:FollowingID"`
	Followers []User    `gorm:"many2many:user_follows;foreignKey:ID;joinForeignKey:FollowingID;joinReferences:FollowerID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
