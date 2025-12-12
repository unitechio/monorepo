package model

// I18nMap holds translations for a field, keyed by language code (e.g., "en", "vi").
type I18nMap map[string]string

// Get returns the string for the requested language, defaulting to "en" or the first available.
func (m I18nMap) Get(lang string) string {
	if val, ok := m[lang]; ok {
		return val
	}
	if val, ok := m["en"]; ok {
		return val
	}
	for _, val := range m {
		return val
	}
	return ""
}

// Page represents a full page structure.
type Page struct {
	Slug   string                 `json:"slug"`
	Title  string                 `json:"title"` // Resolved title
	Config map[string]interface{} `json:"config,omitempty"`
	Blocks []Block                `json:"blocks"`
}

// Block represents a section of the page.
type Block struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"` // e.g., "hero", "grid"
	Properties map[string]interface{} `json:"properties,omitempty"`
	Items      []Item                 `json:"items,omitempty"`
	Title      string                 `json:"title,omitempty"` // Resolved title
}

// Item represents an individual content unit.
type Item struct {
	ID       string                 `json:"id"`
	Title    string                 `json:"title,omitempty"`    // Resolved title
	Subtitle string                 `json:"subtitle,omitempty"` // Resolved subtitle
	Image    string                 `json:"image,omitempty"`
	Link     string                 `json:"link,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
