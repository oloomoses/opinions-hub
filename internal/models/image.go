package models

import "time"

type Image struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	OpinionID int64     `json:"opinion_id" gorm:"index;not null"`
	URL       string    `json:"url" gorm:"not null"`
	MimeType  string    `json:"mime_type" gorm:"not null"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
