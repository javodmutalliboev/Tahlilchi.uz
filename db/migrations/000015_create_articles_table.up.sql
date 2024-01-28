CREATE TABLE IF NOT EXISTS articles (
    id BIGSERIAL PRIMARY KEY,
    title_latin TEXT NOT NULL,
    description_latin TEXT,
    title_cyrillic TEXT NOT NULL,
    description_cyrillic TEXT,
    photos BYTEA[],
    videos BYTEA[],
    cover_image BYTEA,
    tags TEXT[]
);