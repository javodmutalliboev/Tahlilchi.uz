CREATE TABLE IF NOT EXISTS article_photos (
    id SERIAL PRIMARY KEY,
    article INTEGER NOT NULL REFERENCES articles(id),
    file_name TEXT NOT NULL,
    file bytea NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);