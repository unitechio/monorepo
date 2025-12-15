package model

// 3. Block Data Structs (type specific)

// 3.1 Hero Block
type HeroBlockData struct {
	Title    string     `json:"title"`
	Subtitle string     `json:"subtitle,omitempty"`
	Image    ImageSEO   `json:"image"`
	CTA      *CTAButton `json:"cta,omitempty"`
}

// 3.2 Text Block
type TextBlockData struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content_html"` // or markdown
}

// 3.3 Image + Text Block
type ImageTextBlockData struct {
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Image         ImageSEO `json:"image"`
	ImagePosition string   `json:"image_position"` // left, right
}

// 3.4 Feature List Block
type FeatureListBlockData struct {
	Title string        `json:"title"`
	Items []FeatureItem `json:"items"`
}

type FeatureItem struct {
	Icon string `json:"icon"`
	Text string `json:"text"`
}

// 3.5 CTA Block
type CTABlockData struct {
	Title  string    `json:"title"`
	Button CTAButton `json:"button"`
}

// 3.6 Common Structs
type CTAButton struct {
	Text string `json:"text"`
	Link string `json:"link"`
}

type ImageSEO struct {
	URL    string `json:"url"`
	Alt    string `json:"alt"`
	Title  string `json:"title"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
