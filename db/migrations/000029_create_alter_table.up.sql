BEGIN;

CREATE TABLE IF NOT EXISTS bpp_photos (
    id SERIAL PRIMARY KEY,
    bpp INTEGER NOT NULL REFERENCES business_promotional_posts(id),
    file_name TEXT NOT NULL,
    file bytea NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Change the column type
ALTER TABLE business_promotional_posts
ALTER COLUMN videos TYPE text[];

-- Delete the column
ALTER TABLE business_promotional_posts
DROP COLUMN photos;

COMMIT;