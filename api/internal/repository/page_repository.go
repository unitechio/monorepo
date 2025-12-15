package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/unitechio/oss-backend/internal/infrastructure/database"
	"github.com/unitechio/oss-backend/internal/model"
)

type PageRepository struct {
	DB *sql.DB
}

func NewPageRepository() *PageRepository {
	return &PageRepository{DB: database.DB}
}

func (r *PageRepository) GetPageBySlug(ctx context.Context, slug string) (*model.Page, error) {
	var page model.Page
	query := `SELECT id, code, slug, status FROM page WHERE slug = $1 AND status = 'published'`
	err := r.DB.QueryRowContext(ctx, query, slug).Scan(&page.ID, &page.Code, &page.Slug, &page.Status)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepository) GetPageSEO(ctx context.Context, pageID int64, lang string) (*model.PageSEO, error) {
	var seo model.PageSEO
	query := `SELECT id, page_id, lang, title, description, og_title, og_description, og_image, canonical 
			  FROM page_seo WHERE page_id = $1 AND lang = $2`
	err := r.DB.QueryRowContext(ctx, query, pageID, lang).Scan(
		&seo.ID, &seo.PageID, &seo.Lang,
		&seo.Title, &seo.Description,
		&seo.OGTitle, &seo.OGDesc, &seo.OGImage, &seo.Canonical,
	)
	if err == sql.ErrNoRows {
		// Fallback to 'en' or return zero value/nil handling by caller?
		// Caller logic says: "Resolve... if ok return data, else return 'en'".
		// I will handle fallback here or in service. Service is better for logic.
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}
	return &seo, nil
}

func (r *PageRepository) GetPageBlocks(ctx context.Context, pageID int64) ([]model.PageBlock, error) {
	query := `SELECT id, page_id, code, type, block_order, status 
			  FROM page_block WHERE page_id = $1 AND status = 'active' ORDER BY block_order`
	rows, err := r.DB.QueryContext(ctx, query, pageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []model.PageBlock
	for rows.Next() {
		var b model.PageBlock
		if err := rows.Scan(&b.ID, &b.PageID, &b.Code, &b.Type, &b.Order, &b.Status); err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

func (r *PageRepository) GetBlockTranslations(ctx context.Context, blockIDs []int64, lang string) (map[int64]model.PageBlockTranslation, error) {
	if len(blockIDs) == 0 {
		return nil, nil
	}

	// Create placeholder string for IN clause
	placeholders := ""
	args := make([]interface{}, len(blockIDs)+1)
	for i, id := range blockIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	args[len(blockIDs)] = lang

	query := fmt.Sprintf(`SELECT id, block_id, lang, data_json 
			  FROM page_block_translation WHERE block_id IN (%s) AND lang = $%d`, placeholders, len(blockIDs)+1)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]model.PageBlockTranslation)
	for rows.Next() {
		var t model.PageBlockTranslation
		var data []byte
		if err := rows.Scan(&t.ID, &t.BlockID, &t.Lang, &data); err != nil {
			return nil, err
		}
		t.Data = json.RawMessage(data)
		result[t.BlockID] = t
	}
	return result, nil
}
