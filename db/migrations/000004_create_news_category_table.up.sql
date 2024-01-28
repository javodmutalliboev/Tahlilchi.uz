CREATE TABLE IF NOT EXISTS news_category (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    description TEXT UNIQUE
);
