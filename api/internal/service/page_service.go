package service

import (
	"github.com/unitechio/oss-backend/internal/model"
)

// Mock Data Storage (in-memory)
// We store "Source" data using the I18nMap for translatable fields.
type SourceItem struct {
	ID       string
	Title    model.I18nMap
	Subtitle model.I18nMap
	Image    string
	Link     string
}

type SourceBlock struct {
	ID    string
	Type  string
	Title model.I18nMap
	Items []SourceItem
}

type SourcePage struct {
	Slug   string
	Title  model.I18nMap
	Blocks []SourceBlock
}

var mockPages = map[string]SourcePage{
	"home": {
		Slug:  "home",
		Title: model.I18nMap{"en": "Home Page", "vi": "Trang Chủ"},
		Blocks: []SourceBlock{
			{
				ID:    "b1",
				Type:  "hero",
				Title: model.I18nMap{"en": "Welcome", "vi": "Chào mừng"},
				Items: []SourceItem{
					{
						ID:       "i1",
						Title:    model.I18nMap{"en": "Discover Amazing Things", "vi": "Khám Phá Những Điều Tuyệt Vời"},
						Subtitle: model.I18nMap{"en": "We build future", "vi": "Chúng tôi xây dựng tương lai"},
						Image:    "/images/hero.jpg",
						Link:     "/about",
					},
				},
			},
			{
				ID:    "b2",
				Type:  "features",
				Title: model.I18nMap{"en": "Our Features", "vi": "Tính Năng"},
				Items: []SourceItem{
					{
						ID:       "i2",
						Title:    model.I18nMap{"en": "Fast", "vi": "Nhanh"},
						Subtitle: model.I18nMap{"en": "Lightning speed", "vi": "Tốc độ ánh sáng"},
					},
					{
						ID:       "i3",
						Title:    model.I18nMap{"en": "Secure", "vi": "Bảo Mật"},
						Subtitle: model.I18nMap{"en": "Top notch security", "vi": "Bảo mật hàng đầu"},
					},
				},
			},
		},
	},
}

// GetPageByType returns a resolved Page struct with content in the requested language.
func GetPageBySlug(slug string, lang string) *model.Page {
	source, exists := mockPages[slug]
	if !exists {
		return nil
	}

	// Transform SourcePage (with maps) to model.Page (with strings)
	response := &model.Page{
		Slug:   source.Slug,
		Title:  source.Title.Get(lang),
		Blocks: make([]model.Block, len(source.Blocks)),
	}

	for i, b := range source.Blocks {
		resBlock := model.Block{
			ID:    b.ID,
			Type:  b.Type,
			Title: b.Title.Get(lang),
			Items: make([]model.Item, len(b.Items)),
		}
		for j, it := range b.Items {
			resBlock.Items[j] = model.Item{
				ID:       it.ID,
				Title:    it.Title.Get(lang),
				Subtitle: it.Subtitle.Get(lang),
				Image:    it.Image,
				Link:     it.Link,
			}
		}
		response.Blocks[i] = resBlock
	}

	return response
}
