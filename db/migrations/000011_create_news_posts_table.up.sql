CREATE TABLE IF NOT EXISTS news_posts (
    id SERIAL PRIMARY KEY,
    title_latin TEXT NOT NULL,
    description_latin TEXT NOT NULL,
    title_cyrillic TEXT NOT NULL,
    description_cyrillic TEXT NOT NULL,
    photo BYTEA,
    video BYTEA,
    audio BYTEA,
    cover_image BYTEA,
    tags TEXT[]
);
