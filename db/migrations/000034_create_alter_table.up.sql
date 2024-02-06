BEGIN;

CREATE TABLE IF NOT EXISTS e_newspaper_category (
    id SERIAL PRIMARY KEY,
    title_latin TEXT NOT NULL UNIQUE,
    title_cyrillic TEXT NOT NULL UNIQUE
);

ALTER TABLE e_newspapers 
ADD COLUMN category INTEGER,
ADD FOREIGN KEY (category) REFERENCES e_newspaper_category(id);

COMMIT;