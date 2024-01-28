BEGIN;

-- Add new columns
ALTER TABLE news_posts
    ADD COLUMN category INTEGER REFERENCES news_category(id),
    ADD COLUMN subcategory INTEGER REFERENCES news_subcategory(id),
    ADD COLUMN region INTEGER REFERENCES news_regions(id),
    ADD COLUMN top BOOLEAN,
    ADD COLUMN latest BOOLEAN,
    ADD COLUMN related INTEGER REFERENCES news_posts(id);

-- Change column type
ALTER TABLE news_posts
    ALTER COLUMN video TYPE TEXT;

-- Rename column
ALTER TABLE news_posts
    RENAME COLUMN edited_at TO updated_at;

COMMIT;
