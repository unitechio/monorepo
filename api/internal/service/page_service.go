package service

import (
	"context"

	"github.com/unitechio/oss-monorepo/api/internal/model"
	"github.com/unitechio/oss-monorepo/api/internal/repository"
)

var pageRepo *repository.PageRepository

func InitPageService(repo *repository.PageRepository) {
	pageRepo = repo
}

// GetPageBySlug resolves the page content for the given language
func GetPageBySlug(slug string, lang string) (*model.PageResponse, error) {
	ctx := context.Background()

	// 1. Get Page
	page, err := pageRepo.GetPageBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// 2. Get SEO
	seo, err := pageRepo.GetPageSEO(ctx, page.ID, lang)
	if err != nil {
		return nil, err
	}
	if seo == nil && lang != "en" {
		// Fallback to EN if not found
		seo, err = pageRepo.GetPageSEO(ctx, page.ID, "en")
		if err != nil {
			return nil, err
		}
	}

	seoResponse := model.SEOResponse{}
	if seo != nil {
		seoResponse = model.SEOResponse{
			Title:       seo.Title,
			Description: seo.Description,
			OGTitle:     seo.OGTitle,
			OGDesc:      seo.OGDesc,
			OGImage:     seo.OGImage,
			Canonical:   seo.Canonical,
		}
	}

	// 3. Get Blocks
	blocks, err := pageRepo.GetPageBlocks(ctx, page.ID)
	if err != nil {
		return nil, err
	}

	blockIDs := make([]int64, len(blocks))
	for i, b := range blocks {
		blockIDs[i] = b.ID
	}

	// 4. Get Translations
	// Strategy: Load requested lang, and load EN as fallback source
	translations, err := pageRepo.GetBlockTranslations(ctx, blockIDs, lang)
	if err != nil {
		return nil, err
	}

	var enTranslations map[int64]model.PageBlockTranslation
	if lang != "en" {
		// Identify missing blocks
		var missingIDs []int64
		for _, bid := range blockIDs {
			if _, ok := translations[bid]; !ok {
				missingIDs = append(missingIDs, bid)
			}
		}

		if len(missingIDs) > 0 {
			enTranslations, err = pageRepo.GetBlockTranslations(ctx, missingIDs, "en")
			if err != nil {
				return nil, err // or log and ignore
			}
		}
	}

	// 5. Assemble Response
	blockResponses := make([]model.BlockResponse, len(blocks))
	for i, b := range blocks {
		var dataJson = []byte("{}")

		// Try requested lang
		if t, ok := translations[b.ID]; ok {
			dataJson = t.Data
		} else if lang != "en" {
			// Try fallback EN
			if t, ok := enTranslations[b.ID]; ok {
				dataJson = t.Data
			}
		}

		blockResponses[i] = model.BlockResponse{
			Code: b.Code,
			Type: b.Type,
			Data: dataJson,
		}
	}

	resp := &model.PageResponse{
		Code:   page.Code,
		Slug:   page.Slug,
		Lang:   lang,
		SEO:    seoResponse,
		Blocks: blockResponses,
	}

	return resp, nil
}
