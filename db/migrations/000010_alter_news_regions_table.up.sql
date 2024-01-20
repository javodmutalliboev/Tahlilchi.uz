BEGIN;

-- Rename existing columns
ALTER TABLE news_regions RENAME COLUMN name TO name_latin;
ALTER TABLE news_regions RENAME COLUMN description TO description_latin;

-- Add new columns with temporary names
ALTER TABLE news_regions ADD COLUMN temp_name_cyrillic text;
ALTER TABLE news_regions ADD COLUMN temp_description_cyrillic text;

-- Update the new columns with some default values
UPDATE news_regions SET temp_name_cyrillic = 'қиймат_' || id;

-- Now rename the new columns to the desired names
ALTER TABLE news_regions RENAME COLUMN temp_name_cyrillic TO name_cyrillic;
ALTER TABLE news_regions RENAME COLUMN temp_description_cyrillic TO description_cyrillic;

-- Add constraints to the new columns
ALTER TABLE news_regions ALTER COLUMN name_cyrillic SET NOT NULL;
ALTER TABLE news_regions ADD CONSTRAINT unique_name_cyrillic UNIQUE(name_cyrillic);

COMMIT;
