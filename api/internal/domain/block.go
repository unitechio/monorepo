package domain

import "time"

// Block represents a block in a page
type Block struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	PageID     int64     `gorm:"index" json:"page_id"`
	Code       string    `gorm:"size:100" json:"code"`
	Type       string    `gorm:"size:50" json:"type"`
	BlockOrder int       `gorm:"default:0" json:"block_order"`
	Status     string    `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName returns the table name for Block
func (Block) TableName() string {
	return "page_block"
}
