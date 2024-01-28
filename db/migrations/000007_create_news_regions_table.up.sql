CREATE TABLE IF NOT EXISTS news_regions (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT
);
