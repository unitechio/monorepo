-- 1. Table: page
CREATE TABLE IF NOT EXISTS page (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,      -- home, about
    slug VARCHAR(255) NOT NULL,              -- /, /about
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_page_code ON page(code);

-- 2. Table: page_seo (multi-language)
CREATE TABLE IF NOT EXISTS page_seo (
    id BIGSERIAL PRIMARY KEY,
    page_id BIGINT NOT NULL,
    lang VARCHAR(10) NOT NULL,               -- vi, en

    title VARCHAR(255) NOT NULL,
    description TEXT,
    og_title VARCHAR(255),
    og_description TEXT,
    og_image VARCHAR(500),
    canonical VARCHAR(500),

    UNIQUE (page_id, lang),
    FOREIGN KEY (page_id) REFERENCES page(id) ON DELETE CASCADE
);

-- 3. Table: page_block
CREATE TABLE IF NOT EXISTS page_block (
    id BIGSERIAL PRIMARY KEY,
    page_id BIGINT NOT NULL,

    code VARCHAR(100) NOT NULL,               -- home_hero, home_feature
    type VARCHAR(50) NOT NULL,                -- hero, text, image_text
    block_order INT NOT NULL DEFAULT 0, -- renamed from 'order' to Avoid keyword conflict, though 'order' is allowed in quotes
    status VARCHAR(20) NOT NULL DEFAULT 'active',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (page_id, code),
    FOREIGN KEY (page_id) REFERENCES page(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_page_block_page ON page_block(page_id);

-- 4. Table: page_block_translation
CREATE TABLE IF NOT EXISTS page_block_translation (
    id BIGSERIAL PRIMARY KEY,
    block_id BIGINT NOT NULL,
    lang VARCHAR(10) NOT NULL,                -- vi, en

    data_json JSONB NOT NULL, -- Using JSONB for better performance in Postgres

    UNIQUE (block_id, lang),
    FOREIGN KEY (block_id) REFERENCES page_block(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_block_translation_block ON page_block_translation(block_id);

-- 5. Image Tables (Optional)
CREATE TABLE IF NOT EXISTS image (
    id BIGSERIAL PRIMARY KEY,
    url VARCHAR(500) NOT NULL,
    width INT,
    height INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS image_translation (
    image_id BIGINT NOT NULL,
    lang VARCHAR(10) NOT NULL,

    alt VARCHAR(255) NOT NULL,
    title VARCHAR(255),

    PRIMARY KEY (image_id, lang),
    FOREIGN KEY (image_id) REFERENCES image(id) ON DELETE CASCADE
);
