package model

import "encoding/json"

// 2. API Public Response (Next.js consume)

// 2.1 Page API Response
type PageResponse struct {
	Code   string          `json:"code"`
	Slug   string          `json:"slug"`
	Lang   string          `json:"lang"`
	SEO    SEOResponse     `json:"seo"`
	Blocks []BlockResponse `json:"blocks"`
}

// 2.2 SEO Response
type SEOResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	OGTitle     string `json:"og_title,omitempty"`
	OGDesc      string `json:"og_description,omitempty"`
	OGImage     string `json:"og_image,omitempty"`
	Canonical   string `json:"canonical,omitempty"`
}

// 2.3 Block Response (key design)
type BlockResponse struct {
	Code string          `json:"code"`
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
