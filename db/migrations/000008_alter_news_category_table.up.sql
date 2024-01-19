BEGIN;

-- Rename 'title' to 'title_latin'
ALTER TABLE news_category RENAME COLUMN title TO title_latin;

-- Rename 'description' to 'description_latin'
ALTER TABLE news_category RENAME COLUMN description TO description_latin;

-- Add 'title_cyrillic' column with a default value
ALTER TABLE news_category ADD COLUMN title_cyrillic text DEFAULT 'default_value';

-- Add 'description_cyrillic' column
ALTER TABLE news_category ADD COLUMN description_cyrillic text;

-- Update 'title_cyrillic' with unique values
UPDATE news_category SET title_cyrillic = 'қиймат_' || id;

-- Add the UNIQUE and NOT NULL constraints
ALTER TABLE news_category ALTER COLUMN title_cyrillic SET NOT NULL;
ALTER TABLE news_category ADD CONSTRAINT title_cyrillic_unique UNIQUE(title_cyrillic);

COMMIT;
