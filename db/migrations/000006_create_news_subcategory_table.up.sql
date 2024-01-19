CREATE TABLE news_subcategory (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES news_category(id),
    title TEXT NOT NULL UNIQUE,
    description TEXT
);