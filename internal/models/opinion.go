package models

import (
	"errors"
	"time"
)

type Opinion struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (o *Opinion) ValidateContentCannotBeEmpty() error {
	if len(o.Content) == 0 {
		return errors.New("content cannot be empty")
	}
	return nil
}
