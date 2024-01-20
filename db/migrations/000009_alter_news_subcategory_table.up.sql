BEGIN;

-- Rename existing columns
ALTER TABLE news_subcategory RENAME COLUMN title TO title_latin;
ALTER TABLE news_subcategory RENAME COLUMN description TO description_latin;

-- Add new columns
ALTER TABLE news_subcategory ADD COLUMN title_cyrillic text;
ALTER TABLE news_subcategory ADD COLUMN description_cyrillic text;

-- Set a default value for the new columns to avoid null value errors
-- Generate a unique value for each row
DO
$$
DECLARE
   rec RECORD;
BEGIN
   FOR rec IN SELECT * FROM news_subcategory
   LOOP
      UPDATE news_subcategory SET title_cyrillic = 'қиймат_' || rec.id WHERE id = rec.id;
   END LOOP;
END
$$;

-- Add constraints to the new columns
ALTER TABLE news_subcategory ALTER COLUMN title_cyrillic SET NOT NULL;
ALTER TABLE news_subcategory ADD CONSTRAINT unique_title_cyrillic UNIQUE(title_cyrillic);

COMMIT;
