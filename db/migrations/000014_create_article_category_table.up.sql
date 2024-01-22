CREATE TABLE article_category (
    id SERIAL PRIMARY KEY,
    title_latin TEXT NOT NULL UNIQUE,
    description_latin TEXT,
    title_cyrillic TEXT NOT NULL UNIQUE,
    description_cyrillic TEXT
);
