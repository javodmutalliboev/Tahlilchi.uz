-- Changing the 'title' column to 'title_latin'
ALTER TABLE photo_gallery RENAME COLUMN title TO title_latin;

-- Adding a new 'title_cyrillic' column
ALTER TABLE photo_gallery ADD COLUMN title_cyrillic TEXT NOT NULL;
