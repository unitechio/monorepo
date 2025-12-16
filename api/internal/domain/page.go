package domain

import "time"

// Page represents a page in the system
type Page struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"uniqueIndex;size:100" json:"code"`
	Slug      string    `gorm:"size:255" json:"slug"`
	Status    string    `gorm:"size:20;default:'draft'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName returns the table name for Page
func (Page) TableName() string {
	return "page"
}
