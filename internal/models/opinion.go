package models

import "time"

type Opinion struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	Content   string `json:"content"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
