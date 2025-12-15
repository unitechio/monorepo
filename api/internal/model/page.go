package model

import (
	"encoding/json"
)

// 1. Domain / DB layer (Go struct)

// 1.1 Page
type Page struct {
	ID     int64  `db:"id"`
	Code   string `db:"code"`   // home, about, contact
	Slug   string `db:"slug"`   // /, /about
	Status string `db:"status"` // draft, published
}

// 1.2 Page SEO (multi-language)
type PageSEO struct {
	ID          int64  `db:"id"`
	PageID      int64  `db:"page_id"`
	Lang        string `db:"lang"` // vi, en
	Title       string `db:"title"`
	Description string `db:"description"`
	OGTitle     string `db:"og_title"`
	OGDesc      string `db:"og_description"`
	OGImage     string `db:"og_image"`
	Canonical   string `db:"canonical"`
}

// 1.3 Page Block (language independent)
type PageBlock struct {
	ID     int64  `db:"id"`
	PageID int64  `db:"page_id"`
	Code   string `db:"code"` // home_hero, about_intro
	Type   string `db:"type"` // hero, text, image_text
	Order  int    `db:"order"`
	Status string `db:"status"` // active, inactive
}

// 1.4 Page Block Translation (multi-language)
type PageBlockTranslation struct {
	ID      int64           `db:"id"`
	BlockID int64           `db:"block_id"`
	Lang    string          `db:"lang"`
	Data    json.RawMessage `db:"data_json"` // block content by type
}

// 1.5 Image SEO (reusable)
type Image struct {
	ID     int64  `db:"id"`
	URL    string `db:"url"`
	Width  int    `db:"width"`
	Height int    `db:"height"`
}

type ImageTranslation struct {
	ImageID int64  `db:"image_id"`
	Lang    string `db:"lang"`
	Alt     string `db:"alt"`
	Title   string `db:"title"`
}
