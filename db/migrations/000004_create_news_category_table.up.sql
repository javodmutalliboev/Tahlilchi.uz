CREATE TABLE news_category (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    description TEXT UNIQUE
);
